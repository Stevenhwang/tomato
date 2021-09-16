package utils

import (
	"errors"
	"fmt"
	"time"

	"golang.org/x/crypto/ssh"
)

func GetSSHClient(host string, port int, user string, authType int, password string, key string) (*ssh.Client, error) {
	if authType != 1 && authType != 2 {
		return nil, errors.New("认证方式错误")
	}
	var client *ssh.Client
	var errs error
	addr := fmt.Sprintf("%s:%d", host, port)
	if authType == 1 {
		config := &ssh.ClientConfig{
			Timeout:         2 * time.Second,
			User:            user,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Auth:            []ssh.AuthMethod{ssh.Password(password)},
		}
		client, errs = ssh.Dial("tcp", addr, config)
	} else {
		signer, err := ssh.ParsePrivateKey([]byte(key))
		if err != nil {
			return nil, fmt.Errorf("解析私钥错误: %v", err)
		}
		keyConfig := &ssh.ClientConfig{
			Timeout:         2 * time.Second,
			User:            user,
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		}
		client, errs = ssh.Dial("tcp", addr, keyConfig)
	}
	if errs != nil {
		return nil, errs
	}
	return client, nil
}
