package main

import (
	"log"

	"github.com/nats-io/nats.go"
)

func main() {

	nc, err := nats.Connect(nats.DefaultURL)

	if err != nil {

		log.Fatal(err)

	}

	defer nc.Close()

	// Publish a message on "hello" subject
	subj, msg := "hello", []byte("world")
	nc.Publish(subj, msg)

	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", subj, msg)
	}

}
