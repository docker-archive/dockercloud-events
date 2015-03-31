package events

import (
	"github.com/getsentry/raven-go"
	"log"
)

var sentryClient *raven.Client = nil
var DSN string = "https://cc08ab58e20447aabef069775fee97ca:47b6c9460fed424db44a1d9ed65366ec@app.getsentry.com/31604"

func getSentryClient() *raven.Client {
	if sentryClient == nil {
		client, err := raven.NewClient(DSN, nil)
		if err != nil {
			log.Println(err)
		}
		sentryClient = client
	}
	return sentryClient
}

func SendError(err error, msg string) {
	go func() {
		client := getSentryClient()
		packet := &raven.Packet{Message: msg, Interfaces: []raven.Interface{raven.NewException(err, raven.NewStacktrace(0, 5, nil))}}
		_, ch := client.Capture(packet, nil)
		if senderr := <-ch; senderr != nil {
			log.Println(senderr)
		} else {
			log.Println("sent error to sentry successfully")
		}
	}()
}
