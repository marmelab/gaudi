package gaudi

import (
	"bytes"
	"fmt"
	"github.com/marmelab/gaudi/container"
	"github.com/marmelab/gaudi/containerCollection"
	"github.com/marmelab/gaudi/docker"
	"github.com/marmelab/gaudi/util"
	"io/ioutil"
	"launchpad.net/goyaml"
	"os"
	"path"
	"runtime"
	"strings"
	"text/template"
)

const DEFAULT_BASE_IMAGE = "stackbrew/debian"
const DEFAULT_BASE_IMAGE_WITH_TAG = "stackbrew/debian:wheezy"
const PARSED_TEMPLATE_DIR = "/tmp/gaudi/"

type Gaudi struct {
	Applications   containerCollection.ContainerCollection
	Binaries       containerCollection.ContainerCollection
	All            containerCollection.ContainerCollection
	ApplicationDir string
}

type TemplateData struct {
	Collection containerCollection.ContainerCollection
	Container  *container.Container
}

func (gaudi *Gaudi) InitFromFile(file string) {
	gaudi.ApplicationDir = path.Dir(file)

	fileContent, err := ioutil.ReadFile(file)
	if err != nil {
		util.LogError(err)
	}

	gaudi.Init(string(fileContent))
}

func (gaudi *Gaudi) Init(content string) {
	err := goyaml.Unmarshal([]byte(content), &gaudi)
	if err != nil {
		util.LogError(err)
	}

	// Init all containers
	gaudi.All = containerCollection.Merge(gaudi.Applications, gaudi.Binaries)
	if len(gaudi.All) == 0 {
		util.LogError("No application or binary to start")
	}

	hasGaudiManagedContainer := gaudi.All.Init(gaudi.ApplicationDir)

	// Check if base image is pulled
	if hasGaudiManagedContainer && !docker.ImageExists(DEFAULT_BASE_IMAGE) {
		fmt.Println("Pulling base image (this may take a few minutes)  ...")

		docker.Pull(DEFAULT_BASE_IMAGE_WITH_TAG)
	}

	gaudi.build()
}

func (gaudi *Gaudi) StartApplications(rebuild bool) {
	gaudi.Applications.Start(rebuild)
}

func (gaudi *Gaudi) StopApplications() {
	gaudi.Applications.Stop()
}

/**
 * Runs a container as a binary
 */
func (gaudi *Gaudi) Run(name string, arguments []string) {
	gaudi.Binaries[name].BuildAndRun(gaudi.ApplicationDir, arguments)
}

/**
 * Check if all applications are started
 */
