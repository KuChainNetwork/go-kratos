package sdk

import (
	"context"
	"sync"

	"github.com/KuChainNetwork/go-kratos/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/log"
)

type watcher struct {
	wg             sync.WaitGroup
	blockDataChann chan blockData
	stat           int
	mutex          sync.RWMutex

	scanner  *scanner
	cli      *Client
	wsClient *WSClient
	logger   log.Logger
}

func NewWatcher(ctx Context, fromHeight int64) *watcher {
	return &watcher{
		scanner:        NewScanner(ctx, fromHeight),
		cli:            NewClient(ctx),
		logger:         ctx.Logger(),
		blockDataChann: make(chan blockData, 4096),
	}
}

func (w *watcher) SetLogger(l log.Logger) {
	w.logger = l
	w.scanner.SetLogger(l)
}

func (w *watcher) nextStatStep() {
	w.mutex.Lock()
	defer w.mutex.Unlock()

	w.stat++
}

func (w *watcher) Status() int {
	w.mutex.RLock()
	defer w.mutex.RUnlock()

	return w.stat
}

func (w *watcher) Watch(ctx Context, fromHeight int64, h BlockHandler) error {
	w.logger.Debug("start scanner first", "from", fromHeight)

	// init the client
	wsCli, err := NewWSClient(log.NewNopLogger(), ctx.RpcURL())
	if err != nil {
		return errors.Wrapf(err, "start ws client to chain node error")
	}

	if err := wsCli.Start(); err != nil {
		return errors.Wrapf(err, "start ws client error")
	}
	w.wsClient = wsCli

	scannerCtx, cancelScanner := context.WithCancel(context.Background())
	go func() {
		<-ctx.Done()

		// first stop scanner
		cancelScanner()
		if w.scanner != nil {
			w.scanner.Wait()
		}

		// second, stop watcher
		if w.wsClient != nil {
			w.wsClient.Stop()
			w.wsClient.Wait()
		}

		w.stop()
	}()

	// start handler
	w.logger.Debug("start handler")

	// call all handler h in one gorountinue
	w.wg.Add(1)
	go func() {
		defer w.wg.Done()
		handlerLoop(ctx, "watcher", w.blockDataChann, h)
	}()

	w.blockDataChann <- newBlockDataStart(fromHeight)

	// from init to sync
	w.nextStatStep()

	// first scanner all old blocks
	w.scanner.ScanBlocks(ctx.Clone(scannerCtx), fromHeight,
		func(logger log.Logger, height int64, block *types.FullBlock) error {
			if height%100 == 0 {
				logger.Debug("handler blocks", "height", height)
			}

			w.blockDataChann <- newBlockDataBlk(block)
			return nil
		})

	// start watch blocks from no scaned
	currentHasHandled := w.scanner.LastestBlockHeight()
	w.logger.Debug("current block has get", "height", currentHasHandled)

	// from sync to watch
	w.nextStatStep()

	// on watch
	const queryParamString = "tm.event='NewBlock'"

	if err := w.wsClient.SubscribeBlocks(ctx,
		func(block *types.EventNewBlock) error {
			full, err := w.cli.QueryFullBlock(block.Block.Block.Height)
			if err != nil {
				return errors.Wrapf(err, "query full block error")
			}

			w.blockDataChann <- newBlockDataBlk(&full)
			return nil
		}); err != nil {
		return errors.Wrapf(err, "subscribe error")
	}

	return nil
}

func (w *watcher) Wait() {
	w.logger.Debug("watcher start wait stopped")

	w.wsClient.Wait()
	w.wg.Wait()

	w.logger.Debug("watcher stopped")
}

func (w *watcher) stop() {
	w.logger.Debug("watcher start stop")
	w.blockDataChann <- newBlockDataStop()
	close(w.blockDataChann)
}
