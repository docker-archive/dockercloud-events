package main

import (
	"flag"
	. "github.com/tutumcloud/container-events/events"
	"log"
	"os"
	"runtime"
	"strings"
)

func init() {
	runtime.GOMAXPROCS(4)
}

const (
	apiEndpoint = "/api/agent/container/event/"
)

func main() {
	var pTest = flag.Bool("test", false, "test if the execution environment is correct")
	flag.Parse()

	DockerHost = os.Getenv("DOCKER_HOST")
	DockerBinary = os.Getenv("DOCKER_BINARY")
	TutumAuth = os.Getenv("TUTUM_AUTH")
	NodeUUID = os.Getenv("NODE_UUID")
	DSN = os.Getenv("SENTRY_DSN")
	TutumHost := os.Getenv("TUTUM_HOST")

	TutumEndpoint = JoinURL(TutumHost, apiEndpoint)

	if *pTest == false {
		log.Println("Using Tutum Endpoint:", TutumEndpoint)
		log.Printf("Using NodeUUID(%s), TutumAuth(%s)", NodeUUID, TutumAuth)
	}

	client, err := NewDockerClient(DockerHost)
	if err != nil {
		SendError(err, "Fatal: Failed to get docker client", nil)
		log.Fatalf("Docker %s:%s", err.Error(), DockerHost)
	}
	if *pTest == false {
		client.MonitorEvents()
	}
}

func JoinURL(url1 string, url2 string) (url string) {
	if strings.HasSuffix(url1, "/") {
		if strings.HasPrefix(url2, "/") {
			url = url1 + url2[1:]
		} else {
			url = url1 + url2
		}
	} else {
		if strings.HasPrefix(url2, "/") {
			url = url1 + url2
		} else {
			url = url1 + "/" + url2
		}
	}
	if !strings.HasSuffix(url, "/") {
		url = url + "/"
	}
	return
}
