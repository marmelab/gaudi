package maestro

import (
	"bytes"
	"fmt"
	"github.com/marmelab/gaudi/container"
	"github.com/marmelab/gaudi/docker"
	"github.com/marmelab/gaudi/util"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type Maestro struct {
	Applications map[string]*container.Container
}

type TemplateData struct {
	Maestro   *Maestro
	Container *container.Container
}

func (m *Maestro) InitFromFile(file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	m.InitFromString(string(content), filepath.Dir(file))
}

func (maestro *Maestro) InitFromString(content, relativePath string) {
	err := goyaml.Unmarshal([]byte(content), &maestro)
	if err != nil {
		panic(err)
	}
	if maestro.Applications == nil {
		panic("No application to start")
	}

	// Fill name & dependencies
	for name := range maestro.Applications {
		currentContainer := maestro.Applications[name]
		currentContainer.Name = name

		if currentContainer.IsGaudiManaged() {
			currentContainer.Image = "gaudi/" + name
		}

		for _, dependency := range currentContainer.Links {
			if depContainer, exists := maestro.Applications[dependency]; exists {
				currentContainer.AddDependency(depContainer)
			} else {
				panic(name + " references a non existing application : " + dependency)
			}
		}

		// Add relative path to volumes
		for volumeHost, volumeContainer := range currentContainer.Volumes {
			if string(volumeHost[0]) != "/" {
				delete(currentContainer.Volumes, volumeHost)

				if !util.IsDir(relativePath + "/" + volumeHost) {
					panic(relativePath + "/" + volumeHost + " should be a directory")
				}

				currentContainer.Volumes[relativePath+"/"+volumeHost] = volumeContainer
			} else if !util.IsDir(volumeHost) {
				panic(volumeHost + " should be a directory")
			}
		}

		// Check if the beforeScript is a file
		beforeScript := currentContainer.BeforeScript
		if len(beforeScript) != 0 {
			if util.IsFile(beforeScript) {
				currentContainer.BeforeScript = beforeScript
			} else if util.IsFile(relativePath + "/" + beforeScript) {
				currentContainer.BeforeScript = relativePath + "/" + beforeScript
			}
		}
	}
}

func (maestro *Maestro) createHiddenDir() {
	currentDir, _ := os.Getwd()
	err := os.MkdirAll(currentDir+"/.gaudi", 0755)
	if err != nil {
		panic(err)
	}
}

func (maestro *Maestro) parseTemplates() {
	// Running withmock doesn't include templates files in withmock's temporary dir
	path := os.Getenv("GOPATH")
	testPath := os.Getenv("ORIG_GOPATH")
	if len(testPath) > 0 {
		path = testPath
	}

	templateDir := path + "/src/github.com/marmelab/gaudi/templates/"
	parsedTemplateDir := "/tmp/gaudi/"
	templateData := TemplateData{maestro, nil}
	funcMap := template.FuncMap{
		"ToUpper": strings.ToUpper,
		"ToLower": strings.ToLower,
	}

	err := os.MkdirAll(parsedTemplateDir, 0700)
	if err != nil {
		panic(err)
	}

	for _, currentContainer := range maestro.Applications {
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
			destination := parsedTemplateDir + currentContainer.Name + "/" + file.Name()
			if file.IsDir() {
				err := os.MkdirAll(destination, 0755)
				if err != nil {
					panic(err)
				}

				continue
			}

			// Read the template
			filePath := templateDir + currentContainer.Type + "/" + file.Name()
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				panic(err)
			}

			// Parse it (we need to change default delimiters because sometimes we have to parse values like ${{{ .Val }}}
			// which cause an error)
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

			// Create new file
			ioutil.WriteFile(destination, []byte(result.String()), 0644)
		}
	}
}

func (maestro *Maestro) Start() {
	maestro.createHiddenDir()
	maestro.parseTemplates()

	nbApplications := len(maestro.Applications)
	cleanChans := make(chan bool, nbApplications)
	// Clean all applications
	for _, currentContainer := range maestro.Applications {
		go currentContainer.Clean(cleanChans)
	}
	<-cleanChans

	buildChans := make(chan bool, len(maestro.Applications))

	// Build all applications
	for _, currentContainer := range maestro.Applications {
		if currentContainer.IsPreBuild() {
			go currentContainer.Pull(buildChans)
		} else {
			go currentContainer.Build(buildChans)
		}
	}

	for i := 0; i < nbApplications; i++ {
		<-buildChans
	}

	startChans := make(map[string]chan bool)

	// Start all applications
	for name, currentContainer := range maestro.Applications {
		startChans[name] = make(chan bool)

		go maestro.startContainer(currentContainer, startChans)
	}

	// Waiting for all applications to start
	for containerName, _ := range maestro.Applications {
		<-startChans[containerName]
	}
}

func (maestro *Maestro) GetContainer(name string) *container.Container {
	return maestro.Applications[name]
}

func (maestro *Maestro) Check() {
	images, err := docker.SnapshotProcesses()
	if err != nil {
		panic(err)
	}

	for _, currentContainer := range maestro.Applications {
		if containerId, ok := images[currentContainer.Image]; ok {
			currentContainer.Id = containerId
			currentContainer.RetrieveIp()

			fmt.Println("Application", currentContainer.Name, "is running", "("+currentContainer.Ip+":"+currentContainer.GetFirstPort()+")")
		} else {
			fmt.Println("Application", currentContainer.Name, "is not running")
		}
	}
}

func (maestro *Maestro) Stop() {
	killChans := make(chan bool, len(maestro.Applications))

	for _, currentContainer := range maestro.Applications {
		go currentContainer.Kill(killChans, false)
	}

	<-killChans
}

func (maestro *Maestro) startContainer(currentContainer *container.Container, done map[string]chan bool) {
	// Waiting for dependencies to start
	for _, dependency := range currentContainer.Dependencies {
		<-done[dependency.Name]
	}

	currentContainer.Start()

	close(done[currentContainer.Name])
}
