package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func makeCmd(cmds ...*cobra.Command) *cobra.Command {
	for _, cmd := range cmds {
		cmd.Flags().StringP(FlagURL, "u", "http://127.0.0.1:1317/", "lcd server http rpc url")
	}

	rootCmd := &cobra.Command{
		Use:   "kratos-tail",
		Short: "tail datas from chain",
	}

	rootCmd.AddCommand(cmds...)

	return rootCmd
}

func main() {
	rootCmd := makeCmd(
		TailBlocks(),
		TailTxs(),
	)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}
