package main

import (
	"arch-o-matic/maestro"
	"os"
	"flag"
)

/*
docker ps | grep 'ago' | awk '{print $1}' | xargs docker kill || docker ps -a | grep 'ago' | awk '{print $1}' | xargs docker rm
 */
func main(){
	m := maestro.Maestro{}
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	file := flag.String("config", "", "File describing the architecture")
	flag.Parse()

	if len(*file) > 0 {
		filePath := *file

		if string((*file)[0]) != "/" {
			filePath = dir + "/" + *file
		}

		m.InitFromFile(filePath)
		m.Start()
	}
}
