package main

import (
	"flag"
	"github.com/marmelab/arch-o-matic/maestro"
	"os"
)

var (
	config 	= flag.String("config", ".arch-o-matic.yml", "File describing the architecture")
	rebuild = flag.Bool("rebuild", false, "Rebuild all containers ( data not stored in volumes will be lost)")
	check = flag.Bool("check", false, "Check if all containers are running")
)

func main() {
	m := maestro.Maestro{}
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	flag.Parse()

	configPath := *config
	if string(configPath[0]) != "/" {
		configPath = dir + "/" + configPath
	}

	m.InitFromFile(configPath)

	if *check {
		m.Start(*rebuild)
	} else {
		m.Check()
	}
}
