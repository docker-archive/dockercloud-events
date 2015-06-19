package main //import "github.com/tutumcloud/events"

import (
	. "github.com/tutumcloud/events/docker-events"
	"log"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func init() {
	runtime.GOMAXPROCS(4)
}

const (
	apiEndpoint = "/api/agent/container/event/"
	version     = "0.1"
)

func main() {
	TutumAuth = os.Getenv("TUTUM_AUTH")
	NodeUUID = os.Getenv("NODE_UUID")
	if TutumAuth == "**None**" {
		log.Fatal("TUTUM_AUTH must be specified")
	}
	if NodeUUID == "**None**" {
		log.Fatal("NodeUUID must be specified")
	}

	DockerHost = os.Getenv("DOCKER_HOST")
	DockerBinary = os.Getenv("DOCKER_BINARY")
	DSN = os.Getenv("SENTRY_DSN")
	TutumHost := os.Getenv("TUTUM_HOST")

	intervalStr := os.Getenv("REPORT_INTERVAL")
	interval, err := strconv.Atoi(intervalStr)
	if err != nil {
		ReportInterval = 30
	} else {
		ReportInterval = interval
	}

	TutumEndpoint = JoinURL(TutumHost, apiEndpoint)
	UserAgent = "tutum-events/" + version

	log.Println("Using Tutum Endpoint:", TutumEndpoint)
	log.Printf("Using NodeUUID(%s), TutumAuth(%s)", NodeUUID, TutumAuth)

	client, err := NewDockerClient(DockerHost)
	if err != nil {
		SendError(err, "Fatal: Failed to get docker client", nil)
		log.Fatalf("Docker %s:%s", err.Error(), DockerHost)
	}

	client.MonitorEvents()
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
