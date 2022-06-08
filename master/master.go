package master

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func Start() {
	nc, err := nats.Connect("192.168.1.188", nats.UserInfo("nats", "123qwe"))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connect nats-server success")
	defer nc.Close()

	// Send the request
	msg, err := nc.Request("agent1", []byte("fuck hello ,,,"), 5*time.Second)
	if err != nil {
		log.Fatal(err)
	}

	// Use the response
	log.Printf("Reply: %s", msg.Data)

	// Close the connection
	nc.Close()
}
