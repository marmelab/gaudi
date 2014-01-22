package container

import (
	"fmt"
	"strings"
	"time"
	"launchpad.net/goyaml"
	"github.com/marmelab/gaudi/docker"
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
	ID string "ID,omitempty"
	NetworkSettings map[string]string "NetworkSettings,omitempty"
	State map[string]interface{} "State,omitempty"
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

func (c *Container) Remove(done chan bool) {
	docker.Remove(c.Name)
	c.Running = false

	done <- true
}

func (c *Container) Kill(done chan bool, silent bool) {
	if !silent {
		fmt.Println("Stopping", c.Name, "...")
	}
	docker.Kill(c.Name)
	c.Running = false

	done <- true
}

func (c *Container) Clean(done chan bool) {
	fmt.Println("Cleaning", c.Name, "...")

	cleaning := make(chan bool, 2)
	c.Kill(cleaning, true)
	c.Remove(cleaning)

	<-cleaning

	done <- true
}


func (c *Container) Build(done chan bool) {
	fmt.Println("Building", c.Name, "...")
	docker.Build(c.Name)

	done <- true
}

func (c *Container) IsRunning() bool {
	if c.Running {
		return true
	}

	// Check if a container with the same name is already running
	inspect, err := docker.Inspect(c.Name)
	if err != nil {
		return false
	}

	c.retrieveInfoFromInspection(inspect)
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

	if c.IsRunning() {
		fmt.Println("Application", c.Name, "already running", "(" +c.Ip+":"+c.GetFirstPort()+")")
		return
	}

	startResult := docker.Start(c.Name, c.Links, c.Ports, c.Volumes)
	c.Id = strings.TrimSpace(startResult)

	time.Sleep(time.Second)
	c.retrieveIp()
	c.Running = true

	fmt.Println("Application", c.Name, "started", "(" +c.Ip+":"+c.GetFirstPort()+")")
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

func (c *Container) CheckIfRunning() {
	if c.IsRunning() {
		fmt.Println("Application", c.Name, "is running", "(" +c.Ip+":"+c.GetFirstPort()+")")
	} else {
		fmt.Println("Application", c.Name, "is not running")
	}
}

func (c *Container) retrieveIp () {
	inspect, err := docker.Inspect(c.Id)
	if err != nil {
		panic(err)
	}

	c.retrieveInfoFromInspection(inspect)
}

func (c *Container) retrieveInfoFromInspection (inspect []byte) {
	var results []inspection
	goyaml.Unmarshal(inspect, &results)

	c.Running = results[0].State["Running"].(bool)
	c.Ip = results[0].NetworkSettings["IPAddress"]
	c.Id = results[0].ID
}

