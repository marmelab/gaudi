package docker

import (
	"errors"
	"flag"
	"github.com/marmelab/gaudi/util"
	"os"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"
)

var (
	docker, _   = exec.LookPath("docker")
	dockerIo, _ = exec.LookPath("docker.io")
	noCache     = flag.Bool("no-cache", false, "Disable build cache")
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
	buildCmd.Stdin = os.Stdin

	out, err := buildCmd.CombinedOutput()
	if err != nil {
		util.Print(string(out))
		util.LogError("Error while starting container '" + name + "'")
	}

	buildCmd.Wait()

	time.Sleep(1 * time.Second)
}

func Pull(name string) {
	pullCmd := exec.Command(getDockerBinaryPath(), "pull", name)
	pullCmd.Stderr = os.Stderr
	pullCmd.Stdin = os.Stdin

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

	// Inspect container to check status code
	exitCodeBuff, _ := Inspect(name, "--format", "{{.State.ExitCode}}")
	exitCode, _ := strconv.Atoi(strings.TrimSpace(string(exitCodeBuff)))

	if exitCode != 0 {
		error, _ := Logs(name)
		util.LogError("Error while starting container '" + name + "' : " + error)
	}

	return string(out)
}

/**
 * Start a container as binary
 */
func Run(name, currentPath string, arguments []string, ports, environments map[string]string) {
	runFunc := reflect.ValueOf(exec.Command)
	rawArgs := []string{getDockerBinaryPath(), "run", "-v=" + currentPath + ":" + currentPath, "-w=" + currentPath}

	// Add environments
	util.Debug(environments)
	for envName, envValue := range environments {
		rawArgs = append(rawArgs, "-e="+envName+"="+envValue)
	}

	// Add ports
	for portIn, portOut := range ports {
		rawArgs = append(rawArgs, "-p="+string(portIn)+":"+string(portOut))
	}

	rawArgs = append(rawArgs, name)

	// Add user arguments
	for _, argument := range arguments {
		rawArgs = append(rawArgs, argument)
	}

	runCmd := runFunc.Call(util.BuildReflectArguments(rawArgs))[0].Interface().(*exec.Cmd)
	runCmd.Stdout = os.Stdout
	runCmd.Stdin = os.Stdin
	runCmd.Stderr = os.Stderr

	util.Debug("Run command:", runCmd.Args)

	if err := runCmd.Run(); err != nil {
		util.LogError(err)
	}
}

func Exec(args []string, hasStdout bool) {
	execFunc := reflect.ValueOf(exec.Command)
	execCmd := execFunc.Call(util.BuildReflectArguments(args))[0].Interface().(*exec.Cmd)

	if hasStdout {
		execCmd.Stdout = os.Stdout
	}

	execCmd.Stdin = os.Stdin
	execCmd.Stderr = os.Stderr

	util.Debug("Exec command:", execCmd.Args)

	if err := execCmd.Start(); err != nil {
		util.LogError(err)
	}
}

func Enter(name string) {
	var pid string
	var imageExists bool
	nsenter, _ := exec.LookPath("nsenter")

	ps, _ := SnapshotProcesses()
	if pid, imageExists = ps[name]; !imageExists {
		util.LogError("Image " + name + " doesn't exists")
	}

	statePidBuff, _ := Inspect(pid, "--format", "{{.State.Pid}}")
	statePid := strings.TrimSpace(string(statePidBuff))

	enterCmd := exec.Command("sudo", nsenter, "--target", statePid, "--mount", "--uts", "--ipc", "--net", "--pid")
	enterCmd.Stdout = os.Stdout
	enterCmd.Stdin = os.Stdin
	enterCmd.Stderr = os.Stderr

	util.Debug("Run command:", enterCmd.Args)

	if err := enterCmd.Run(); err != nil {
		util.LogError(err)
	}
}

func Inspect(params ...string) ([]byte, error) {
	inspectFunc := reflect.ValueOf(exec.Command)
	rawArgs := []string{getDockerBinaryPath(), "inspect"}

	// Add extra arguments
	if len(params) > 1 {
		for _, arg := range params[1:] {
			rawArgs = append(rawArgs, arg)
		}
	}

	// Add the id at the end
	rawArgs = append(rawArgs, params[0])

	inspectCmd := inspectFunc.Call(util.BuildReflectArguments(rawArgs))[0].Interface().(*exec.Cmd)
	out, err := inspectCmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func Logs(name string) (string, error) {
	logsCmd := exec.Command(getDockerBinaryPath(), "logs", name)

	out, err := logsCmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out), nil
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

func ShouldRebuild(imageName string) bool {
	images, err := GetImages()

	if err != nil {
		return false
	}

	_, ok := images[imageName]
	return !ok
}

func GetImages() (map[string]string, error) {
	images := make(map[string]string)

	imagesCommand := exec.Command(getDockerBinaryPath(), "images")
	out, err := imagesCommand.CombinedOutput()
	if err != nil {
		return nil, errors.New(string(out))
	}

	// Retrieve lines & remove first and last one
	lines := strings.Split(string(out), "\n")
	lines = lines[1 : len(lines)-1]

	for _, line := range lines {
		fields := strings.Fields(line)
		if fields[0] == "<none>" {
			continue
		}

		images[fields[0]] = fields[2]
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
