package main

import (
	"encoding/base64"
	"fmt"
	"strconv"

	"github.com/KuChainNetwork/kuchain/chain/types"
	"github.com/spf13/cobra"
	"github.com/tendermint/go-amino"
)

func parseCmd(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "parse",
		Short: "Parse datas marshal by cdc in kuchain",
	}

	cmd.AddCommand(
		ParseCoins(cdc),
		ParseKeyName(cdc),
	)

	return cmd
}

func ParseCoins(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coins [base64Data]",
		Short: "parse coins from base64 data",
		RunE: func(cmd *cobra.Command, args []string) error {
			var v types.Coins

			if err := UnmarshalBinaryBareStr(cdc, args[0], &v); err != nil {
				return err
			}

			fmt.Println(v.String())
			return nil
		},
		Args: cobra.ExactArgs(1),
	}

	return cmd
}

func ParseKeyName(cdc *amino.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "keyname [prefixlen] [str]",
		Short: "parse keyname from base64 data in kratos",
		RunE: func(cmd *cobra.Command, args []string) error {
			bz, err := base64.StdEncoding.DecodeString(args[1])
			if err != nil {
				return fmt.Errorf("hex decode error %s", err.Error())
			}

			prefixLen, err := strconv.Atoi(args[0])
			if err != nil {
				return fmt.Errorf("prefixLen error %s", err.Error())
			}

			v := types.Name{
				Value: bz[prefixLen : prefixLen+types.NameBytesLen],
			}

			if prefixLen != 0 {
				fmt.Printf("prefix : %x, name: %s, key: %x\naccountIDKey: %s\n",
					bz[:prefixLen],
					v.String(),
					bz[prefixLen+types.NameBytesLen+1:],
					types.NewAccountIDFromByte(bz[prefixLen+types.NameBytesLen+1:]).String())
			} else {
				fmt.Printf("name %s, key %x\naccountIDKey %s\n",
					v.String(),
					bz[prefixLen+types.NameBytesLen+1:],
					types.NewAccountIDFromByte(bz[prefixLen+types.NameBytesLen+1:]).String())
			}
			return nil
		},
		Args: cobra.ExactArgs(2),
	}

	return cmd
}
