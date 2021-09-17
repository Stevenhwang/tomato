package modules

import (
	"io/ioutil"
	"os"
	"path"
	"strings"
	"sync"
	"tomato/utils"

	"github.com/pkg/sftp"

	"golang.org/x/crypto/ssh"
)

func ExecCopy(mode string, groups map[string]interface{}, args []string) {
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
				if len(args) <= 1 {
					rp.PrintFail(host, "必需参数不足(src和dest)")
					wg.Done()
					return
				}
				var src string
				var dest string
				for _, arg := range args {
					if strings.Contains(arg, "src") {
						src = strings.ReplaceAll(arg, "src=", "")
					}
					if strings.Contains(arg, "dest") {
						dest = strings.ReplaceAll(arg, "dest=", "")
					}
				}
				if len(src) == 0 || len(dest) == 0 {
					rp.PrintFail(host, "必需参数不足(src和dest)")
					wg.Done()
					return
				}
				if !utils.PathExists(src) {
					rp.PrintFail(host, "src路径不存在")
					wg.Done()
					return
				}
				if utils.IsDir(src) {
					rp.PrintFail(host, "src必须是文件而不是文件夹")
					wg.Done()
					return
				}
				if strings.HasSuffix(dest, "/") {
					rp.PrintFail(host, "dest必须是文件而不是文件夹")
					wg.Done()
					return
				}
				sftpClient, err := sftp.NewClient(client)
				if err != nil {
					rp.PrintFail(host, "创建 sftp client 失败: %v", err)
					wg.Done()
					return
				}
				defer sftpClient.Close()
				srcFile, err := os.Open(src)
				if err != nil {
					rp.PrintFail(host, "打开 src 文件失败: %v", err)
					wg.Done()
					return
				}
				defer srcFile.Close()
				// 创建远端文件夹
				destDir := path.Dir(dest)
				err = sftpClient.MkdirAll(destDir)
				if err != nil {
					rp.PrintFail(host, "创建 dest 文件夹失败: %v", err)
					wg.Done()
					return
				}
				destFile, err := sftpClient.Create(dest)
				if err != nil {
					rp.PrintFail(host, "创建 dest 文件失败: %v", err)
					wg.Done()
					return
				}
				defer destFile.Close()
				// 开始复制
				buf := make([]byte, 1024)
				for {
					n, _ := srcFile.Read(buf)
					if n == 0 {
						break
					}
					destFile.Write(buf)
				}
				rp.PrintSucc(host, "文件传输成功，远端文件路径: %s", dest)
				wg.Done()
			}()
		}
		wg.Wait()
	}
}
