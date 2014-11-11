package events

import (
	dc "github.com/fsouza/go-dockerclient"
	"log"
	"time"
)

func (self DockerClient) MonitorEvents() {
	//self.getExistingContainerStatus()
	log.Println("Start monitoring container events ...")
	listener, err := self.addEventListener()
	if err != nil {
		SendError(err)
		log.Fatalf("Failed to add event listener: %s\n", err)
	}
	defer self.removeEventListener(listener)
	for {
		select {
		case event := <-listener:
			self.handleEvents(event)
		}
	}
}

func (self DockerClient) handleEvents(event *dc.APIEvents) {
	log.Printf("Docker event:%s %s at %d\n", event.ID, event.Status, event.Time)
	switch event.Status {
	case "start", "die":
		inspect, err := self.inspect(event.ID)
		var isRunning bool
		var exitCode int
		if err != nil {
			isRunning, exitCode = false, 0
		} else {
			isRunning = inspect.State.Running
			exitCode = inspect.State.ExitCode
		}
		go SendContainerEvent(event.ID, isRunning, exitCode, event.Time)
	}
}

func (self DockerClient) getExistingContainerStatus() {
	log.Println("Collecting existing container status")
	containers, err := self.ps(&dc.ListContainersOptions{All: true})
	if err != nil {
		SendError(err)
		log.Fatal(err)
	}
	for _, container := range containers {
		inspect, err := self.inspect(container.ID)
		if err != nil {
			log.Println(err)
			continue
		}
		go SendContainerEvent(container.ID, inspect.State.Running,
			inspect.State.ExitCode, time.Now().UTC().Unix())
	}
	log.Println("Finish collecting existing container status")
}
