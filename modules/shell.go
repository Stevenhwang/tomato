package modules

import (
	"io/ioutil"
	"strings"
	"sync"
	"tomato/utils"

	"golang.org/x/crypto/ssh"
)

func ExecShell(mode string, groups map[string]interface{}, args []string) {
	rp := utils.ResultPrinter{}
	fillG := utils.FillParams(groups)
	if mode == "ssh" {
		var wg sync.WaitGroup
		for h, v := range fillG {
			host := h
			vals := v
			wg.Add(1)
			go func() {
				// 获取ssh client
				var client *ssh.Client
				var errs error
				if len(vals["password"].(string)) > 0 {
					client, errs = utils.GetSSHClient(host, vals["port"].(int), vals["user"].(string), 1, vals["password"].(string), "")
				} else {
					key, err := ioutil.ReadFile(vals["key"].(string))
					if err != nil {
						rp.PrintFail(host, "私钥文件读取失败")
						wg.Done()
						return
					}
					client, errs = utils.GetSSHClient(host, vals["port"].(int), vals["user"].(string), 2, "", string(key))
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
					rp.PrintFail(host, "获取 ssh session 失败: %v", err)
					wg.Done()
					return
				}
				defer session.Close()
				//执行远程命令
				cmd := strings.Join(args, " ")
				combo, err := session.CombinedOutput(cmd)
				if err != nil {
					rp.PrintFail(host, "远程执行cmd 失败: %s %v", string(combo), err)
					wg.Done()
					return
				}
				rp.PrintSucc(host, string(combo))
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
