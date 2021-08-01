package main

import (
	"log"
	"github.com/nats-io/nats.go"
)

func main() {
	// Create a connection to the NATS server over TLS with the RootCA
	nc, err := nats.Connect("tls://nats.example.com:4222,tls://nats.example.com:5222,tls://nats.example.com:6222", nats.RootCAs("../ca.pem"))
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()

	// Publish a message on "patient.profile" subject
	subj, msg := "patient.profile", []byte("{\"name\":\"parakum\"}")
	nc.Publish(subj, msg)
	nc.Flush()
	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", subj, msg)
	}
}


