package events

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

type ContainerStatus struct {
	NodeUUID    string `json:"node_uuid"`
	ContainerID string `json:"docker_id"`
	IsRunning   bool   `json:"is_running"`
	ExitCode    int    `json:"exit_code"`
	TimeStamp   int64  `json:"timestamp"`
}

var (
	TutumEndpoint string
	NodeUUID      string
	TutumToken    string
)

func SendContainerEvent(id string, isRunning bool, exitCode int, timestamp int64) {
	data := ContainerStatus{NodeUUID, id, isRunning, exitCode, timestamp}
	form, err := json.Marshal(data)
	if err != nil {
		log.Printf("Cannot marshal the posting data: %v\n", data)
	}

	log.Printf("Sending container event: %v\n", data)
	if err := sendData(TutumEndpoint, form); err != nil {
		log.Printf("Error when Posting %s:%s\n", form, err.Error())
	}
}

func sendData(url string, form []byte) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(form))
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "TutumAgentToken "+TutumToken)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201, 202:
		log.Println(resp.Status)
		return nil
	default:
		return errors.New(resp.Status)
	}
}
