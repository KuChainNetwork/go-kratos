package main

import (
	"context"
	"encoding/json"
	"flag"

	sdk "github.com/KuChainNetwork/go-kratos"
	"github.com/KuChainNetwork/kuchain/utils/log"
	"go.uber.org/zap"
)

var addr = flag.String("addr", "tcp://127.0.0.1:26657", "http service address")
var query = flag.String("query", "tm.event='NewBlock'", "http service address")

func main() {
	logger := log.NewLogger(zap.NewExample())

	wsClient, err := sdk.NewWSClient(logger, *addr)
	if err != nil {
		logger.Error("create ws client", "error", err)
		return
	}

	if err := wsClient.Start(); err != nil {
		logger.Error("start ws client", "err", err)
		return
	}

	defer func() {
		if err := wsClient.Stop(); err != nil {
			logger.Error("err stop", "error", err)
			return
		}
	}()

	if err := wsClient.Subscribe(context.TODO(), *query, func(typ string, data json.RawMessage) error {
		logger.Info("data", "typ", typ)
		return nil
	}); err != nil {
		logger.Error("Subscribe ws client", "err", err)
		return
	}

	wsClient.Wait()
}
