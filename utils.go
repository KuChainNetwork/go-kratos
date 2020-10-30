package sdk

import (
	"os"
	"os/signal"
	"syscall"
)

func SetupCloseHandler(waitFunc func()) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		waitFunc()
		os.Exit(0)
	}()
}

func HoldToClose(waitFunc func()) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	waitFunc()
	os.Exit(0)
}
