package events

import (
	"log"
	"time"

	dc "github.com/fsouza/go-dockerclient"
)

type Event struct {
	Node       string `json:"node"`
	Status     string `json:"status"`
	ID         string `json:"id"`
	From       string `json:"from"`
	Time       int64  `json:"time"`
	HandleTime int64  `json:"handletime"`
	Inspect    string `json:"inspect"`
}

var NodeUUID string

func (self DockerClient) MonitorEvents() {
	log.Println("Start monitoring container events ...")
	listener, err := self.addEventListener()
	if err != nil {
		SendError(err)
		log.Fatalf("Failed to add event listener: %s\n", err)
	}
	defer self.removeEventListener(listener)
	for {
		select {
		case apiEvent := <-listener:
			self.handleEvents(apiEvent)
		}

	}
}

func (self DockerClient) handleEvents(apiEvent *dc.APIEvents) {
	handle_time := time.Now().UnixNano()
	log.Printf("Docker event(handled at %d):%s %s at %d\n",
		handle_time, apiEvent.ID, apiEvent.Status, apiEvent.Time)

	inspect := self.inspect(apiEvent.ID)
	event := Event{NodeUUID, apiEvent.Status, apiEvent.ID, apiEvent.From,
		apiEvent.Time, handle_time, inspect}
	go SendContainerEvent(event)
}
