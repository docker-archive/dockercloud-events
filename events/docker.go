package events

import (
	"fmt"
	dc "github.com/fsouza/go-dockerclient"
	"strings"
)

type DockerClient struct{ client *dc.Client }

var (
	DockerHost   string
	DockerBinary string
)

func NewDockerClient(host string) (DockerClient, error) {
	c, err := dc.NewClient(host)
	if err != nil {
		return DockerClient{}, err
	}
	return DockerClient{c}, nil
}

func (self DockerClient) addEventListener() (listener chan *dc.APIEvents, err error) {
	listener = make(chan *dc.APIEvents)
	return listener, self.client.AddEventListener(listener)
}

func (self DockerClient) removeEventListener(listener chan *dc.APIEvents) error {
	return self.client.RemoveEventListener(listener)
}

func (self DockerClient) inspect(id string) (restart bool, exitcode string) {
	restart = false
	exitcode = ""

	container, err := self.client.InspectContainer(id)
	if err == nil && container != nil {
		name := strings.ToLower(container.HostConfig.RestartPolicy.Name)
		if strings.HasPrefix(name, "on-failure") || strings.HasPrefix(name, "always") {
			restart = true
		}

		exitcode = fmt.Sprintf("%d", container.State.ExitCode)
	}
	return
}

func (self DockerClient) ps(opts *dc.ListContainersOptions) ([]dc.APIContainers, error) {
	return self.client.ListContainers(*opts)
}
