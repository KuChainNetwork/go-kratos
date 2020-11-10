package sdk

import (
	"context"
	"sync"

	"github.com/KuChainNetwork/go-kratos/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/log"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

type watcher struct {
	wg             sync.WaitGroup
	blockDataChann chan *types.FullBlock

	scanner *scanner
	logger  tmlog.Logger
}

func NewWatcher(fromHeight int64) *watcher {
	return &watcher{
		scanner:        NewScanner(fromHeight),
		logger:         tmlog.NewNopLogger(),
		blockDataChann: make(chan *types.FullBlock, 128),
	}
}

func (w *watcher) SetLogger(l tmlog.Logger) {
	w.logger = l
	w.scanner.SetLogger(l)
}

func (w *watcher) Watch(lcdURL, rpcURL string, fromHeight int64, h BlockHandler) error {
	w.logger.Debug("start scanner first", "from", fromHeight)

	cli := NewClient(lcdURL)

	// init the ws client
	wsClient, err := NewWSClient(log.NewNopLogger(), rpcURL)
	if err != nil {
		return errors.Wrapf(err, "start ws client to chain node error")
	}

	if err := wsClient.Start(); err != nil {
		return errors.Wrapf(err, "start ws client error")
	}

	w.logger.Debug("start handler")

	// call all handler h in one gorountinue
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		w.logger.Debug("start handler gorountinue")
		for {
			data, ok := <-w.blockDataChann
			if !ok {
				// closed scanner
				w.logger.Info("stop handler gorountinue")
				return
			}

			if err := h(w.logger, data.Height, data); err != nil {
				w.logger.Error("handler err", "err", err.Error())
			}
		}
	}()

	// first scanner all old blocks
	if err := w.scanner.ScanBlocks(lcdURL, fromHeight,
		func(logger tmlog.Logger, height int64, block *types.FullBlock) error {
			if height%100 == 0 {
				logger.Debug("handler blocks", "height", height)
			}

			w.blockDataChann <- block
			return nil
		}); err != nil {
		return errors.Wrapf(err, "scanner last blocks error")
	}

	// start watch blocks from no scaned
	currentHasHandled := w.scanner.CurrentBlockHeight()
	w.logger.Debug("current block to get", "height", currentHasHandled)

	// on watch
	const queryParamString = "tm.event='NewBlock'"

	if err := wsClient.SubscribeBlocks(context.Background(),
		func(block *types.EventNewBlock) error {
			full, err := cli.QueryFullBlock(block.Block.Block.Height)
			if err != nil {
				return errors.Wrapf(err, "query full block error")
			}

			w.blockDataChann <- &full
			return nil
		}); err != nil {
		return errors.Wrapf(err, "subscribe error")
	}

	wsClient.Wait()

	return nil
}

func (w *watcher) Wait() {
	w.wg.Wait()
}

func (w *watcher) Stop() {
	// TODO: stop scan and watch at first then wait all processed
	close(w.blockDataChann)
}
