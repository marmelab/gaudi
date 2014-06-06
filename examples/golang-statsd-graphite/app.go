package main

import (
	"fmt"
	"math/rand"
	"github.com/cactus/go-statsd-client/statsd"
	"os"
	"time"
)

func main() {
	client, err := statsd.New(os.Getenv("STATSD_PORT_8125_TCP_ADDR") + ":8125", "test-client")
	
	// handle any errors
	if err != nil {
		panic(err)
	}
	
	// make sure to clean up
	defer client.Close()
	
	for true {
		rand.Seed(time.Now().Unix())
		stat := int64(rand.Intn(10))
		
		// Send a stat
		err = client.Inc("stat1", stat, 1.0)
		// handle any errors
		if err != nil {
			panic(err)
		}

		fmt.Println("Stat", stat, "sent")
		
		time.Sleep(1 * time.Second)
	}
}
