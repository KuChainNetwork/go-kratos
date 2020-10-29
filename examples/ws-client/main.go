package main

import (
	"context"

	"github.com/KuChainNetwork/kuchain/app"
	"github.com/KuChainNetwork/kuchain/utils/log"
	"github.com/tendermint/tendermint/rpc/jsonrpc/client"
	"go.uber.org/zap"
)

func main() {
	logger := log.NewLogger(zap.NewExample())

	wsClient, err := client.NewWS("tcp://127.0.0.1:26657", "/websocket")
	if err != nil {
		logger.Error("err client", "error", err)
		return
	}

	wsClient.SetCodec(app.MakeCodec())
	wsClient.SetLogger(logger)

	if err := wsClient.Start(); err != nil {
		logger.Error("err client", "error", err)
		return
	}

	defer func() {
		if err := wsClient.Stop(); err != nil {
			logger.Error("err stop", "error", err)
			return
		}
	}()

	if err := wsClient.CallWithArrayParams(context.TODO(), "subscribe", []interface{}{"tm.event='NewBlock'"}); err != nil {
		logger.Error("err Subscribe", "error", err)
		return
	}

	for {
		resp, ok := <-wsClient.ResponsesCh
		if !ok {
			return
		}

		logger.Info("resp", "id", resp.ID, "data", resp.Result)
	}
}
