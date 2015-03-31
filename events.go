package main

import (
	"flag"
	. "github.com/tutumcloud/container-events/events"
	"log"
	"runtime"
	"strings"
)

func init() {
	runtime.GOMAXPROCS(4)
}

var pTest = flag.Bool("test", false, "test if the running environment is correct")
var pDockerBinary = flag.String("dockerBinary", "/docker", "docker binary")
var pDockerHost = flag.String("dockerHost", "unix:///var/run/docker.sock", "docker host")
var pTutumHost = flag.String("tutumHost", "https://dashboard.tutum.co/", "tutum host")
var pTutumAuth = flag.String("auth", "", "tutum auth")
var pNodeUUID = flag.String("uuid", "", "node uuid")

const (
	apiEndpoint = "/api/agent/container/event/"
	configFile  = "/etc/tutum/agent/tutum-agent.conf"
)

func main() {
	flag.Parse()

	DockerHost = *pDockerHost
	DockerBinary = *pDockerBinary
	TutumAuth = *pTutumAuth
	NodeUUID = *pNodeUUID

	TutumEndpoint = JoinURL(*pTutumHost, apiEndpoint)

	if *pTest == false {
		log.Println("Using Tutum Endpoint:", TutumEndpoint)
		log.Printf("Using NodeUUID(%s), TutumAuth(%s)", NodeUUID, TutumAuth)
	}

	client, err := NewDockerClient(DockerHost)
	if err != nil {
		SendError(err, "Fatal: Failed to get docker client")
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
