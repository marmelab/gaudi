package gaudi

import (
	"bytes"
	"fmt"
	"github.com/marmelab/gaudi/container"
	"github.com/marmelab/gaudi/containerCollection"
	"github.com/marmelab/gaudi/docker"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"text/template"
)

const DEFAULT_BASE_IMAGE = "stackbrew/debian"
const DEFAULT_BASE_IMAGE_WITH_TAG = "stackbrew/debian:wheezy"

type Gaudi struct {
	Applications containerCollection.ContainerCollection
	Binaries     containerCollection.ContainerCollection
	All          containerCollection.ContainerCollection
	Path         string
}

type TemplateData struct {
	Collection containerCollection.ContainerCollection
	Container  *container.Container
}

func (gaudi *Gaudi) Init(file string) {
	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	content := string(fileContent)
	gaudi.Path = filepath.Dir(file)

	err = goyaml.Unmarshal([]byte(content), &gaudi)
	if err != nil {
		panic(err)
	}

	// Init all containers
	gaudi.All = containerCollection.Merge(gaudi.Applications, gaudi.Binaries)
	hasGaudiManagedContainer := gaudi.All.Init(gaudi.Path)

	// Check if base image is pulled
	if hasGaudiManagedContainer && !docker.ImageExists(DEFAULT_BASE_IMAGE) {
		fmt.Println("Pulling base image (this may take a few minutes)  ...")

		docker.Pull(DEFAULT_BASE_IMAGE_WITH_TAG)
	}

	gaudi.build()
}

func (gaudi *Gaudi) StartApplications() {
	gaudi.Applications.Start()
}

func (gaudi *Gaudi) StopApplications() {
	gaudi.Applications.Stop()
}

/**
 * Runs a container as a binary
 */
func (gaudi *Gaudi) Run(name string, arguments []string) {
	gaudi.Binaries[name].Run(gaudi.Path, arguments)
}

/**
 * Check if all applications are started
 */
func (gaudi *Gaudi) Check() {
	images, err := docker.SnapshotProcesses()
	if err != nil {
		panic(err)
	}

	for _, currentContainer := range gaudi.Applications {
		if containerId, ok := images[currentContainer.Image]; ok {
			currentContainer.Id = containerId
			currentContainer.RetrieveIp()

			fmt.Println("Application", currentContainer.Name, "is running", "("+currentContainer.Ip+":"+currentContainer.GetFirstPort()+")")
		} else {
			fmt.Println("Application", currentContainer.Name, "is not running")
		}
	}
}

func (gaudi *Gaudi) build() {
	// Retrieve application Path
	templateDir := getApplicationDir() + "/templates/"
	parsedTemplateDir := "/tmp/gaudi/"

	err := os.MkdirAll(parsedTemplateDir, 0700)
	if err != nil {
		panic(err)
	}

	for _, currentContainer := range gaudi.All {
		if !currentContainer.IsGaudiManaged() {
			continue
		}

		files, err := ioutil.ReadDir(templateDir + currentContainer.Type)
		if err != nil {
			panic("Template not found for application : " + currentContainer.Type)
		}

		err = os.MkdirAll(parsedTemplateDir+currentContainer.Name, 0755)
		if err != nil {
			panic(err)
		}

		// Parse & copy files
		for _, file := range files {
			gaudi.parseFile(templateDir, parsedTemplateDir, file, currentContainer)
		}
	}
}

func (gaudi *Gaudi) parseFile(sourceDir, destinationDir string, file os.FileInfo, currentContainer *container.Container) {
	templateData := TemplateData{gaudi.All, nil}
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
	}

	// Create destination directory if needed
	destination := destinationDir + currentContainer.Name + "/" + file.Name()
	if file.IsDir() {
		err := os.MkdirAll(destination, 0755)
		if err != nil {
			panic(err)
		}

		return
	}

	// Read the template
	filePath := sourceDir + currentContainer.Type + "/" + file.Name()
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		panic(err)
	}

	// Parse it : we need to change default delimiters because sometimes we have to parse values like ${{{ .Val }}}
	// which cause an error
	tmpl, err := template.New(filePath).Funcs(funcMap).Delims("[[", "]]").Parse(string(content))
	if err != nil {
		panic(err)
	}

	templateData.Container = currentContainer
	var result bytes.Buffer
	err = tmpl.Execute(&result, templateData)
	if err != nil {
		panic(err)
	}

	// Create the destination file
	ioutil.WriteFile(destination, []byte(result.String()), 0644)
}

func getApplicationDir() string {
	_, currentFile, _, _ := runtime.Caller(0)
	pathParts := strings.Split(currentFile, "/")

	return strings.Join(pathParts[0:len(pathParts)-2], "/")
}
