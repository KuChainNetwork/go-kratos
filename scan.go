package sdk

import (
	"sync"

	"github.com/KuChainNetwork/go-kratos/types"
	"github.com/pkg/errors"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

type BlockHandler func(logger tmlog.Logger, height int64, block *types.FullBlock) error

type scanner struct {
	wg             sync.WaitGroup
	logger         tmlog.Logger
	cli            *Client
	blockDataChann chan blockData

	mutex             sync.RWMutex
	latestBlockHeight int64
}

func NewScanner(fromHeight int64) *scanner {
	return &scanner{
		logger:         tmlog.NewNopLogger(),
		blockDataChann: make(chan blockData, 4096),
	}
}

func (s *scanner) SetLogger(l tmlog.Logger) {
	s.logger = l
}

func (s *scanner) LastestBlockHeight() int64 {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	return s.latestBlockHeight
}

func (s *scanner) setToHeight(height int64) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.logger.Debug("setToHeight", "height", height)
	s.latestBlockHeight = height
}

func (s *scanner) ScanBlocks(ctx Context, fromHeight int64, h BlockHandler) {
	s.cli = NewClient(ctx)

	s.wg.Add(1)
	go func() {
		defer s.wg.Done()
		handlerLoop("scanner", s.logger, s.blockDataChann, h)
	}()

	s.blockDataChann <- newBlockDataStart(fromHeight)

	s.wg.Add(1)
	go func() {
		defer func() {
			s.wg.Done()

			// stop the scanner handler
			close(s.blockDataChann)
		}()

		if err := s.scanBlocksImp(ctx.LcdURL(), fromHeight, h); err != nil {
			s.logger.Error("scan block error", "err", err)
		}
	}()

}

func (s *scanner) scanBlocksImp(url string, fromHeight int64, h BlockHandler) error {
	currentBlockHeight := fromHeight
	if currentBlockHeight < 1 {
		currentBlockHeight = 1
	}

	last, err := s.cli.QueryLatestBlock()
	if err != nil {
		return errors.Wrapf(err, "get latest block err")
	}

	if last.DecodeBlock.Height <= currentBlockHeight {
		// has scan all
		return nil
	}

	currToBlockHeight := last.DecodeBlock.Height
	s.setToHeight(currToBlockHeight)

	for {
		curr := currentBlockHeight
		if curr > currToBlockHeight {
			// has to the last
			return s.scanBlocksImp(url, curr, h)
		}

		block, err := s.cli.QueryFullBlock(curr)
		if err != nil {
			return errors.Wrapf(err, "query block %d", curr)
		}

		// to handler loop
		s.blockDataChann <- newBlockDataBlk(&block)

		currentBlockHeight = block.Height + 1
	}
}

func (s *scanner) Wait() {
	s.wg.Wait()
}

func (s *scanner) Stop() {

}
