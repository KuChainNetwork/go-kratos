package main

import (
	"context"
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

func ScanAllBlocks() *cobra.Command {
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

			url := cmd.Flag(FlagURL).Value.String()
			logger := log.NewLoggerByZap(true, "*:debug")

			ctxCancel, cancel := context.WithCancel(context.Background())
			ctx := sdk.NewCtx(ctxCancel).
				WithLogger(logger).
				WithUrls(url, "")

			scanner := sdk.NewScanner(ctx, int64(fromHeight))
			scanner.ScanBlocks(ctx, int64(fromHeight), func(l tmlog.Logger, height int64, block *types.FullBlock) error {
				l.Info("block", "height", block.Height, "id", block.BlockID, "appHash", block.AppHash.String(), "txs", len(block.TxDatas))
				return nil
			})

			sdk.HoldToClose(func() {
				cancel()

				logger.Info("cancel scanner, waiting for stopped")
				scanner.Wait()

				logger.Info("scanner stopped")
			})

			return nil
		},
		Args: cobra.ExactArgs(1),
	}

	return cmd
}

func ScanAllTxs() *cobra.Command {
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

			logger := log.NewLoggerByZap(true, "*:info")
			url := cmd.Flag(FlagURL).Value.String()

			ctxCancel, cancel := context.WithCancel(context.Background())
			ctx := sdk.NewCtx(ctxCancel).
				WithLogger(logger).
				WithUrls(url, "")

			scanner := sdk.NewScanner(ctx, int64(fromHeight))
			scanner.ScanBlocks(ctx, int64(fromHeight), func(l tmlog.Logger, height int64, block *types.FullBlock) error {
				l.Debug("block", "height", block.Height, "id", block.BlockID, "appHash", block.AppHash.String(), "txs", len(block.TxDatas))
				for _, tx := range block.TxDatas {
					l.Info("txs", "height", block.Height, "tx", tx.TxHash)
				}
				return nil
			})

			sdk.HoldToClose(func() {
				cancel()

				logger.Info("cancel scanner, waiting for stopped")
				scanner.Wait()

				logger.Info("scanner stopped")
			})

			return nil
		},
		Args: cobra.ExactArgs(1),
	}

	return cmd
}
