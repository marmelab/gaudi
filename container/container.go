package container

import (
	"fmt"
	"github.com/marmelab/gaudi/docker"
	"github.com/marmelab/gaudi/util"
	"launchpad.net/goyaml"
	"strings"
	"time"
)

type Ambassador struct {
	Type   string
	Remote string
	Port   string
}

type Container struct {
	Name         string
	Type         string
	Image        string
	Path         string
	Running      bool
	Id           string
	Ip           string
	BeforeScript string            "before_script"
	AfterScript  string            "after_script"
	AptPackets   []string          "apt_get"
	Add          map[string]string "add"
	Links        []string
	Dependencies []*Container
	Ambassador   Ambassador
	Ports        map[string]string
	Volumes      map[string]string
	Environments map[string]string "environments"
	Custom       map[string]interface{}
}

type inspection struct {
	ID              string                 "ID,omitempty"
	NetworkSettings map[string]string      "NetworkSettings,omitempty"
	State           map[string]interface{} "State,omitempty"
}

func (c *Container) Init() {
	if c.Ports == nil {
		c.Ports = make(map[string]string)
	}
	if c.Links == nil {
		c.Links = make([]string, 0)
	}
	if c.AptPackets == nil {
		c.AptPackets = make([]string, 0)
	}
	if c.Environments == nil {
		c.Environments = make(map[string]string)
	}
	if c.Add == nil {
		c.Add = make(map[string]string)
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
		// remote type is deprecated
		if c.Type == "remote" {
			fmt.Println("WARN: 'remote' type is deprecated, use 'github' instead")
		}

		buildName = c.Image
		buildPath = c.Path
	}

	fmt.Println("Building", buildName, "...")
	docker.Build(buildName, buildPath)

	done <- true
}

func (c *Container) Pull(done chan bool) {
	fmt.Println("Pulling", c.Image, "...")

	// prebuild type is deprecated
	if c.Type == "prebuild" {
		fmt.Println("WARN: 'prebuild' type is deprecated, use 'index' instead")
	}

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
func (c *Container) Start(rebuild bool) {
	c.Init()

	// Check if the container is already running
	if !rebuild {
		if c.IsRunning() {
			fmt.Println("Application", c.Name, "is already running", "("+c.Ip+":"+c.GetFirstPort()+")")
			return
		}

		cleanChan := make(chan bool, 1)
		go c.Clean(cleanChan)
		<-cleanChan
	}

	fmt.Println("Starting", c.Name, "...")

	startResult := docker.Start(c.Name, c.Image, c.Links, c.Ports, c.Volumes, c.Environments)
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
	c.Init()

	fmt.Println("Running", c.Name, strings.Join(arguments, " "), "...")

	out := docker.Run(c.Image, currentPath, arguments)

	fmt.Println(out)
}

func (c *Container) GetCustomValue(params ...string) interface{} {
	if value, ok := c.Custom[params[0]]; ok {
		return value
	}

	if len(params) == 2 {
		return params[1]
	}

	return nil
}

func (c *Container) GetCustomValueAsString(params ...string) string {
	if value, ok := c.Custom[params[0]]; ok {
		return value.(string)
	}

	if len(params) == 2 {
		return params[1]
	}

	return ""
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
	return c.Type == "index" || c.Type == "prebuild"
}

func (c *Container) IsRemote() bool {
	return c.Type == "github" || c.Type == "remote"
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
		util.LogError(err)
	}

	c.retrieveInfoFromInspection(inspect)
}

func (c *Container) retrieveInfoFromInspection(inspect []byte) {
	var results []inspection
	goyaml.Unmarshal(inspect, &results)

	isRunning := false

	if len(results) > 0 {
		rawRunning := results[0].State["Running"]
		if rawRunning != nil {
			isRunning = rawRunning.(bool)
		}

		if isRunning {
			c.Ip = results[0].NetworkSettings["IPAddress"]
			c.Id = results[0].ID
		}
	}

	c.Running = isRunning
}
