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

func UnmarshalJSON(bz []byte, ptr interface{}) error {
	return Codec.UnmarshalJSON(bz, ptr)
}

func MustMarshalJSON(o interface{}) []byte {
	return Codec.MustMarshalJSON(o)
}
