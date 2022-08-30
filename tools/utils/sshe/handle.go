package sshe

import (
	"dataProcess/constants"
	"dataProcess/logger"
	"golang.org/x/crypto/ssh"
	"net"
	"os"
	"time"
)

var tryConnNum = 5

func TryGetSshConn(userName, password, containerIp, network, containerPort string) (cli *ssh.Client, err error) {
	cli, err = SshConn(userName, password, containerIp, network, containerPort)
	for {
		if err != nil {
			logger.Error(err.Error())
			if tryConnNum != 0 {
				cli, err = SshConn(userName, password, containerIp, network, containerPort)
				tryConnNum--
			} else {
				return nil, err
			}
		} else {
			break
		}
	}
	return
}

func SshConn(userName, password, containerIp, network, containerPort string) (cli *ssh.Client, err error) {
	conf := &ssh.ClientConfig{
		User: userName,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		Timeout: 30 * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	cli, err = ssh.Dial(network, containerIp+constants.Colon+containerPort, conf)
	if err != nil {
		return
	}
	return
}

func ExecCmd(cli *ssh.Client, cmdInfo string, sign chan bool) (err error) {
	defer cli.Close()

	session, err := cli.NewSession()
	if err != nil {
		return
	}
	session.Stderr = os.Stderr
	defer session.Close()

	sign <- true
	err = session.Run(cmdInfo)
	if err != nil {
		return
	}
	return
}
