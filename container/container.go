package container

import (
	"fmt"
	"github.com/marmelab/gaudi/docker"
	"launchpad.net/goyaml"
	"strings"
	"time"
)

type Container struct {
	Name         string
	Type         string
	Image        string
	Path         string
	Running      bool
	Id           string
	Ip           string
	Binary       bool
	BeforeScript string   "before_script"
	AfterScript  string   "after_script"
	AptPackets   []string "apt_get"
	Links        []string
	Dependencies []*Container
	Ports        map[string]string
	Volumes      map[string]string
	Custom       map[string]interface{}
}

type inspection struct {
	ID              string                 "ID,omitempty"
	NetworkSettings map[string]string      "NetworkSettings,omitempty"
	State           map[string]interface{} "State,omitempty"
}

func (c *Container) init() {
	if c.Ports == nil {
		c.Ports = make(map[string]string)
	}
	if c.Links == nil {
		c.Links = make([]string, 0)
	}
	if c.AptPackets == nil {
		c.AptPackets = make([]string, 0)
	}
}

func (c *Container) Remove() {
	docker.Remove(c.Name)
	c.Running = false
}

func (c *Container) Kill(silent bool, done chan bool) {
	if !silent {
		fmt.Println("Killing", c.Name, "...")
	}

	docker.Kill(c.Name)
	c.Running = false

	if done != nil {
		done <- true
	}
}

func (c *Container) Clean(done chan bool) {
	fmt.Println("Cleaning", c.Name, "...")

	c.Kill(true, nil)
	c.Remove()

	done <- true
}

func (c *Container) BuildOrPull(buildChans chan bool) {
	if c.IsPreBuild() {
		c.Pull(buildChans)
	} else {
		c.Build(buildChans)
	}
}

func (c *Container) Build(done chan bool) {
	buildName := "gaudi/" + c.Name
	buildPath := "/tmp/gaudi/" + c.Name

	if c.IsRemote() {
		buildName = c.Image
		buildPath = c.Path
	}

	fmt.Println("Building", buildName, "...")
	docker.Build(buildName, buildPath)

	done <- true
}

func (c *Container) Pull(done chan bool) {
	fmt.Println("Pulling", c.Image, "...")
	docker.Pull(c.Image)

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

/**
 * Starts a container as a server
 */
func (c *Container) Start() {
	c.init()

	fmt.Println("Starting", c.Name, "...")

	if c.IsRunning() {
		fmt.Println("Application", c.Name, "already running", "("+c.Ip+":"+c.GetFirstPort()+")")
		return
	}

	startResult := docker.Start(c.Name, c.Image, c.Links, c.Ports, c.Volumes)
	c.Id = strings.TrimSpace(startResult)

	time.Sleep(3 * time.Second)
	c.RetrieveIp()
	c.Running = true

	fmt.Println("Application", c.Name, "started", "("+c.Ip+":"+c.GetFirstPort()+")")
}

func (c *Container) BuildAndRun(currentPath string, arguments []string) {
	buildChans := make(chan bool, 1)
	go c.BuildOrPull(buildChans)
	<-buildChans

	c.Run(currentPath, arguments)
}

/**
 * Starts a container as a binary file
 */
func (c *Container) Run(currentPath string, arguments []string) {
	c.init()

	fmt.Println("Running", c.Name, strings.Join(arguments, " "), "...")

	out := docker.Run(c.Image, currentPath, arguments)

	fmt.Println(out)
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

func (c *Container) GetFirstLocalPort() string {
	for _, localPort := range c.Ports {
		return localPort
	}

	return ""
}

func (c *Container) IsGaudiManaged() bool {
	return !c.IsPreBuild() && !c.IsRemote()
}

func (c *Container) IsPreBuild() bool {
	return c.Type == "prebuild"
}

func (c *Container) IsRemote() bool {
	return c.Type == "remote"
}

func (c *Container) HasBeforeScript() bool {
	return len(c.BeforeScript) != 0
}

func (c *Container) HasBeforeScriptFile() bool {
	return c.HasBeforeScript() && (c.BeforeScript[0] == '.' || c.BeforeScript[0] == '/')
}

func (c *Container) HasAfterScript() bool {
	return len(c.AfterScript) != 0
}

func (c *Container) HasAfterScriptFile() bool {
	return c.HasAfterScript() && (c.AfterScript[0] == '.' || c.AfterScript[0] == '/')
}

func (c *Container) RetrieveIp() {
	inspect, err := docker.Inspect(c.Id)
	if err != nil {
		panic(err)
	}

	c.retrieveInfoFromInspection(inspect)
}

func (c *Container) retrieveInfoFromInspection(inspect []byte) {
	var results []inspection
	goyaml.Unmarshal(inspect, &results)

	var isRunning bool
	rawRunning := results[0].State["Running"]
	if rawRunning != nil {
		isRunning = rawRunning.(bool)
	} else {
		isRunning = false
	}

	c.Running = isRunning

	c.Ip = results[0].NetworkSettings["IPAddress"]
	c.Id = results[0].ID
}
