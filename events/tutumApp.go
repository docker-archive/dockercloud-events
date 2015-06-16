package events

import (
	"bytes"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"
)

var (
	TutumEndpoint string
	TutumAuth     string
	UserAgent     string
)

func SendContainerAutoRestartEvents(events []Event) {
	data, err := json.Marshal(events)
	if err != nil {
		log.Printf("Cannot marshal the posting data: %s\n", events)
	}
	sendData(data)
}

func SendContainerEvent(event Event) {
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
		err := send(TutumEndpoint, data)
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
		SendError(err, "Failed to create http.NewRequest", nil)
		return err
	}
	req.Header.Add("Authorization", TutumAuth)
	req.Header.Add("User-Agent", UserAgent)
	resp, err := client.Do(req)
	if err != nil {
		extra := map[string]interface{}{"data": string(data)}
		SendError(err, "Failed to POST the http request", extra)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		log.Printf("Send event failed: %s - %s", resp.Status, string(data))
		extra := map[string]interface{}{"data": string(data)}
		SendError(errors.New(resp.Status), "http error", extra)
		if resp.StatusCode >= 500 {
			return errors.New(resp.Status)
		}
	}
	return nil
}
