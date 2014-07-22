package container

import (
	"github.com/marmelab/gaudi/docker"
	"github.com/marmelab/gaudi/util"
	"launchpad.net/goyaml"
	"reflect"
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
	Extends      string
	Image        string
	Path         string
	Template     string
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

	if c.Custom == nil {
		c.Custom = make(map[string]interface{})
	}
}

func (c *Container) ExtendsContainer(parent *Container) {
	reflectedParent := reflect.ValueOf(parent).Elem()
	reflectedContainer := reflect.ValueOf(c).Elem()
	typeOfT := reflectedContainer.Type()

	// Check if fields are empty, take the value of the parent if so
	for i := 0; i < reflectedContainer.NumField(); i++ {
		field := reflectedContainer.Field(i)
		fieldName := typeOfT.Field(i).Name
		fieldType := field.Type().Name()

		if fieldType == "bool" || fieldType == "Ambassador" {
			continue
		}

		if field.Len() == 0 {
			field.Set(reflectedParent.FieldByName(fieldName))
		}
	}
}

func (c *Container) Remove() {
	docker.Remove(c.Name)
	c.Running = false
}

func (c *Container) Kill(silent bool, done chan bool) {
	if !silent {
		util.PrintGreen("Killing", c.Name, "...")
	}

	docker.Kill(c.Name)
	c.Running = false

	if done != nil {
		done <- true
	}
}

func (c *Container) Clean(done chan bool) {
	util.PrintGreen("Cleaning", c.Name, "...")

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
			util.PrintRed("WARN: 'remote' type is deprecated, use 'github' instead")
		}

		buildName = c.Image
		buildPath = c.Path
	}

	util.PrintGreen("Building", buildName, "...")
	docker.Build(buildName, buildPath)

	done <- true
}

func (c *Container) Pull(done chan bool) {
	util.PrintGreen("Pulling", c.Image, "...")

	// prebuild type is deprecated
	if c.Type == "prebuild" {
		util.PrintRed("WARN: 'prebuild' type is deprecated, use 'index' instead")
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
	// Check if the container is already running
	if !rebuild {
		if c.IsRunning() {
			util.PrintGreen("Application", c.Name, "is already running", "("+c.Ip+":"+c.GetFirstPort()+")")
			return
		}

		cleanChan := make(chan bool, 1)
		go c.Clean(cleanChan)
		<-cleanChan
	}

	util.PrintGreen("Starting", c.Name, "...")

	startResult := docker.Start(c.Name, c.Image, c.Links, c.Ports, c.Volumes, c.Environments)
	c.Id = strings.TrimSpace(startResult)

	time.Sleep(3 * time.Second)
	c.RetrieveIp()
	c.Running = true

	util.PrintGreen("Application", c.Name, "started", "("+c.Ip+":"+c.GetFirstPort()+")")
}

func (c *Container) BuildAndRun(currentPath string, arguments []string) {
	if docker.ShouldRebuild(c.Image) {
		buildChans := make(chan bool, 1)
		go c.BuildOrPull(buildChans)
		<-buildChans
	}

	c.Run(currentPath, arguments)
}

/**
 * Starts a container as a binary file
 */
func (c *Container) Run(currentPath string, arguments []string) {
	util.PrintGreen("Running", c.Name, strings.Join(arguments, " "), "...")

	docker.Run(c.Image, currentPath, arguments, c.Ports, c.Environments)
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

func (c *Container) SetCustomValue(name, value string) string {
	c.Custom[name] = value

	return ""
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

func (c *Container) GetFirstLocalPort(args ...string) string {
	for _, localPort := range c.Ports {
		return localPort
	}

	if len(args) == 1 {
		return args[0]
	} else {
		return ""
	}
}

func (c *Container) FirstLinked() *Container {
	for _, dep := range c.Dependencies {
		return dep
	}

	return nil
}

func (c *Container) GetFirstMountedDir() string {
	for _, volume := range c.Volumes {
		return volume
	}

	return "/"
}

func (c *Container) DependsOf(otherComponentType string) bool {
	for _, dep := range c.Dependencies {
		if dep.Type == otherComponentType {
			return true
		}
	}

	return false
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

func (c *Container) GetFullName() string {
	name := "gaudi/" + c.Name
	if !c.IsGaudiManaged() {
		name = c.Image
	}

	return name
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
