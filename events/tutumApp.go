package events

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
)

var (
	TutumEndpoint string
	TutumAuth     string
)

func SendContainerEvent(event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Cannot marshal the posting data: %s\n", event)
	}

	log.Printf("Sending container event: %s on %s \n", event.Status, event.ID)
	if err := sendData(TutumEndpoint, data); err != nil {
		log.Printf("Error when Posting %s:%s\n", data, err.Error())
	}
}

func sendData(url string, data []byte) error {
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(data))
	if err != nil {
		SendError(err, "Failed to create http.NewRequest")
		return err
	}
	req.Header.Add("Authorization", TutumAuth)
	resp, err := client.Do(req)
	if err != nil {
		SendError(err, "Failed to do the http request")
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201, 202:
		return nil
	default:
		SendError(errors.New(resp.Status), "http error")
		return errors.New(resp.Status)
	}
}
