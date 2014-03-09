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

	g := gaudi.Gaudi{}
	g.InitFromFile(retrieveConfigPath(*config))

	if len(os.Args) == 1 {
		// Start all applications
		g.StartApplications()
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
		}
	}
}

func retrieveConfigPath(configFile string) string {
	if len(configFile) == 0 {
		panic("Config file name cannot be empty.")
	}

	if string(configFile[0]) != "/" {
		currentDir, err := os.Getwd()
		if err != nil {
			panic(err)
		}

		configFile = currentDir + "/" + configFile
	}

	if !util.IsFile(configFile) {
		panic("Config file must be a file.")
	}

	return configFile
}
