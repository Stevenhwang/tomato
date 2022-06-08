package master

import (
	"log"
	"time"

	"github.com/nats-io/nats-server/v2/server"
)

func Start() {
	opts := &server.Options{Username: "nats", Password: "123qwe"}
	ns, err := server.NewServer(opts)
	if err != nil {
		panic(err)
	}
	go ns.Start()
	// start monitoring subscriptions(all agents)
	copts := &server.ConnzOptions{SubscriptionsDetail: true}
	for {
		time.Sleep(5 * time.Second)
		connz, _ := ns.Connz(copts)
		subs := []string{}
		if len(connz.Conns) > 0 {
			for _, cz := range connz.Conns {
				if len(cz.SubsDetail) > 0 {
					subs = append(subs, cz.SubsDetail[0].Subject)
				}
			}
		}
		log.Println(subs)
	}

	// nc, err := nats.Connect("192.168.1.188", nats.UserInfo("nats", "123qwe"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Println("connect nats-server success")
	// defer nc.Close()

	// // Send the request
	// msg, err := nc.Request("agent1", []byte("fuck hello ,,,"), 5*time.Second)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// // Use the response
	// log.Printf("Reply: %s", msg.Data)

	// // Close the connection
	// nc.Close()
}
