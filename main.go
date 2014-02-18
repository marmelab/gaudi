package main

import (
	"flag"
	"github.com/marmelab/gaudi/maestro"
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
	runArgs stringSlice
	config  = flag.String("config", ".gaudi.yml", "File describing the architecture")
	run     = flag.String("run", "", "Run a container as a binary file")
	stop    = flag.Bool("stop", false, "Stop all applications ( data not stored in volumes will be lost)")
	check   = flag.Bool("check", false, "Check if all applications are running")
)

func main() {
	flag.Parse()

	m := maestro.Maestro{}
	m.InitFromFile(retrieveConfigPath(*config))

	if len(*run) > 0 {
		runArgs := strings.Split(*run, " ")

		// Run a specific command
		m.Run(runArgs[0], runArgs[1:])
	} else {
		if *check {
			m.Check()
		} else if *stop {
			m.Stop()
		} else {
			m.Start()
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
