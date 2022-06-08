package agent

import (
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/nats-io/nats.go"
)

func Start() {
	nc, err := nats.Connect("127.0.0.1",
		nats.UserInfo("nats", "123qwe"),    // 会一直重连，密码不能错
		nats.RetryOnFailedConnect(true),    // 连接失败重试
		nats.ReconnectWait(30*time.Second), // 重试间隔30秒
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Println("reconnect success")
		}))
	if err != nil {
		log.Fatal(err)
	}
	log.Println("connect nats-server success")
	defer nc.Close()
	// Subscribe
	if _, err := nc.Subscribe("agent1", func(m *nats.Msg) {
		// 监听到master的消息，回复
		log.Println(string(m.Data))
		m.Respond([]byte(time.Now().String()))
	}); err != nil {
		log.Println(err)
	}
	// shutdown signal
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)
	time.Sleep(5 * time.Second)
	s := <-sig
	log.Println("Got signal:", s)
}
