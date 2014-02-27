package containerCollection

import (
	"github.com/marmelab/gaudi/container"
	"github.com/marmelab/gaudi/util"
	"os"
)

type ContainerCollection map[string]*container.Container

func Merge(c1, c2 ContainerCollection) ContainerCollection {
	result := make(ContainerCollection)

	for name, currentContainer := range c1 {
		result[name] = currentContainer
	}
	for name, currentContainer := range c2 {
		result[name] = currentContainer
	}

	return result
}

func (collection ContainerCollection) Init(relativePath string) bool {
	hasGaudiManagedContainer := false

	// Fill name & dependencies
	for name, currentContainer := range collection {
		currentContainer.Name = name

		if currentContainer.IsGaudiManaged() {
			hasGaudiManagedContainer = true
			currentContainer.Image = "gaudi/" + name
		}

		for _, dependency := range currentContainer.Links {
			if depContainer, exists := collection[dependency]; exists {
				currentContainer.AddDependency(depContainer)
			} else {
				panic(name + " references a non existing application : " + dependency)
			}
		}

		// Add relative path to volumes
		for volumeHost, volumeContainer := range currentContainer.Volumes {
			// Relative volume host
			if string(volumeHost[0]) != "/" {
				delete(currentContainer.Volumes, volumeHost)
				volumeHost = relativePath + "/" + volumeHost

				currentContainer.Volumes[volumeHost] = volumeContainer
			}

			// Create directory if needed
			if !util.IsDir(volumeHost) {
				err := os.MkdirAll(volumeHost, 0755)
				if err != nil {
					panic(err)
				}
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

	return hasGaudiManagedContainer
}

func (collection ContainerCollection) Get(name string) *container.Container {
	return collection[name]
}

func (collection ContainerCollection) Start() {
	collection.CheckIfNotEmpty()
	collection.Clean()

	startChans := make(map[string]chan bool, len(collection))

	// Start all applications
	for name, currentContainer := range collection {
		startChans[name] = make(chan bool)

		go startOne(currentContainer, startChans)
	}

	// Waiting for all applications to start
	for name, _ := range collection {
		<-startChans[name]
	}
}

func (collection ContainerCollection) Stop() {
	collection.CheckIfNotEmpty()

	nbContainers := len(collection)
	killChans := make(chan bool, nbContainers)

	for _, currentContainer := range collection {
		go currentContainer.Kill(false, killChans)
	}

	// Waiting for all applications to stop
	waitForIt(killChans)
}

func (collection ContainerCollection) CheckIfNotEmpty() {
	// Check if there is at least a container
	if collection == nil || len(collection) == 0 {
		panic("Gaudi requires at least an application to be defined to start anything")
	}
}

func (collection ContainerCollection) Clean() {
	nbContainers := len(collection)
	cleanChans := make(chan bool, nbContainers)
	buildChans := make(chan bool, nbContainers)

	// Clean all applications
	for _, currentContainer := range collection {
		go currentContainer.Clean(cleanChans)
	}
	waitForIt(cleanChans)

	// Build all applications & binaries
	for _, currentContainer := range collection {
		if currentContainer.IsPreBuild() {
			go currentContainer.Pull(buildChans)
		} else {
			go currentContainer.Build(buildChans)
		}
	}
	waitForIt(buildChans)
}

func waitForIt(channels chan bool) {
	nbContainers := cap(channels)

	for i := 0; i < nbContainers; i++ {
		<-channels
	}
}

func startOne(currentContainer *container.Container, done map[string]chan bool) {
	// Waiting for dependencies to be started
	for _, dependency := range currentContainer.Dependencies {
		<-done[dependency.Name]
	}

	currentContainer.Start()

	close(done[currentContainer.Name])
}
