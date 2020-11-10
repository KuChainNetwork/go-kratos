package main

import (
	"fmt"
	"os"

	"github.com/KuChainNetwork/kuchain/app"
	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "set-parse-tool",
		Short: "Command line interface parser datas in kuchain",
	}

	cdc := app.MakeCodec()

	rootCmd.AddCommand(
		parseCmd(cdc),
	)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}
