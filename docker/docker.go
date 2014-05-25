package docker

import (
	"errors"
	"flag"
	"github.com/marmelab/gaudi/util"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"time"
)

var (
	docker, _   = exec.LookPath("docker")
	dockerIo, _ = exec.LookPath("docker.io")
	noCache     = flag.Bool("no-cache", false, "Disable build cache")
	quiet       = flag.Bool("quiet", false, "Do not display build output")
)

func main() {
	flag.Parse()
}

func HasDocker() bool {
	return len(docker) > 0 || len(dockerIo) > 0
}

func ImageExists(name string) bool {
	imagesCmd := exec.Command(getDockerBinaryPath(), "images", name)

	out, err := imagesCmd.CombinedOutput()
	if err != nil {
		return false
	}

	// Retrieve lines & remove first and last one
	lines := strings.Split(string(out), "\n")
	lines = lines[1 : len(lines)-1]

	return len(lines) > 0
}

func Remove(name string) {
	removeCmd := exec.Command(getDockerBinaryPath(), "rm", name)
	removeErr := removeCmd.Start()
	if removeErr != nil {
		util.LogError(removeErr)
	}

	time.Sleep(1 * time.Second)
}

func Kill(name string) {

	killCommand := exec.Command(getDockerBinaryPath(), "kill", name)
	killErr := killCommand.Start()
	if killErr != nil {
		util.LogError(killErr)
	}

	time.Sleep(1 * time.Second)
}

func Build(name, path string) {
	buildFunc := reflect.ValueOf(exec.Command)
	rawArgs := []string{getDockerBinaryPath(), "build"}

	if *noCache {
		rawArgs = append(rawArgs, "--no-cache")
	}

	rawArgs = append(rawArgs, "-t", name, path)

	util.Debug(rawArgs)

	buildCmd := buildFunc.Call(util.BuildReflectArguments(rawArgs))[0].Interface().(*exec.Cmd)
	buildCmd.Stderr = os.Stderr

	if !*quiet {
		buildCmd.Stdout = os.Stdout
	}

	if err := buildCmd.Run(); err != nil {
		util.LogError(err)
	}

	buildCmd.Wait()

	time.Sleep(1 * time.Second)
}

func Pull(name string) {
	pullCmd := exec.Command(getDockerBinaryPath(), "pull", name)
	pullCmd.Stderr = os.Stderr

	if !*quiet {
		pullCmd.Stdout = os.Stdout
	}

	util.Debug("Pull command:", pullCmd.Args)

	if err := pullCmd.Run(); err != nil {
		util.LogError(err)
	}

	pullCmd.Wait()
}

/**
 * Start a container as a server
 */
func Start(name, image string, links []string, ports, volumes, environments map[string]string) string {
	runFunc := reflect.ValueOf(exec.Command)
	rawArgs := []string{getDockerBinaryPath(), "run", "-d", "-i", "-t", "--privileged", "--name=" + name}

	// Add environments
	util.Debug(environments)
	for envName, envValue := range environments {
		rawArgs = append(rawArgs, "-e="+envName+"="+envValue)
	}

	// Add links
	for _, link := range links {
		rawArgs = append(rawArgs, "--link="+link+":"+link)
	}

	// Add ports
	for portIn, portOut := range ports {
		rawArgs = append(rawArgs, "-p="+string(portIn)+":"+string(portOut))
	}

	// Add volumes
	for volumeHost, volumeContainer := range volumes {
		rawArgs = append(rawArgs, "-v="+volumeHost+":"+volumeContainer)
	}

	rawArgs = append(rawArgs, image)

	// Initiate the command with several arguments
	runCmd := runFunc.Call(util.BuildReflectArguments(rawArgs))[0].Interface().(*exec.Cmd)
	util.Debug("Start command:", runCmd.Args)

	out, err := runCmd.CombinedOutput()
	if err != nil {
		util.LogError(string(out))
	}

	return string(out)
}

/**
 * Start a container as binary
 */
func Run(name, currentPath string, arguments []string) {
	runFunc := reflect.ValueOf(exec.Command)
	rawArgs := []string{getDockerBinaryPath(), "run", "-v=" + currentPath + ":" + currentPath, "-w=" + currentPath, name}

	for _, argument := range arguments {
		rawArgs = append(rawArgs, argument)
	}

	runCmd := runFunc.Call(util.BuildReflectArguments(rawArgs))[0].Interface().(*exec.Cmd)
	runCmd.Stdout = os.Stdout
	runCmd.Stdin = os.Stdin
	runCmd.Stderr = os.Stderr

	util.Debug("Run command:", runCmd.Args)

	if err := runCmd.Start(); err != nil {
		util.LogError(err)
	}
}

func Inspect(id string) ([]byte, error) {
	inspectCmd := exec.Command(getDockerBinaryPath(), "inspect", id)

	out, err := inspectCmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func SnapshotProcesses() (map[string]string, error) {
	images := make(map[string]string)

	psCommand := exec.Command(getDockerBinaryPath(), "ps")
	out, err := psCommand.CombinedOutput()
	if err != nil {
		return nil, errors.New(string(out))
	}

	// Retrieve lines & remove first and last one
	lines := strings.Split(string(out), "\n")
	lines = lines[1 : len(lines)-1]

	for _, line := range lines {
		fields := strings.Fields(line)
		nameParts := strings.Split(fields[1], ":")

		images[nameParts[0]] = fields[0]
	}

	return images, nil
}

func getDockerBinaryPath() string {
	if len(docker) != 0 {
		return docker
	} else {
		return dockerIo
	}
}
