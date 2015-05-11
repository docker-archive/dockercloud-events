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
		SendError(err, "Fatal: Failed to add event listener", nil)
		log.Fatal("Failed to add event listener:", err)
	}
	defer self.removeEventListener(listener)

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
	handle_time := time.Now().UnixNano()

	inspect := self.inspect(apiEvent.ID)
	event := Event{NodeUUID, apiEvent.Status, apiEvent.ID, apiEvent.From,
		apiEvent.Time, handle_time, inspect}
	SendContainerEvent(event)
}
