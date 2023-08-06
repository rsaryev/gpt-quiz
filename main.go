package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"gpt-quiz/cmd"
	"gpt-quiz/config"
	"os"
)

var rootCmd = &cobra.Command{
	Use: "gpt-quiz",
}

func main() {
	err := config.Load()
	if err != nil {
		fmt.Printf("Failed to load environment variables: %v\n", err)
		os.Exit(1)
	}

	rootCmd.AddCommand(cmd.Start())
	rootCmd.AddCommand(cmd.Reset())
	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
