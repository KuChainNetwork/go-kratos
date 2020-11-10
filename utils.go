package sdk

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/KuChainNetwork/kuchain/utils/log"
	tmflags "github.com/tendermint/tendermint/libs/cli/flags"
	tmlog "github.com/tendermint/tendermint/libs/log"
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

func NewLoggerByZap(isTrace bool, logLevelStr string) tmlog.Logger {
	zapLogger := log.NewZapLogger(isTrace)

	// warp zap log to logger, it will add caller skip 1
	logger := log.NewLogger(zapLogger)

	// add caller skip by 2, as warp and level log
	logger = logger.WithCallerSkip(1)

	// process log level for cosmos-sdk, , it will add caller skip 1
	loggerByLevel, err := tmflags.ParseLogLevel(logLevelStr, logger, "*:info")
	if err != nil {
		panic(err)
	}

	return loggerByLevel
}
