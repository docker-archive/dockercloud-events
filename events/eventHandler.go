package events

import (
	dc "github.com/fsouza/go-dockerclient"
	"log"
	"strings"
	"time"
)

type Event struct {
	Node       string `json:"node"`
	Status     string `json:"status"`
	ID         string `json:"id"`
	From       string `json:"from"`
	Time       int64  `json:"time"`
	HandleTime int64  `json:"handletime"`
	ExitCode   string `json:"exitcode"`
}

var (
	AutorestartEvents = make([]Event, 0)
	NodeUUID          string
	ReportInterval    int
)

func (self DockerClient) MonitorEvents() {
	log.Println("Start monitoring container events ...")
	listener, err := self.addEventListener()
	if err != nil {
		SendError(err, "Fatal: Failed to add event listener", nil)
		log.Fatal("Failed to add event listener:", err)
	}
	defer self.removeEventListener(listener)

	ticker := time.NewTicker(time.Second * time.Duration(ReportInterval))
	go func() {
		for {
			select {
			case <-ticker.C:
				events := AutorestartEvents
				AutorestartEvents = make([]Event, 0)
				if len(events) > 0 {
					go SendContainerAutoRestartEvents(events)
				}
			}
		}
	}()

	timeout := time.After(1 * time.Second)
	for {
		select {
		case apiEvent := <-listener:
			go self.handleEvents(apiEvent)
		case <-timeout:
			break
		}
	}
}

func (self DockerClient) handleEvents(apiEvent *dc.APIEvents) {
	if strings.ToLower(apiEvent.Status) == "start" ||
		strings.ToLower(apiEvent.Status) == "die" {
		handle_time := time.Now().UnixNano()
		autorestart, exitcode := self.inspect(apiEvent.ID)
		event := Event{NodeUUID, apiEvent.Status, apiEvent.ID, apiEvent.From,
			apiEvent.Time, handle_time, exitcode}
		if autorestart {
			AutorestartEvents = append(AutorestartEvents, event)
		} else {
			SendContainerEvent(event)
		}
	}
}
