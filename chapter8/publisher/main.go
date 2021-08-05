package main

import (
	"log"
	"runtime"
	"github.com/nats-io/nats.go"
)

func main() {
	// Create a connection to the NATS server over TLS with the RootCA
	nc, err := nats.Connect("tls://nats.example.com:4222", nats.Name("registration-service"),
							nats.RootCAs("../ca.pem"), nats.UserInfo("registration_service", "reg123"))
	if err != nil {
		log.Fatal(err)
	}
	//defer nc.Close()

	// Publish a message on "patient.profile" subject
	subj, msg := "patient.register", []byte("{\"name\":\"parakum\"}")
	nc.Publish(subj, msg)
	nc.Flush()
	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	} else {
		log.Printf("Published [%s] : '%s'\n", subj, msg)
	}
	runtime.Goexit()
}


