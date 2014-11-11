package events

import (
	dc "github.com/fsouza/go-dockerclient"
)

type DockerClient struct{ client *dc.Client }

func NewDockerClient(host string) (DockerClient, error) {
	c, err := dc.NewClient(host)
	if err != nil {
		return DockerClient{}, err
	}
	return DockerClient{c}, nil
}

func (self DockerClient) addEventListener() (listener chan *dc.APIEvents, err error) {
	listener = make(chan *dc.APIEvents, 10)
	return listener, self.client.AddEventListener(listener)
}

func (self DockerClient) removeEventListener(listener chan *dc.APIEvents) error {
	return self.client.RemoveEventListener(listener)
}

func (self DockerClient) inspect(id string) (*dc.Container, error) {
	return self.client.InspectContainer(id)
}

func (self DockerClient) ps(opts *dc.ListContainersOptions) ([]dc.APIContainers, error) {
	return self.client.ListContainers(*opts)
}
