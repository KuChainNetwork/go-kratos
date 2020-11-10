package main

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	sdk "github.com/KuChainNetwork/go-kratos"
	"github.com/KuChainNetwork/go-kratos/types"
	"github.com/KuChainNetwork/kuchain/chain/config"
	"github.com/KuChainNetwork/kuchain/utils/log"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

const (
	FlagURL = "url"
)

func TailBlocks() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blocks [fromHeight]",
		Short: "scan all blocks and log to json",
		RunE: func(cmd *cobra.Command, args []string) error {
			config.SealChainConfig()

			fromHeight, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.Wrapf(err, "fromHeight should be an integer %s", args[0])
			}

			if fromHeight < 1 {
				fromHeight = 1
			}

			scanner := sdk.NewScanner(int64(fromHeight))
			scanner.SetLogger(log.NewLoggerByZap(true, "*:info"))

			return scanner.ScanBlocks(cmd.Flag(FlagURL).Value.String(), int64(fromHeight), func(l tmlog.Logger, height int64, block *types.FullBlock) error {
				l.Info("block", "height", block.Height, "id", block.BlockID, "appHash", block.AppHash.String(), "txs", len(block.TxDatas))
				return nil
			})
		},
		Args: cobra.ExactArgs(1),
	}

	return cmd
}

func TailTxs() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "txs [fromHeight]",
		Short: "scan all txs and log to json",
		RunE: func(cmd *cobra.Command, args []string) error {
			config.SealChainConfig()

			fromHeight, err := strconv.Atoi(args[0])
			if err != nil {
				return errors.Wrapf(err, "fromHeight should be an integer %s", args[0])
			}

			if fromHeight < 1 {
				fromHeight = 1
			}

			scanner := sdk.NewScanner(int64(fromHeight))
			scanner.SetLogger(log.NewLoggerByZap(true, "*:info"))

			return scanner.ScanBlocks(cmd.Flag(FlagURL).Value.String(), int64(fromHeight), func(l tmlog.Logger, height int64, block *types.FullBlock) error {
				l.Debug("block", "height", block.Height, "id", block.BlockID, "appHash", block.AppHash.String(), "txs", len(block.TxDatas))
				for _, tx := range block.TxDatas {
					l.Info("txs", "height", block.Height, "tx", tx.TxHash)
				}
				return nil
			})
		},
		Args: cobra.ExactArgs(1),
	}

	return cmd
}