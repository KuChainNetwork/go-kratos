package sdk

import "github.com/KuChainNetwork/go-kratos/types"

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

const (
	handlerStatInit = iota + 0
	handlerStatSync
	handlerStatWatch
	handlerStatEnd
)

func handlerLoop(ctx Context, name string, ch <-chan blockData, h BlockHandler) {
	logger := ctx.Logger().With("name", name)
	logger.Debug("start handler gorountinue")

	var currHeightToHandler int64

	for {
		data, ok := <-ch
		if !ok {
			// closed scanner, closed by last
			logger.Info("stop handler gorountinue")
			return
		}

		switch data.typ {
		case blockDataTypStart:
			logger.Info("start handler by gorountinue", "height", data.height)
			currHeightToHandler = data.height
		case blockDataTypStop:
			logger.Info("stop handler by cmd gorountinue")
			return
		case blockDataTypBlock:
			if currHeightToHandler == data.blk.Height {
				if err := h(logger, data.blk.Height, data.blk); err != nil {
					logger.Error("handler err", "err", err.Error())
				}
				currHeightToHandler++
			}
		default:
			logger.Error("error data type", "data", data.typ)
		}
	}
}
