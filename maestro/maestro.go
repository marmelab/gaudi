package maestro

import (
	"launchpad.net/goyaml"
	"arch-o-matic/container"
	"io/ioutil"
	"path/filepath"
)

type Maestro struct {
	Containers map[string] *container.Container
	listeners map[string]func()
}

func (m *Maestro) InitFromFile(file string) {
	content, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	m.InitFromString(string(content), filepath.Dir(file))
}

func (maestro *Maestro) InitFromString (content, relativePath string) {
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

				currentContainer.Volumes[relativePath + "/" + volumeHost] = volumeContainer
			}
		}
	}

	maestro.listeners = make(map[string]func(), 0)
}

func (maestro *Maestro) Start() {
	buildChans := make(chan bool, len(maestro.Containers))
	startChans := make(map[string] chan bool)

	// Build all containers
	for _, currentContainer := range maestro.Containers {
		go currentContainer.Build(buildChans)
	}
	<- buildChans


	// Start all containers
	for name, currentContainer := range maestro.Containers {
		startChans[name] = make(chan bool)

		go maestro.startContainer(currentContainer, startChans)
	}

	// Waiting for all containers to start
	for containerName, _ := range maestro.Containers {
		<- startChans[containerName]
	}
}

func (maestro *Maestro) startContainer (currentContainer *container.Container, done map[string] chan bool) {
	// Waiting for dependencies to start
	for _, dependency := range currentContainer.Dependencies {
		<-done[dependency.Name]
	}

	currentContainer.Start()

	close(done[currentContainer.Name])
}
