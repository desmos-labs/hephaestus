package cmd

import (
	"github.com/spf13/cobra"

	"github.com/desmos-labs/discord-bot/consts"
)

// RootCmd returns a Cobra command allowing to perform various operations
func RootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   consts.AppName,
		Short: "Official Hephaestus",
	}

	rootCmd.AddCommand(
		StartCmd(),
	)

	return rootCmd
}
