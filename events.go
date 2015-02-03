package main

import (
	"flag"
	"github.com/tutumcloud/container-events/events"
	"log"
	"runtime"
	"strings"
)

func init() {
	runtime.GOMAXPROCS(4)
}

var pDockerHost = flag.String("dockerHost", "unix:///var/run/docker.sock", "docker host")
var pTest = flag.Bool("test", false, "test if the running environment is correct")
var pDockerBinary = flag.String("dockerBinary", "/docker", "docker binary")

const (
	apiEndpoint = "/api/agent/container/event/"
	configFile  = "/etc/tutum/agent/tutum-agent.conf"
)

func main() {
	flag.Parse()

	events.DockerHost = *pDockerHost
	events.DockerBinary = *pDockerBinary

	conf := events.GetConf(configFile)
	events.TutumEndpoint = JoinURL(conf.TutumHost, apiEndpoint)
	events.NodeUUID = conf.TutumUUID
	events.TutumToken = conf.TutumToken

	if *pTest == false {
		log.Println("Using Tutum Endpoint:", events.TutumEndpoint)
		log.Printf("Using NodeUUID(%s), TutumToken(%s)\n", events.NodeUUID, events.TutumToken)
	}

	client, err := events.NewDockerClient(events.DockerHost)
	if err != nil {
		events.SendError(err)
		log.Fatalf("Docker %s:%s", err.Error(), events.DockerHost)
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
