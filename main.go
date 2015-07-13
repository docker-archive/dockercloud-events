package main //import "github.com/tutumcloud/events"

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"github.com/getsentry/raven-go"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"
)

type Event struct {
	Node       string `json:"node",omitempty`
	Status     string `json:"status"`
	ID         string `json:"id"`
	From       string `json:"from"`
	Time       int64  `json:"time"`
	HandleTime int64  `json:"handletime"`
	ExitCode   string `json:"exitcode"`
}

func init() {
	runtime.GOMAXPROCS(4)
}

const (
	VERSION     = "0.1"
	DockerPath  = "/usr/bin/docker"
	ApiEndpoint = "api/agent/container/event/"
)

var (
	AutorestartEvents = make([]Event, 0)
	UserAgent         = "tutum-events/" + VERSION
	ReportInterval    int
	TutumAuth         string
	TutumUrl          string
	sentryClient      *raven.Client = nil
	DSN               string
	NodeUUID          string
)

func main() {
	TutumAuth = os.Getenv("TUTUM_AUTH")
	TutumUrl = os.Getenv("TUTUM_URL")
	if TutumAuth == "**None**" {
		log.Fatal("TUTUM_AUTH must be specified")
	}
	if TutumUrl == "**None**" {
		TutumHost := os.Getenv("TUTUM_HOST")
		NodeUUID = os.Getenv("NODE_UUID")
		if strings.HasSuffix(TutumHost, "/") {
			TutumUrl = TutumHost + ApiEndpoint
		} else {
			TutumUrl = TutumHost + "/" + ApiEndpoint
		}
		if TutumUrl == "" {
			log.Fatal("TUTUM_URL must be specified")
		}
	}

	DSN = os.Getenv("SENTRY_DSN")

	if !fileExist(DockerPath) {
		log.Fatal("docker client is not mounted to", DockerPath)
	}

	intervalStr := os.Getenv("REPORT_INTERVAL")

	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		ReportInterval = 30
	} else {
		ReportInterval = interval
	}

	log.Println("POST docker event to:", TutumUrl)

	cmd := exec.Command(DockerPath, "version")
	if err := cmd.Start(); err != nil {
		sendError(err, "Fatal: Failed to run docker version", nil)
		log.Println(err)
		return
	}
	cmd.Wait()

	monitorEvents()
}

func monitorEvents() {
	ticker := time.NewTicker(time.Second * time.Duration(ReportInterval))
	go func() {
		for {
			select {
			case <-ticker.C:
				events := AutorestartEvents
				AutorestartEvents = make([]Event, 0)
				if len(events) > 0 {
					go sendContainerAutoRestartEvents(events)
				}
			}
		}
	}()

	for {
		log.Println("docker events starts")
		cmd := exec.Command(DockerPath, "events")
		cmdReader, err := cmd.StdoutPipe()
		if err != nil {
			log.Fatal("Error creating StdoutPipe for Cmd", err)
		}

		scanner := bufio.NewScanner(cmdReader)
		go func() {
			for scanner.Scan() {
				eventStr := scanner.Text()
				if eventStr != "" {
					re := regexp.MustCompile("(.*) (.{64}): \\(from (.*)\\) (.*)")
					terms := re.FindStringSubmatch(eventStr)
					if len(terms) == 5 {
						var event Event
						if NodeUUID != "" {
							event.Node = NodeUUID
						}
						eventTime, err := time.Parse(time.RFC3339Nano, terms[1])
						if err == nil {
							event.Time = eventTime.Unix()
						} else {
							event.Time = time.Now().Unix()
						}
						event.ID = terms[2]
						event.From = terms[3]
						event.Status = terms[4]
						event.HandleTime = time.Now().UnixNano()
						go eventHandler(event)
					}
				}
			}
		}()

		err = cmd.Start()
		if err != nil {
			log.Print("Error starting docker evets", err)
			break
		}

		err = cmd.Wait()
		if err != nil {
			log.Print("Error waiting for docker events", err)
			break
		}
		log.Println("docker events stops")
	}
}

func eventHandler(event Event) {
	event.ExitCode = "0"
	isRestart := false
	if strings.ToLower(event.Status) == "start" ||
		strings.ToLower(event.Status) == "die" {

		result, err := exec.Command(DockerPath, "inspect", "-f",
			"{{index .HostConfig.RestartPolicy.Name}} {{index .State.ExitCode}}",
			event.ID).Output()

		if err == nil && len(result) > 0 {
			terms := strings.Split(string(result), " ")
			if len(terms) == 2 {
				if strings.HasPrefix(strings.ToLower(terms[0]), "on-failure") ||
					strings.HasPrefix(strings.ToLower(terms[0]), "always") {
					isRestart = true
				}
				event.ExitCode = strings.Trim(terms[1], "\n")
			}
		}
		if isRestart {
			AutorestartEvents = append(AutorestartEvents, event)
		} else {
			sendContainerEvent(event)
		}
	}
}

func sendContainerAutoRestartEvents(events []Event) {
	data, err := json.Marshal(events)
	if err != nil {
		log.Printf("Cannot marshal the posting data: %s\n", events)
	}
	sendData(data)
}

func sendContainerEvent(event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Cannot marshal the posting data: %s\n", event)
	}
	sendData(data)
}

func sendData(data []byte) {
	counter := 1
	for {
		log.Println("sending event: ", string(data))
		err := send(TutumUrl, data)
		if err == nil {
			break
		} else {
			if counter > 100 {
				log.Println("Too many reties, give up")
				break
			} else {
				counter *= 2
				log.Printf("%s: Retry in %d seconds", err, counter)
				time.Sleep(time.Duration(counter) * time.Second)
			}
		}
	}
}

func send(url string, data []byte) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		sendError(err, "Failed to create http.NewRequest", nil)
		return err
	}
	req.Header.Add("Authorization", TutumAuth)
	req.Header.Add("User-Agent", UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		extra := map[string]interface{}{"data": string(data)}
		sendError(err, "Failed to POST the http request", extra)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Printf("Send event failed: %s - %s", resp.Status, string(data))
		extra := map[string]interface{}{"data": string(data)}
		sendError(errors.New(resp.Status), "http error", extra)
		if resp.StatusCode >= 500 {
			return errors.New(resp.Status)
		}
	}
	return nil
}

func fileExist(file string) bool {
	if _, err := os.Stat(file); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

func getSentryClient() *raven.Client {
	if sentryClient == nil && DSN != "" {
		client, _ := raven.NewClient(DSN, nil)
		sentryClient = client
	}
	return sentryClient
}

func sendError(err error, msg string, extra map[string]interface{}) {
	go func() {
		client := getSentryClient()
		if sentryClient != nil {
			packet := &raven.Packet{Message: msg, Interfaces: []raven.Interface{raven.NewException(err, raven.NewStacktrace(0, 5, nil))}}
			if extra != nil {
				packet.Extra = extra
			}
			_, ch := client.Capture(packet, nil)
			<-ch
		}
	}()
}
