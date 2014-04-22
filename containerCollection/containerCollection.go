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
				util.LogError(name + " references a non existing application : " + dependency)
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
					util.LogError(err)
				}
			}
		}
	}

	return hasGaudiManagedContainer
}

func (collection ContainerCollection) Get(name string) *container.Container {
	return collection[name]
}

func (collection ContainerCollection) Start(rebuild bool) {
	collection.CheckIfNotEmpty()

	if rebuild {
		collection.Clean()
		collection.Build()
	}

	startChans := make(map[string]chan bool, len(collection))

	// Start all applications
	for name, currentContainer := range collection {
		startChans[name] = make(chan bool)

		go startOne(currentContainer, rebuild, startChans)
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
		util.LogError("Gaudi requires at least an application to be defined to start anything")
	}
}

func (collection ContainerCollection) Clean() {
	nbContainers := len(collection)
	cleanChans := make(chan bool, nbContainers)

	// Clean all applications
	for _, currentContainer := range collection {
		go currentContainer.Clean(cleanChans)
	}
	waitForIt(cleanChans)
}

func (collection ContainerCollection) Build() {
	nbContainers := len(collection)
	buildChans := make(chan bool, nbContainers)

	// Build all
	for _, currentContainer := range collection {
		go currentContainer.BuildOrPull(buildChans)
	}
	waitForIt(buildChans)
}

func waitForIt(channels chan bool) {
	nbContainers := cap(channels)

	for i := 0; i < nbContainers; i++ {
		<-channels
	}
}

func (collection ContainerCollection) AddAmbassasor() {
	for name, currentContainer := range collection {
		if currentContainer.Ambassador.Type == "" {
			continue
		}

		// Add the ambassador
		ambassadorName := "ambassasor-" + name
		ambassador := &container.Container{Name: ambassadorName, Type: "ambassador"}
		ambassador.Init()

		ambassador.Links = append(ambassador.Links, name)
		ambassador.Ports[currentContainer.Ambassador.Port] = currentContainer.Ambassador.Port

		collection[ambassadorName] = ambassador
	}
}

func startOne(currentContainer *container.Container, rebuild bool, done map[string]chan bool) {
	// Waiting for dependencies to be started
	for _, dependency := range currentContainer.Dependencies {
		<-done[dependency.Name]
	}

	currentContainer.Start(rebuild)

	close(done[currentContainer.Name])
}
