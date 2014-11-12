package events

import (
	dc "github.com/fsouza/go-dockerclient"
	"log"
)

type Event struct {
	Node    string `json:"node"`
	Status  string `json:"status"`
	ID      string `json:"id"`
	From    string `json:"from"`
	Time    int64  `json:"time"`
	Inspect string `json:"inspect"`
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
	log.Printf("Docker event:%s %s at %d\n", apiEvent.ID, apiEvent.Status, apiEvent.Time)

	inspect := self.inspect(apiEvent.ID)
	event := Event{NodeUUID, apiEvent.Status, apiEvent.ID, apiEvent.From, apiEvent.Time, inspect}
	go SendContainerEvent(event)
}
