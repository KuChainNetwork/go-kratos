package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "scan",
		Short: "scan datas from chain",
	}

	rootCmd.AddCommand(
		ScanAllBlocks(),
	)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("Failed executing CLI command: %s, exiting...\n", err)
		os.Exit(1)
	}
}
