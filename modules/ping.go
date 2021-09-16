package modules

import (
	"io/ioutil"
	"sync"
	"tomato/hosts"
	"tomato/utils"

	"golang.org/x/crypto/ssh"
)

func ExecPing(mode string, groups map[string]interface{}) {
	rp := utils.ResultPrinter{}
	if mode == "ssh" {
		var wg sync.WaitGroup
		defaultuser := hosts.Hosts.GetString("default.user")
		defaultkeyFile := hosts.Hosts.GetString("default.key")
		defaultport := hosts.Hosts.GetInt("default.port")
		for h, v := range groups {
			host := h
			vals := v
			wg.Add(1)
			go func() {
				var user string
				var keyFile string
				var port int
				var password string
				if vals == nil {
					user = defaultuser
					keyFile = defaultkeyFile
					port = defaultport
					password = ""
				} else {
					values := vals.(map[string]interface{})
					if u, ok := values["user"]; ok {
						user = u.(string)
					} else {
						user = defaultuser
					}
					if k, ok := values["key"]; ok {
						keyFile = k.(string)
					} else {
						keyFile = defaultkeyFile
					}
					if p, ok := values["port"]; ok {
						port = p.(int)
					} else {
						port = defaultport
					}
					if pass, ok := values["password"]; ok {
						password = pass.(string)
					} else {
						password = ""
					}
				}
				// 获取ssh client
				var client *ssh.Client
				var errs error
				if len(password) > 0 {
					client, errs = utils.GetSSHClient(host, port, user, 1, password, "")
				} else {
					key, err := ioutil.ReadFile(keyFile)
					if err != nil {
						rp.PrintFail(host, "私钥文件读取失败")
						wg.Done()
						return
					}
					client, errs = utils.GetSSHClient(host, port, user, 2, "", string(key))
				}
				if errs != nil {
					rp.PrintFail(host, "获取 ssh client 失败: %v", errs)
					wg.Done()
					return
				}
				defer client.Close()
				//创建ssh session
				session, err := client.NewSession()
				if err != nil {
					rp.PrintFail(host, "获取 ssh session 失败: %v", errs)
					wg.Done()
					return
				}
				defer session.Close()
				//执行远程命令
				_, err = session.CombinedOutput("echo 1")
				if err != nil {
					rp.PrintFail(host, "远程执行cmd 失败: %v", err)
					wg.Done()
					return
				}
				rp.PrintSucc(host, "PONG")
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
