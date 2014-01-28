package docker

import (
	"os/exec"
	"reflect"
	"time"
)

var docker, _ = exec.LookPath("docker")

func Remove(name string) {
	removeCmd := exec.Command(docker, "rm", name)
	removeErr := removeCmd.Start()
	if removeErr != nil {
		panic(removeErr)
	}
	time.Sleep(1 * time.Second)
}

func Kill(name string) {
	killCommand := exec.Command(docker, "kill", name)
	killErr := killCommand.Start()
	if killErr != nil {
		panic(killErr)
	}

	time.Sleep(1 * time.Second)
}

func Build(name, path string) {
	var buildCmd *exec.Cmd

	buildCmd = exec.Command(docker, "build", "-t", name, path)

	out, err := buildCmd.CombinedOutput()
	if err != nil {
		panic(out)
	}
}

func Pull(name string) {
	pullCmd := exec.Command(docker, "pull", name)

	out, err := pullCmd.CombinedOutput()
	if err != nil {
		panic(string(out))
	}
}

func Start(name, image string, links []string, ports, volumes map[string]string) string {
	runFunc := reflect.ValueOf(exec.Command)
	rawArgs := []string{docker, "run", "-d", "-i", "-t", "-name", name}

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
	runCmd := runFunc.Call(buildArguments(rawArgs))[0].Interface().(*exec.Cmd)
	out, err := runCmd.CombinedOutput()
	if err != nil {
		panic(string(out))
	}

	return string(out)
}

func Inspect(id string) ([]byte, error) {
	inspectCmd := exec.Command(docker, "inspect", id)

	out, err := inspectCmd.CombinedOutput()
	if err != nil {
		return nil, err
	}

	return out, nil
}

func buildArguments(rawArgs []string) []reflect.Value {
	args := make([]reflect.Value, 0)

	for _, arg := range rawArgs {
		args = append(args, reflect.ValueOf(arg))
	}

	return args
}
