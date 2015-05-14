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
)

func SendContainerEvent(event Event) {
	data, err := json.Marshal(event)
	if err != nil {
		log.Printf("Cannot marshal the posting data: %s\n", event)
	}

	counter := 0
	for {
		err := sendData(TutumEndpoint, data)
		if err == nil {
			break
		} else {
			if counter > 720 {
				log.Println("Too many reties, give up")
				break
			} else {
				counter += 1
				log.Println("Retry in 5 seconds")
				time.Sleep(5 * time.Second)
			}
		}
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

	if resp.StatusCode >= 400 {
		log.Printf("Send metrics failed: %s - %s", resp.Status, string(data))
		extra := map[string]interface{}{"data": string(data)}
		SendError(errors.New(resp.Status), "http error", extra)
		if resp.StatusCode >= 500 {
			return errors.New(resp.Status)
		}
	}
	return nil
}
