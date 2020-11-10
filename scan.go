package sdk

import (
	"sync"

	"github.com/KuChainNetwork/go-kratos/types"
	"github.com/pkg/errors"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

type BlockHandler func(logger tmlog.Logger, height int64, block *types.FullBlock) error

type scanner struct {
	wg     sync.WaitGroup
	logger tmlog.Logger

	currentBlockHeight int64
	latestBlockHeight  int64
}

func NewScanner(fromHeight int64) *scanner {
	return &scanner{
		currentBlockHeight: fromHeight,
		logger:             tmlog.NewNopLogger(),
	}
}

func (s *scanner) SetLogger(l tmlog.Logger) {
	s.logger = l
}

func (s *scanner) CurrentBlockHeight() int64 {
	return s.currentBlockHeight
}

func (s *scanner) setToHeight(height int64) {
	s.logger.Debug("setToHeight", "height", height)
	s.latestBlockHeight = height
}

func (s *scanner) onBlockGot(block *types.FullBlock) error {
	s.logger.Debug("onBlockGot", "height", block.Height)

	s.currentBlockHeight = block.Height + 1
	return nil
}

func (s *scanner) ScanBlocks(url string, fromHeight int64, h BlockHandler) error {
	// no create a go routine

	cli := NewClient(url)
	s.currentBlockHeight = fromHeight
	if s.currentBlockHeight < 1 {
		s.currentBlockHeight = 1
	}

	last, err := cli.QueryLatestBlock()
	if err != nil {
		return errors.Wrapf(err, "get latest block err")
	}

	if last.DecodeBlock.Height <= s.currentBlockHeight {
		// has scan all
		return nil
	}

	s.setToHeight(last.DecodeBlock.Height)

	for {
		curr := s.currentBlockHeight
		if curr > s.latestBlockHeight {
			// has to the last
			return s.ScanBlocks(url, curr, h)
		}

		block, err := cli.QueryFullBlock(curr)
		if err != nil {
			return errors.Wrapf(err, "query block %d", curr)
		}

		if err := h(s.logger, curr, &block); err != nil {
			return errors.Wrapf(err, "handler block err %d", curr)
		}

		s.onBlockGot(&block)
	}
}
