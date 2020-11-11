package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

func makeCmd(cmds ...*cobra.Command) *cobra.Command {
	for _, cmd := range cmds {
		cmd.Flags().StringP(FlagURL, "u", "http://127.0.0.1:1317/", "lcd server http rpc url")
		cmd.Flags().StringP(FlagRPCURL, "r", "http://127.0.0.1:26657", "node http rpc url")
	}

	rootCmd := &cobra.Command{
		Use:   "kratos-tail",
		Short: "tail datas from chain",
	}

	rootCmd.AddCommand(cmds...)

	return rootCmd
}

func main() {
	runtime.GOMAXPROCS(4)

	rootCmd := makeCmd(
		TailBlocks(),
		TailEvents(),
	)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}
