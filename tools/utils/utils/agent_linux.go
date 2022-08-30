package utils

import (
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"jinnDataProcessing/tools/utils/log"
)

func ServerAgent(pid string, start func(), stop func(), logSwitch func()) {
	f, err := os.OpenFile(pid, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal("open pid filej %s failed: %s", pid, err)
	}
	if err := syscall.Flock(int(f.Fd()), syscall.LOCK_EX|syscall.LOCK_NB); err != nil {
		log.Fatal("lock filej %s failed: %s", pid, err)
	}
	_, _ = f.WriteString(strconv.Itoa(os.Getpid()))
	start()
	exit := make(chan os.Signal)
	signal.Notify(exit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1)
	for sig := range exit {
		if sig == syscall.SIGUSR1 {
			logSwitch()
		} else {
			stop()
			break
		}
	}
	_ = syscall.Flock(int(f.Fd()), syscall.LOCK_UN)
	f.Close()
	_ = os.Remove(pid)
}

func ServerExit() {
	_ = syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
}
