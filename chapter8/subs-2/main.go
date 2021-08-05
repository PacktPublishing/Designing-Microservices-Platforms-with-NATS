package main

import (
	"fmt"
	"log"
	"runtime"
	"github.com/nats-io/nats.go"
)

func main() {
	// Create a connection to the NATS server over TLS with the RootCA
	nc, err := nats.Connect("tls://nats.example.com:6222", nats.Name("inspection-service-2"),
							nats.RootCAs("../ca.pem"),nats.UserInfo("inspection_service", "ins123"))	
	if err != nil {
		log.Fatal(err)
	}
	//defer nc.Close()


	fmt.Println("Listening on [patient.register] subject")

	// Subscribe
	nc.Subscribe("patient.register", func(msg *nats.Msg) {
		fmt.Printf("Received on [%s]: '%s' \n", msg.Subject, string(msg.Data))
	})
	nc.Flush()

	if err := nc.LastError(); err != nil {
		log.Fatal(err)
	}
	runtime.Goexit()



	// // Use a WaitGroup to wait for a message to arrive
	// wg := sync.WaitGroup{}
	// wg.Add(1)
	// fmt.Println("Listening on [patient.register] subject")

	// // Subscribe
	// if _, err := nc.Subscribe("patient.register", func(m *nats.Msg) {
	// 	wg.Done()
	// 	fmt.Printf("Received on [%s]: '%s'", m.Subject, string(m.Data))
	// }); err != nil {
	// 	log.Fatal(err)
	// }

	// // Wait for a message to come in
	// wg.Wait()
}
