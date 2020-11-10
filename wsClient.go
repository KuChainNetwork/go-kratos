package sdk

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/KuChainNetwork/go-kratos/types"
	"github.com/KuChainNetwork/kuchain/app"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/rpc/jsonrpc/client"
)

type WsHandler func(typ string, val json.RawMessage) error

type WSClient struct {
	wsCli *client.WSClient
	mutex sync.RWMutex
	wg    sync.WaitGroup

	handlers      map[string]WsHandler
	handlerForAll []WsHandler
}

func NewWSClient(logger log.Logger, addr string) (*WSClient, error) {
	ws, err := client.NewWS(addr, "/websocket")
	if err != nil {
		return nil, errors.Wrapf(err, "new ws error")
	}

	ws.SetCodec(app.MakeCodec())
	ws.SetLogger(logger)

	return &WSClient{
		wsCli:         ws,
		handlers:      make(map[string]WsHandler, 32),
		handlerForAll: make([]WsHandler, 0, 32),
	}, nil
}

func (w *WSClient) SetLogger(l log.Logger) {
	w.wsCli.SetLogger(l)
}

func (w *WSClient) Start() error {
	if err := w.wsCli.Start(); err != nil {
		return err
	}

	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		for {
			resp, ok := <-w.wsCli.ResponsesCh
			if !ok {
				w.wsCli.Logger.Info("stop wait resp for ws client")
				return
			}

			var respData struct {
				Query  string              `json:"query"`
				Data   json.RawMessage     `json:"data"`
				Events map[string][]string `json:"events"`
			}

			if err := json.Unmarshal(resp.Result, &respData); err != nil {
				w.wsCli.Logger.Error("unmarshal json err", "err", err)
				continue
			}

			func() {
				w.mutex.RLock()
				defer w.mutex.RUnlock()

				h, ok := w.handlers[respData.Query]
				if ok {
					err := h(respData.Query, respData.Data)
					if err != nil {
						w.wsCli.Logger.Error("handler error", "err", err)
					}
				}

				for _, h := range w.handlerForAll {
					err := h(respData.Query, respData.Data)
					if err != nil {
						w.wsCli.Logger.Error("handler error", "err", err)
					}
				}
			}()
		}
	}()

	return nil
}

func (w *WSClient) Stop() error {
	return w.wsCli.Stop()
}

func (w *WSClient) Wait() {
	w.wsCli.Wait()
	w.wg.Wait()
}

func (w *WSClient) AddHandler(handler WsHandler) {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.handlerForAll = append(w.handlerForAll, handler)
}

func (w *WSClient) Subscribe(ctx context.Context, query string, handler WsHandler) error {
	w.mutex.Lock()

	if _, ok := w.handlers[query]; ok {
		w.mutex.Unlock()
		return errors.New("subscribe query already subscribed")
	}

	w.handlers[query] = handler

	w.mutex.Unlock()

	// start subscribe
	if err := w.wsCli.CallWithArrayParams(ctx, "subscribe", []interface{}{query}); err != nil {
		w.wsCli.Logger.Error("err Subscribe", "error", err)
		return errors.Wrapf(err, "call subscribe error %s", query)
	}

	w.wsCli.Logger.Info("subscribe", "query", query)

	return nil
}

func (w *WSClient) SubscribeBlocks(ctx context.Context, handler func(evtBlock *types.EventNewBlock) error) error {
	const queryParamString = "tm.event='NewBlock'"
	return w.Subscribe(context.Background(), queryParamString,
		func(typ string, val json.RawMessage) error {
			evt := types.EventNewBlock{}

			if err := Codec.UnmarshalJSON(val, &evt); err != nil {
				return errors.Wrapf(err, "unmarlshal block from ws error")
			}

			return handler(&evt)
		})
}
