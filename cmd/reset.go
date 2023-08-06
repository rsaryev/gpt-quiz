package cmd

import (
	"github.com/spf13/cobra"
	"gpt-quiz/config"
)

func Reset() *cobra.Command {
	var cmd = &cobra.Command{
		Use:     "reset",
		Aliases: []string{"r"},
		Short:   "Reset config",
		RunE: func(cmd *cobra.Command, args []string) error {
			err := config.Remove()
			if err != nil {
				return err
			}
			return nil
		},
	}

	return cmd
}
