package main

import (
	"dataProcess/tcp_socket/db_do"
	"dataProcess/tcp_socket/server_app/service"
	"dataProcess/tools/utils/log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Warn("%s\n", "using default path: config/settings.yaml")
		dbDo.ConfigureSetUp("config/settings.yaml")
	} else {
		dbDo.ConfigureSetUp(os.Args[1])
	}
	service.InitServerApp(dbDo.LinkConfig.Network, dbDo.LinkConfig.Host, dbDo.LinkConfig.Port)
}
