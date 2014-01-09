package maestro

import (
	"bytes"
	"github.com/marmelab/arch-o-matic/container"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"path/filepath"
	"text/template"
)

type Maestro struct {
	Containers map[string]*container.Container
	listeners  map[string]func()
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

	// Fill name & dependencies
	for name := range maestro.Containers {
		currentContainer := maestro.Containers[name]
		currentContainer.Name = name

		for _, dependency := range currentContainer.Links {
			currentContainer.AddDependency(maestro.Containers[dependency])
		}

		// Add relative path to volumes
		for volumeHost, volumeContainer := range currentContainer.Volumes {
			if string(volumeHost[0]) != "/" {
				delete(currentContainer.Volumes, volumeHost)

				currentContainer.Volumes[relativePath+"/"+volumeHost] = volumeContainer
			}
		}
	}

	maestro.listeners = make(map[string]func(), 0)
}

func (maestro *Maestro) parseTemplates() {
	templateDir := os.Getenv("GOPATH") + "/src/github.com/marmelab/arch-o-matic/templates/"
	parsedTemplateDir := "/tmp/arch-o-matic/"

	err := os.MkdirAll(parsedTemplateDir, 0700)
	if err != nil {
		panic(err)
	}

	for _, currentContainer := range maestro.Containers {
		files, err := ioutil.ReadDir(templateDir + currentContainer.Type)
		if err != nil {
			panic(err)
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
			tmpl, err := template.New(filePath).Delims("[[", "]]").Parse(string(content))
			if err != nil {
				panic(err)
			}

			templateDate := TemplateData{maestro, currentContainer}
			var result bytes.Buffer
			err = tmpl.Execute(&result, templateDate)
			if err != nil {
				panic(err)
			}

			// Create new file
			ioutil.WriteFile(destination, []byte(result.String()), 0644)
		}
	}
}

func (maestro *Maestro) Start() {
	maestro.parseTemplates()

	cleanChans := make(chan bool, len(maestro.Containers))
	buildChans := make(chan bool, len(maestro.Containers))
	startChans := make(map[string]chan bool)

	// Clean all containers
	for _, currentContainer := range maestro.Containers {
		go currentContainer.Clean(cleanChans)
	}
	<-cleanChans

	// Build all containers
	for _, currentContainer := range maestro.Containers {
		go currentContainer.Build(buildChans)
	}
	<-buildChans

	// Start all containers
	for name, currentContainer := range maestro.Containers {
		startChans[name] = make(chan bool)

		go maestro.startContainer(currentContainer, startChans)
	}

	// Waiting for all containers to start
	for containerName, _ := range maestro.Containers {
		<-startChans[containerName]
	}
}

func (maestro *Maestro) startContainer(currentContainer *container.Container, done map[string]chan bool) {
	// Waiting for dependencies to start
	for _, dependency := range currentContainer.Dependencies {
		<-done[dependency.Name]
	}

	currentContainer.Start()

	close(done[currentContainer.Name])
}

func (maestro *Maestro) GetContainer(name string) *container.Container {
	return maestro.Containers[name]
}
