package agent

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"time"
	"tomato/utils"

	"github.com/dustin/go-humanize"
	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/shirou/gopsutil/v3/disk"
	"github.com/shirou/gopsutil/v3/host"
	"github.com/shirou/gopsutil/v3/load"
	"github.com/shirou/gopsutil/v3/mem"
)

type response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func shellExec(cmd string) string {
	res := exec.Command("bash", "-c", cmd)
	output, _ := res.CombinedOutput()
	return string(output)
}

func Start() {
	v, _ := mem.VirtualMemory()
	log.Printf("MEM: Total: %s, Free:%s, UsedPercent:%.2f%%\n", humanize.Bytes(v.Total), humanize.Bytes(v.Free), v.UsedPercent)
	h, _ := host.Info()
	log.Println(h)
	d, _ := disk.Usage("/")
	log.Printf("DISK: Total: %s, Free:%s, UsedPercent:%.2f%%\n", humanize.Bytes(d.Total), humanize.Bytes(d.Free), d.UsedPercent)
	l, _ := load.Avg()
	log.Printf("LOAD: load1: %f, load5: %.2f, load15: %.2f", l.Load1, l.Load5, l.Load15)
	// register agent
	var id string
	reglock := "./register.lock"
	if utils.PathExists(reglock) {
		c, _ := os.ReadFile(reglock)
		id = string(c)
	} else {
		id = uuid.New().String()
	}
	info := utils.Agent{Name: "agent1",
		Info: utils.Info{
			ID:   id,
			Mem:  utils.MD{Total: humanize.Bytes(v.Total), Free: humanize.Bytes(v.Free), UsedPercent: fmt.Sprintf("%.2f%%", v.UsedPercent)},
			Disk: utils.MD{Total: humanize.Bytes(d.Total), Free: humanize.Bytes(d.Free), UsedPercent: fmt.Sprintf("%.2f%%", d.UsedPercent)},
			Load: utils.LD{Load1: fmt.Sprintf("%.2f", l.Load1), Load5: fmt.Sprintf("%.2f", l.Load5), Load15: fmt.Sprintf("%.2f", l.Load15)},
		}}
	post, err := json.Marshal(&info)
	if err != nil {
		panic(err)
	}
	pbody := bytes.NewBuffer(post)
	httpClient := http.Client{Timeout: 5 * time.Second}
	resp, err := httpClient.Post("http://192.168.1.106:1323/register", "application/json", pbody)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	r := &response{}
	json.Unmarshal(body, r)
	if r.Code != 0 {
		err := errors.New(r.Message)
		panic(err)
	}
	log.Println(r.Message)
	// write lock file
	os.WriteFile(reglock, []byte(id), 0600)
	// start subscribe
	nc, err := nats.Connect("192.168.1.106",
		nats.UserInfo("nats", "123qwe"),    // ?????????????????????????????????
		nats.RetryOnFailedConnect(true),    // ??????????????????
		nats.ReconnectWait(30*time.Second), // ????????????30???
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
		// ?????????master??????????????????
		log.Println(string(m.Data))
		res := shellExec(string(m.Data))
		m.Respond([]byte(res))
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
