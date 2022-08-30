package clientApp

import (
	rpcCli "dataProcess/jinn/src/cgo_rpc_client"
	dbDo "dataProcess/tcp_socket/db_do"
)

var client rpcCli.Client

func initClientConn(user, token string) {
	client = rpcCli.Client{User: user, Token: token}
	client.Connect(dbDo.JinnConfig.Host, dbDo.JinnConfig.Port)
}
