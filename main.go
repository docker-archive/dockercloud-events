package main //import "github.com/tutumcloud/events"

import (
	"bufio"
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/getsentry/raven-go"
)

type Event struct {
	Status     string `json:"status"`
	ID         string `json:"id"`
	From       string `json:"from"`
	Time       int64  `json:"time"`
	HandleTime int64  `json:"handletime"`
	ExitCode   string `json:"exitcode"`
}

type ContainerState struct {
	isRunning bool
	created   int64
	updated   int64
}

func init() {
	runtime.GOMAXPROCS(4)
}

const (
	VERSION    = "1.0"
	DockerPath = "/usr/bin/docker"
)

var (
	UserAgent      = "events-daemon/" + VERSION
	Interval       int
	Auth           string
	ApiUrl         string
	sentryClient   *raven.Client = nil
	DSN            string
	Container      = make(map[string]*ContainerState)
	FlagStandalone *bool
)

func main() {
	FlagStandalone = flag.Bool("standalone", false, "Standalone mode")
	flag.Parse()

	Auth = os.Getenv("DOCKERCLOUD_AUTH")
	ApiUrl = os.Getenv("EVENTS_API_URL")
	if Auth == "**None**" {
		log.Fatal("DOCKERCLOUD_AUTH must be specified")
	}
	if ApiUrl == "**None**" {
		log.Fatal("EVENTS_API_URL must be specified")
	}

	DSN = os.Getenv("SENTRY_DSN")

	if !fileExist(DockerPath) {
		log.Fatal("docker client is not mounted to", DockerPath)
	}

	intervalStr := os.Getenv("REPORT_INTERVAL")

	if interval, err := strconv.Atoi(intervalStr); err == nil {
		Interval = interval
	} else {
		Interval = 5
	}

	if *FlagStandalone {
		log.Print("Running in standalone mode")
	} else {
		log.Print("POST docker event to: ", ApiUrl)
	}

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
	log.Println("docker events starts")
	cmd := exec.Command(DockerPath, "events")
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal("Error creating StdoutPipe for Cmd", err)
	}

	err = cmd.Start()
	if err != nil {
		log.Fatal("Error starting docker evets", err)
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

					state := strings.ToLower(event.Status)
					if state == "start" || state == "die" {
						updateContainerState(&event)
						go eventHandler(&event)
					}
				}
			}
		}
		if scanner.Err() == nil {
			log.Fatal("The scanner returns an error:", "EOF")
		} else {
			log.Fatal("The scanner returns an error:", scanner.Err())
		}
	}()

	err = cmd.Wait()
	if err != nil {
		log.Fatal("Error waiting for docker events", err)
	}
	log.Println("docker events stops")
}

func updateContainerState(event *Event) {
	isRunning := false
	if strings.ToLower(event.Status) == "start" {
		isRunning = true
	}
	container := Container[event.ID]
	if container == nil {
		Container[event.ID] = &ContainerState{isRunning: isRunning, updated: event.HandleTime, created: event.HandleTime}
	} else {
		container.updated = event.HandleTime
		container.isRunning = isRunning
	}
}

func eventHandler(event *Event) {
	exitcode, isAutoRestart := getContainerStatus(event)
	event.ExitCode = exitcode

	if isAutoRestart {
		status := strings.ToLower(event.Status)
		if status == "die" {
			container := Container[event.ID]
			if container != nil && container.created == container.updated {
				sendContainerEvent(event)
			}
		}
		if status == "start" {
			container := Container[event.ID]
			if container == nil {
				sendContainerEvent(event)
			} else {
				if container.created == container.updated {
					delete(Container, event.ID)
					sendContainerEvent(event)
				} else {
					go delaySendContainerEvent(event)
				}
			}
		}
	} else {
		delete(Container, event.ID)
		sendContainerEvent(event)
	}

}

func getContainerStatus(event *Event) (exitcode string, isAutoRestart bool) {
	exitcode = "0"
	isAutoRestart = false
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
					isAutoRestart = true
				}
				exitcode = strings.Trim(terms[1], "\n")
			}
		}
	}
	return
}

func sendContainerEvent(event *Event) {
	data, err := json.Marshal(*event)
	if err != nil {
		log.Printf("Cannot marshal the posting data: %s\n", event)
	}
	if *FlagStandalone {
		log.Print("Send event: ", string(data))
	} else {
		sendData(data)
	}
}

func delaySendContainerEvent(event *Event) {
	time.Sleep(time.Duration(Interval) * time.Second)
	container := Container[event.ID]
	if container == nil {
		log.Print("No container found")
		sendContainerEvent(event)
	} else {
		currentTime := time.Now().UnixNano()
		if currentTime-container.updated >= int64(Interval)*1000000000 && container.isRunning {
			delete(Container, event.ID)
			log.Printf("Run over %d", currentTime-container.updated)
			sendContainerEvent(event)
		}
	}
}

func sendData(data []byte) {
	counter := 1
	for {
		log.Println("sending event: ", string(data))
		err := send(ApiUrl, data)
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
	req.Header.Add("Authorization", Auth)
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
