package main

import (
	"github.com/gocql/gocql"
	"fmt"
	"os"
)

func main () {
	var id gocql.UUID
	var text string

	cluster := gocql.NewCluster(os.Getenv("DB_PORT_9160_TCP_ADDR"))
	cluster.Keyspace = "myApp"
	cluster.Consistency = gocql.Quorum
	session, err := cluster.CreateSession()
	defer session.Close()

	if err != nil {
		panic(err)
	}

	if err := session.Query(`INSERT INTO tweet (timeline, id, text) VALUES (?, ?, ?)`, "me", gocql.TimeUUID(), "hello world").Exec(); err != nil {
		panic(err)
	}

	if err := session.Query(`SELECT id, text FROM tweet WHERE timeline = ? LIMIT 1`, "me").Consistency(gocql.One).Scan(&id, &text); err != nil {
		panic(err)
	}
	fmt.Println("Tweet:", id, text)
}
