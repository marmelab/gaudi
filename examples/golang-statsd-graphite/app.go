package main

import (
	"fmt"
	"os"
	"github.com/cactus/go-statsd-client/statsd"
)

func main() {
	client, err := statsd.New(os.Getenv("STATSD_PORT_8125_TCP_ADDR") + ":8125", "test-client")
	
	// handle any errors
	if err != nil {
		panic(err)
	}
	
	// make sure to clean up
	defer client.Close()

	// Send a stat
	err = client.Inc("stat1", 42, 1.0)
	// handle any errors
	
	if err != nil {
		panic(err)
	}
	
	fmt.Println("Stat sent")
}
