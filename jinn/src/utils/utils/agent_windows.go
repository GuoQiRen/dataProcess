package utils

import (
	"jinnDataProcessing/jinn/src/utils/log"
	"os"
	"os/signal"
	"syscall"
)

func ServerAgent(_ string, start func(), stop func(), _ func()) {
	start()
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM)
	<-exit
	stop()
}

func ServerExit() {
	log.Fatal("ServerExit")
}
