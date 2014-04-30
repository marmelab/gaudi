package main

import (
	"flag"
	"github.com/marmelab/gaudi/gaudi"
	"github.com/marmelab/gaudi/util"
	"os"
	"strings"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, " ")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)

	return nil
}

var (
	config = flag.String("config", ".gaudi.yml", "File describing the architecture")
)

func main() {
	flag.Parse()

	rebuild := len(flag.Args()) > 0 && flag.Args()[0] == "rebuild"
	g := gaudi.Gaudi{}
	g.InitFromFile(retrieveConfigPath(*config))

	if len(flag.Args()) == 0 || rebuild {
		// Start all applications
		g.StartApplications(rebuild)
	} else {
		switch os.Args[1] {
		case "run":
			// Run a specific command
			g.Run(os.Args[2], os.Args[3:])
			break
		case "stop":
			// Stop all applications
			g.StopApplications()
			break
		case "check":
			// Check if all applications are running
			g.Check()
			break
		case "clean":
			// Clean application containers
			g.Clean()
			break
		default:
			util.LogError("Argument " + os.Args[1] + " was not found")
			break
		}
	}
}

func retrieveConfigPath(configFile string) string {
	if len(configFile) == 0 {
		util.LogError("Config file name cannot be empty.")
	}

	if string(configFile[0]) != "/" {
		currentDir, err := os.Getwd()
		if err != nil {
			util.LogError(err)
		}

		configFile = currentDir + "/" + configFile
	}

	if !util.IsFile(configFile) {
		util.LogError("Config file must be a file.")
	}

	return configFile
}
