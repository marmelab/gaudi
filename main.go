package main

import (
	"flag"
	"github.com/marmelab/gaudi/maestro"
	"github.com/marmelab/gaudi/util"
	"os"
)

var (
	config 	= flag.String("config", ".gaudi.yml", "File describing the architecture")
	rebuild = flag.Bool("rebuild", false, "Rebuild all containers ( data not stored in volumes will be lost)")
	stop = flag.Bool("stop", false, "Stop all containers ( data not stored in volumes will be lost)")
	check = flag.Bool("check", false, "Check if all containers are running")
)

func main() {
	flag.Parse()

	m := maestro.Maestro{}
	m.InitFromFile(retrieveConfigPath(*config))

	if *check {
		m.Check()
	} else if *stop{
		m.Stop()
	} else {
		m.Start(*rebuild)
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
