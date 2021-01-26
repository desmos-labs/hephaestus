package cmd

import (
	"github.com/desmos-labs/discord-bot/consts"
	"github.com/spf13/cobra"
)

// RootCmd returns a Cobra command allowing to perform various operations
func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   consts.AppName,
		Short: "Official Desmos Discord bot",
	}

	rootCmd.AddCommand(
		StartCmd(),
	)

	return rootCmd
}
