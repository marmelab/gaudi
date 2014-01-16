package container

import (
	"fmt"
	"strings"
	"time"
	"launchpad.net/goyaml"
	"github.com/marmelab/arch-o-matic/docker"
)

type Container struct {
	Name         string
	Type         string
	InstanceType string
	Running      bool
	Id           string
	Ip           string
	Links        []string
	Dependencies []*Container
	Ports        map[string]string
	Volumes      map[string]string
	Custom       map[string]interface{}
}

type inspection struct {
	NetworkSettings map[string]string "NetworkSettings,omitempty"
}

func (c *Container) init() {
	if c.Ports == nil {
		c.Ports = make(map[string]string)
	}
	if c.Links == nil {
		c.Links = make([]string, 0)
	}
	if c.Dependencies == nil {
		c.Dependencies = make([]*Container, 0)
	}
}

func (c *Container) Clean(done chan bool) {
	fmt.Println("Cleaning", c.Name, "...")
	docker.Clean(c.Name)
	c.Running = false

	done <- true
}

func (c *Container) Build(done chan bool) {
	fmt.Println("Building", c.Name, "...")
	docker.Build(c.Name)

	done <- true
}

func (c *Container) IsRunning() bool {
	return c.Running
}

func (c *Container) IsReady() bool {
	ready := true

	for _, dependency := range c.Dependencies {
		ready = ready && dependency.IsRunning()
	}

	return ready
}

func (c *Container) AddDependency(container *Container) {
	c.Dependencies = append(c.Dependencies, container)
}

func (c *Container) Start() {
	c.init()

	startResult := docker.Start(c.Name, c.Links, c.Ports, c.Volumes)
	c.Id = strings.TrimSpace(startResult)
	c.Running = true

	time.Sleep(2 * time.Second)
	c.retrieveIp()

	fmt.Println("Container", c.Name, "started", "(" +c.Ip+":"+c.GetFirstPort()+") :", c.Id)
}

func (c *Container) GetCustomValue(name string) interface{} {
	return c.Custom[name]
}

func (c *Container) GetCustomValueAsString(name string) string {
	return c.Custom[name].(string)
}

func (c *Container) GetFirstPort() string {
	for key, _ := range c.Ports {
		return key
	}

	return ""
}

func (c *Container) retrieveIp () {
	inspect := docker.Inspect(c.Id)

	var results []inspection
	goyaml.Unmarshal(inspect, &results)

	c.Ip = results[0].NetworkSettings["IPAddress"]
}
