package main

import (
	"context"
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"

	sdk "github.com/KuChainNetwork/go-kratos"
	"github.com/KuChainNetwork/go-kratos/types"
	"github.com/KuChainNetwork/kuchain/chain/config"
	tmlog "github.com/tendermint/tendermint/libs/log"
)

const (
	FlagURL    = "url"
	FlagRPCURL = "rpc"
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

			lcdURL := cmd.Flag(FlagURL).Value.String()
			rpcURL := cmd.Flag(FlagRPCURL).Value.String()

			ctx := sdk.NewCtx(context.Background()).
				WithUrls(lcdURL, rpcURL).
				WithLogger(sdk.NewLoggerByZap(true, "*:debug"))

			watcher := sdk.NewWatcher(int64(fromHeight))
			watcher.SetLogger(sdk.NewLoggerByZap(true, "*:debug"))

			if err := watcher.Watch(ctx, int64(fromHeight),
				func(logger tmlog.Logger, height int64, block *types.FullBlock) error {
					logger.Info("on block", "height", height, "id", block.BlockID, "appHash", block.AppHash.String())
					return nil
				}); err != nil {
				return errors.Wrapf(err, "watcher error")
			}

			watcher.Wait()

			return nil
		},
		Args: cobra.ExactArgs(1),
	}

	return cmd
}
