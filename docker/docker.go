package docker

import (
	"os/exec"
	"reflect"
	"time"
	"fmt"
)

var docker, _ = exec.LookPath("docker")

func Remove (name string) {
	removeCmd := exec.Command(docker, "rm", name)
	removeErr := removeCmd.Start()
	if removeErr != nil {
		panic (removeErr)
	}
	time.Sleep(1 * time.Second)
}

func Kill (name string) {
	killCommand := exec.Command(docker, "kill", name)
	killErr := killCommand.Start()
	if killErr != nil {
		panic (killErr)
	}
	time.Sleep(1 * time.Second)
}

func Build (name string) {
	buildCmd := exec.Command(docker, "build", "-rm", "-t", "gaudi/"+name, "/tmp/gaudi/"+name)

	out, err := buildCmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}
}

func Start (name string, links []string, ports, volumes map[string]string) string {
	runFunc := reflect.ValueOf(exec.Command)
	rawArgs := []string{docker, "run", "-d", "-name", name}

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

	rawArgs = append(rawArgs, "gaudi/"+name)

	// Initiate the command with several arguments
	runCmd := runFunc.Call(buildArguments(rawArgs))[0].Interface().(*exec.Cmd)

	out, err := runCmd.CombinedOutput()
	if err != nil {
		fmt.Println(string(out))
		panic(err)
	}

	return string(out)
}

func Inspect (id string) ([]byte, error) {
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
