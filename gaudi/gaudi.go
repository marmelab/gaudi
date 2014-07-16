package gaudi

import (
	"archive/tar"
	"bytes"
	"flag"
	"github.com/marmelab/gaudi/container"
	"github.com/marmelab/gaudi/containerCollection"
	"github.com/marmelab/gaudi/docker"
	"github.com/marmelab/gaudi/util"
	"io"
	"io/ioutil"
	"launchpad.net/goyaml"
	"net/http"
	"os"
	"path"
	"strings"
	"text/template"
)

var (
	emptyCmdFlag          = flag.String("empty-cmd", "", "Use /bin/bash for specified component (use all to target all components)")
	emptyCmdForContainers []string
)

func main() {
	flag.Parse()
}

const DEFAULT_BASE_IMAGE = "stackbrew/debian"
const DEFAULT_BASE_IMAGE_WITH_TAG = "stackbrew/debian:wheezy"
const TEMPLATE_DIR = "/var/tmp/gaudi/templates/"
const TEMPLATE_REMOTE_PATH = "http://gaudi.io/apt/templates.tar"
const PARSED_TEMPLATE_DIR = "/tmp/gaudi/"
const VERSION = "0.2.1"

type Gaudi struct {
	Applications      containerCollection.ContainerCollection
	Binaries          containerCollection.ContainerCollection
	All               containerCollection.ContainerCollection
	ApplicationDir    string
	ConfigurationPath string
}

type TemplateData struct {
	Collection containerCollection.ContainerCollection
	Container  *container.Container
	EmptyCmd   bool
}

func (gaudi *Gaudi) InitFromFile(file string) {
	gaudi.ConfigurationPath = file
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

	emptyCmdForContainers = strings.Split(*emptyCmdFlag, ",")

	// Init all containers
	gaudi.Applications.AddAmbassadors()
	gaudi.All = containerCollection.Merge(gaudi.Applications, gaudi.Binaries)
	if len(gaudi.All) == 0 {
		util.LogError("No application or binary to start. Are you missing a 'applications' or 'binaries' field in your configuration ?")
	}

	hasGaudiManagedContainer := gaudi.All.Init(gaudi.ApplicationDir)

	// Apply extends
	gaudi.applyInheritance()

	// Check if docker is installed
	if !docker.HasDocker() {
		util.LogError("Docker should be installed to use Gaudi (see: https://www.docker.io/gettingstarted/).")
	}

	// Check if base image is pulled
	if hasGaudiManagedContainer && !docker.ImageExists(DEFAULT_BASE_IMAGE) {
		util.PrintGreen("Pulling base image (this may take a few minutes) ...")

		docker.Pull(DEFAULT_BASE_IMAGE_WITH_TAG)
	}

	if gaudi.useNewVersion() {
//		os.RemoveAll(TEMPLATE_DIR)
	}

	// Check if templates are present
	if !util.IsDir(TEMPLATE_DIR) {
		util.PrintGreen("Retrieving templates ...")

		retrieveTemplates()
		extractTemplates()
	}

	gaudi.build()
}

func (gaudi *Gaudi) StartApplications(rebuild bool) {
	// Force rebuild if needed
	if rebuild == false {
		rebuild = gaudi.shouldRebuild()

		if rebuild {
			util.PrintOrange("Changes detected in configuration file, rebuilding containers ...")
		}
	}

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
 * Enter in a specific container
 */
func (gaudi *Gaudi) Enter(name string) {
	// Check if nsenter exists
	images, err := docker.GetImages()
	if err != nil {
		util.LogError(err)
	}

	if _, ok := images["jpetazzo/nsenter"]; !ok {
		// Pull ns-enter image
		util.PrintGreen("Retrieving ns-enter image ...")
		docker.Exec([]string{"run", "--rm", "-v", "/usr/local/bin:/target", "jpetazzo/nsenter"})
	}

	container := gaudi.All[name]
	docker.Enter(container.GetFullName())
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

			util.PrintOrange("Application", currentContainer.Name, "is running", "("+currentContainer.Ip+":"+currentContainer.GetFirstPort()+")")
		} else {
			util.PrintOrange("Application", currentContainer.Name, "is not running")
		}
	}
}

