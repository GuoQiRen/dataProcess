package main

import (
	clientApp "dataProcess/tcp_socket/client_app/client"
	"dataProcess/tcp_socket/db_do"
	"dataProcess/tools/utils/log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Warn("%s\n", "using default path: config/settings.yaml")
		dbDo.ConfigureSetUp("config/settings.yaml")
	} else if len(os.Args) == 2 {
		log.Fatal("%s\n", "taskId is invalid")
	} else {
		dbDo.ConfigureSetUp(os.Args[1])
	}
	clientApp.InitClientApp(dbDo.LinkConfig.Network, dbDo.LinkConfig.Host, dbDo.LinkConfig.Port, os.Args[2]) // network, ip, port, taskId
}
