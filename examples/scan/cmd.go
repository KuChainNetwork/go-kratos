package main

import (
	"strconv"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

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

			scanner := sdk.NewScanner(int64(fromHeight))
			scanner.SetLogger(log.NewLoggerByZap(true, "*:info"))

			return scanner.ScanBlocks(viper.GetString(FlagURL), int64(fromHeight), func(l tmlog.Logger, height int64, block *types.FullBlock) error {
				l.Info("block", "height", block.Height, "id", block.BlockID, "appHash", block.AppHash.String(), "txs", len(block.TxDatas))
				return nil
			})
		},
		Args: cobra.ExactArgs(1),
	}

	cmd.Flags().String(FlagURL, "http://127.0.0.1:1317/", "lcd server http rpc url")
	viper.BindPFlag(FlagURL, cmd.Flags().Lookup(FlagURL))

	return cmd
}