/**
 * Kill & Remove all containers
 */
func (gaudi *Gaudi) Clean() {
	// Clean all containers
	gaudi.Applications.Clean()
}

func (gaudi *Gaudi) GetApplication(name string) *container.Container {
	if application, ok := gaudi.Applications[name]; ok {
		return application
	}

	return nil
}

func (gaudi *Gaudi) build() {
	// Retrieve application Path
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
			util.LogError("Container " + currentContainer.Name + " should have a field called 'type'.")
		}

		if !currentContainer.IsGaudiManaged() {
			continue
		}

		// Check if the beforeScript is a file
		beforeScript := currentContainer.BeforeScript
		if len(beforeScript) != 0 {
			copied := gaudi.copyRelativeFiles(beforeScript, PARSED_TEMPLATE_DIR+currentContainer.Name+"/")
			if copied {
				currentContainer.BeforeScript = "./" + currentContainer.BeforeScript
			}
		}

		// Check if the afterScript is a file
		afterScript := currentContainer.AfterScript
		if len(afterScript) != 0 {
			copied := gaudi.copyRelativeFiles(afterScript, PARSED_TEMPLATE_DIR+currentContainer.Name+"/")
			if copied {
				currentContainer.AfterScript = "./" + currentContainer.AfterScript
			}
		}

		templateDir, isCustom := gaudi.GetContainerTemplate(currentContainer)

		files, err := ioutil.ReadDir(templateDir)
		if err != nil {
			util.LogError("Application '" + currentContainer.Type + "' is not supported. Check http://gaudi.io/components.html for a list of supported applications.")
		}

		err = os.MkdirAll(PARSED_TEMPLATE_DIR+currentContainer.Name, 0755)
		if err != nil {
			util.LogError(err)
		}

		sourceTemplateDir := TEMPLATE_DIR + currentContainer.Type
		if isCustom {
			sourceTemplateDir = templateDir
		}

		// Parse & copy files
		for _, file := range files {
			gaudi.parseTemplate(sourceTemplateDir, PARSED_TEMPLATE_DIR, file, includes, currentContainer)
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

func (gaudi *Gaudi) GetContainerTemplate(container *container.Container) (string, bool) {
	templatePath := TEMPLATE_DIR + container.Type
	isCustom := container.Type == "custom"

	// Check if the application has a custom template
	if isCustom {
		templatePath = container.Template

		// Allows relative path with ./
		if len(templatePath) > 1 && templatePath[0] == '.' && templatePath[1] == '/' {
			templatePath = gaudi.ApplicationDir + "/" + strings.Join(strings.Split(templatePath, "/")[1:], "/")
		}

		// Handle relative patch without ./
		if templatePath[0] != '/' {
			templatePath = gaudi.ApplicationDir + "/" + templatePath
		}

		// Template path should be a directory
		if util.IsFile(templatePath) {
			templateParts := strings.Split(templatePath, "/")
			templatePath = strings.Join(templateParts[0:len(templateParts)-1], "/")
		}
	}

	return templatePath, isCustom
}

func (gaudi *Gaudi) parseTemplate(sourceDir, destinationDir string, file os.FileInfo, includes map[string]string, currentContainer *container.Container) {
	emptyCmd := shouldEmptyCmdForContainer(currentContainer.Name)
	templateData := TemplateData{gaudi.All, nil, emptyCmd}
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
	filePath := sourceDir + "/" + file.Name()
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

func (g *Gaudi) copyRelativeFiles(filePath, destination string) bool {
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

func shouldEmptyCmdForContainer(containerName string) bool {
	if len(emptyCmdForContainers) == 1 && emptyCmdForContainers[0] == "all" {
		return true
	}

	for _, name := range emptyCmdForContainers {
		if name == containerName {
			return true
		}
	}

	return false
}

func retrieveTemplates() {
	os.MkdirAll(TEMPLATE_DIR, 0755)

	archive, err := os.Create(TEMPLATE_DIR + "templates.tar")
	if err != nil {
		util.LogError(err)
	}
	defer archive.Close()

	content, err := http.Get(TEMPLATE_REMOTE_PATH)
	if err != nil {
		util.LogError(err)
	}
	defer content.Body.Close()

	_, err = io.Copy(archive, content.Body)
	if err != nil {
		util.LogError(err)
	}
}

func extractTemplates() {
	tarFile, _ := os.Open(TEMPLATE_DIR + "templates.tar")
	defer tarFile.Close()

	tar := tar.NewReader(tarFile)

	for {
		header, err := tar.Next()
		if err == io.EOF {
			break
		}

		if err != nil {
			util.LogError(err)
		}

		// Remove first path part
		filePath := strings.Join(strings.Split(header.Name, "/")[1:], "/")

		// Check if we should create a folder or a file
		if header.Size == 0 {
			err := os.MkdirAll(TEMPLATE_DIR+filePath, 0755)
			if err != nil {
				util.LogError(err)
			}
		} else {
			f, err := os.Create(TEMPLATE_DIR + filePath)
			if err != nil {
				util.LogError(err)
			}
			defer f.Close()

			_, err = io.Copy(f, tar)
			if err != nil {
				util.LogError(err)
			}
		}
	}

	os.Remove(TEMPLATE_DIR + "templates.tar")
}

func (gaudi *Gaudi) shouldRebuild() bool {
	shouldRebuild := true
	checkSumFile := gaudi.ApplicationDir + "/.gaudi/.gaudi.sum"
	currentCheckSum := util.GetFileCheckSum(gaudi.ConfigurationPath)

	if util.IsFile(checkSumFile) {
		oldCheckSum, _ := ioutil.ReadFile(checkSumFile)

		shouldRebuild = string(oldCheckSum) != currentCheckSum
	}

	// Write new checksum
	ioutil.WriteFile(checkSumFile, []byte(currentCheckSum), 775)

	return shouldRebuild
}

func (gaudi *Gaudi) useNewVersion() bool {
	useNewVersion := true
	versionFile := gaudi.ApplicationDir + "/.gaudi/version.txt"

	if util.IsFile(versionFile) {
		oldVersion, _ := ioutil.ReadFile(versionFile)
		useNewVersion = string(oldVersion) != VERSION
	}

	// Write new version
	ioutil.WriteFile(versionFile, []byte(VERSION), 775)

	return useNewVersion
}

func (gaudi *Gaudi) applyInheritance() {
	extendsChan := make(map[string]chan bool)

	for name, currentContainer := range gaudi.Applications {
		extendsChan[name] = make(chan bool)

		if currentContainer.Extends == "" {
			close(extendsChan[name])
			continue
		}

		if parentContainer, exists := gaudi.All[currentContainer.Extends]; exists {
			go extendsOne(currentContainer, parentContainer, extendsChan)
		} else {
			util.LogError(currentContainer.Name + " extends a non existing application : " + currentContainer.Extends)
		}
	}

	// Waiting for all applications to be extended
	for name, _ := range extendsChan {
		<-extendsChan[name]
	}
}

func extendsOne(child, parent *container.Container, extendsChan map[string]chan bool) {
	// Waiting for dependencies to be extended
	<-extendsChan[parent.Name]

	// Extends the application
	child.ExtendsContainer(parent)
	close(extendsChan[child.Name])
}

func getIncludes() map[string]string {
	includesDir := TEMPLATE_DIR + "_includes/"
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