func (gaudi *Gaudi) Check() {
	images, err := docker.SnapshotProcesses()
	if err != nil {
		util.LogError(err)
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

func (gaudi *Gaudi) GetApplication(name string) *container.Container {
	if application, ok := gaudi.Applications[name]; ok {
		return application
	}

	return nil
}

func (gaudi *Gaudi) build() {
	// Retrieve application Path
	templateDir := getGaudiDirectory() + "/templates/"
	currentDirctory, _ := os.Getwd()

	err := os.MkdirAll(PARSED_TEMPLATE_DIR, 0700)
	if err != nil {
		util.LogError(err)
	}

	// Retrieve includes
	includes := getIncludes()

	for _, currentContainer := range gaudi.All {
		// Check if the container has a type
		if currentContainer.Type == "" {
			util.LogError("Container " + currentContainer.Name + " should have a type")
		}

		if !currentContainer.IsGaudiManaged() {
			continue
		}

		// Check if the beforeScript is a file
		beforeScript := currentContainer.BeforeScript
		if len(beforeScript) != 0 {
			copied := gaudi.copyRelativeFile(beforeScript, PARSED_TEMPLATE_DIR+currentContainer.Name+"/")
			if copied {
				currentContainer.BeforeScript = "./" + currentContainer.BeforeScript
			}
		}

		// Check if the afterScript is a file
		afterScript := currentContainer.AfterScript
		if len(afterScript) != 0 {
			copied := gaudi.copyRelativeFile(afterScript, PARSED_TEMPLATE_DIR+currentContainer.Name+"/")
			if copied {
				currentContainer.AfterScript = "./" + currentContainer.AfterScript
			}
		}

		files, err := ioutil.ReadDir(templateDir + currentContainer.Type)
		if err != nil {
			util.LogError("Template not found for application : " + currentContainer.Type)
		}

		err = os.MkdirAll(PARSED_TEMPLATE_DIR+currentContainer.Name, 0755)
		if err != nil {
			util.LogError(err)
		}

		// Parse & copy files
		for _, file := range files {
			gaudi.parseFile(templateDir, PARSED_TEMPLATE_DIR, file, includes, currentContainer)
		}

		// Copy all files marked as Add
		for fileToAdd := range currentContainer.Add {
			filePath := currentDirctory + "/" + fileToAdd

			directories := strings.Split(fileToAdd, "/")
			if len(directories) > 1 {
				os.MkdirAll(PARSED_TEMPLATE_DIR+currentContainer.Name+"/"+strings.Join(directories[0:len(directories)-1], "/"), 0755)
			}

			err := util.Copy(PARSED_TEMPLATE_DIR+currentContainer.Name+"/"+fileToAdd, filePath)
			if err != nil {
				util.LogError(err)
			}
		}
	}
}

func (gaudi *Gaudi) parseFile(sourceDir, destinationDir string, file os.FileInfo, includes map[string]string, currentContainer *container.Container) {
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
			util.LogError(err)
		}

		return
	}

	// Read the template
	filePath := sourceDir + currentContainer.Type + "/" + file.Name()
	rawContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		util.LogError(err)
	}
	content := string(rawContent)

	// Add includes
	for name, include := range includes {
		content = strings.Replace(content, "[[ "+name+" ]]", include, -1)
	}

	// Parse it
	// We need to change default delimiters because sometimes we have to parse values like ${{{ .Val }}} which cause an error
	tmpl, templErr := template.New(filePath).Funcs(funcMap).Delims("[[", "]]").Parse(content)
	if templErr != nil {
		util.LogError(templErr)
	}

	templateData.Container = currentContainer
	var result bytes.Buffer
	err = tmpl.Execute(&result, templateData)
	if err != nil {
		util.LogError(err)
	}

	// Create the destination file
	ioutil.WriteFile(destination, []byte(result.String()), 0644)
}

func (g *Gaudi) copyRelativeFile(filePath, destination string) bool {
	// File cannot be absolute
	if util.IsFile(filePath) && filePath[0] == '/' {
		util.LogError("File '" + filePath + "' cannot be an absolute path")
	}

	// Check if the relative file exists
	absolutePath := g.ApplicationDir + "/" + filePath
	if util.IsFile(absolutePath) {

		// Move file to the build context (and keep the same file tree)
		directories := strings.Split(filePath, "/")
		if len(directories) > 1 {
			os.MkdirAll(destination+strings.Join(directories[0:len(directories)-1], "/"), 0755)
		}

		err := util.Copy(destination+filePath, absolutePath)
		if err != nil {
			util.LogError(err)
		}

		return true
	}

	return false
}

func getGaudiDirectory() string {
	// withmock copy only test and tested file, so we need to retrieve the template from the real app path in test env
	testPath := os.Getenv("ORIG_GOPATH")
	if len(testPath) > 0 {
		return testPath + "/src/github.com/marmelab/gaudi/"
	}

	_, currentFile, _, _ := runtime.Caller(0)
	pathParts := strings.Split(currentFile, "/")

	return strings.Join(pathParts[0:len(pathParts)-2], "/")
}

func getIncludes() map[string]string {
	includesDir := getGaudiDirectory() + "/templates/_includes/"
	result := make(map[string]string)

	files, err := ioutil.ReadDir(includesDir)
	if err != nil {
		util.LogError(err)
	}

	for _, file := range files {
		name := strings.Split(file.Name(), ".")[0]
		content, _ := ioutil.ReadFile(includesDir + file.Name())

		result[name] = string(content)
	}

	return result
}
