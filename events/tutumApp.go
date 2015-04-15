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
		SendError(err, "Failed to create http.NewRequest", nil)
		return err
	}
	req.Header.Add("Authorization", TutumAuth)
	resp, err := client.Do(req)
	if err != nil {
		extra := map[string]interface{}{"data": string(data)}
		SendError(err, "Failed to POST the http request", extra)
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 200, 201, 202:
		return nil
	default:
		extra := map[string]interface{}{"data": string(data)}
		SendError(errors.New(resp.Status), "http error", extra)
		return errors.New(resp.Status)
	}
}
