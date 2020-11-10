package sdk

import (
	"context"
	"sync"

	"github.com/KuChainNetwork/go-kratos/types"
	"github.com/pkg/errors"
	"github.com/tendermint/tendermint/libs/log"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

const (
	blockDataTypStart = iota + 1
	blockDataTypBlock
	blockDataTypStop
)

type blockData struct {
	typ    int
	height int64
	blk    *types.FullBlock
}

func newBlockDataStart(fromHeight int64) blockData {
	return blockData{
		typ:    blockDataTypStart,
		height: fromHeight,
	}
}

func newBlockDataBlk(block *types.FullBlock) blockData {
	return blockData{
		typ:    blockDataTypBlock,
		height: block.Height,
		blk:    block,
	}
}

func newBlockDataStop() blockData {
	return blockData{
		typ: blockDataTypStop,
	}
}

type watcher struct {
	wg             sync.WaitGroup
	blockDataChann chan blockData

	scanner *scanner
	logger  tmlog.Logger
}

func NewWatcher(fromHeight int64) *watcher {
	return &watcher{
		scanner:        NewScanner(fromHeight),
		logger:         tmlog.NewNopLogger(),
		blockDataChann: make(chan blockData, 4096),
	}
}

func (w *watcher) SetLogger(l tmlog.Logger) {
	w.logger = l
	w.scanner.SetLogger(l)
}

func (w *watcher) handlerLoop(h BlockHandler) {
	w.logger.Debug("start handler gorountinue")

	var currHeightToHandler int64

	for {
		data, ok := <-w.blockDataChann
		if !ok {
			// closed scanner
			w.logger.Info("stop handler gorountinue")
			return
		}

		switch data.typ {
		case blockDataTypStart:
			w.logger.Info("start handler by gorountinue", "height", data.height)
			currHeightToHandler = data.height
		case blockDataTypStop:
			w.logger.Info("stop handler by cmd gorountinue")
			return
		case blockDataTypBlock:
			if currHeightToHandler == data.blk.Height {
				if err := h(w.logger, data.blk.Height, data.blk); err != nil {
					w.logger.Error("handler err", "err", err.Error())
				}
				currHeightToHandler++
			}
		default:
			w.logger.Error("error data type", "data", data.typ)
		}
	}
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
		w.handlerLoop(h)
	}()

	w.blockDataChann <- newBlockDataStart(fromHeight)

	// first scanner all old blocks
	if err := w.scanner.ScanBlocks(lcdURL, fromHeight,
		func(logger tmlog.Logger, height int64, block *types.FullBlock) error {
			if height%100 == 0 {
				logger.Debug("handler blocks", "height", height)
			}

			w.blockDataChann <- newBlockDataBlk(block)
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

			w.blockDataChann <- newBlockDataBlk(&full)
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
