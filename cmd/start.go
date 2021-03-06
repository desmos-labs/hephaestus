package cmd

import (
	"github.com/spf13/cobra"

	"github.com/desmos-labs/hephaestus/types"

	"github.com/desmos-labs/hephaestus/bot"
	"github.com/desmos-labs/hephaestus/cosmos"
)

// StartCmd returns a Cobra command allowing to start the bot
func StartCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "start [config-file]",
		Short: "Starts the bot using the provided configuration file",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// Read the configuration
			cfg, err := types.Parse(args[0])
			if err != nil {
				return err
			}

			// Crete cosmos client
			cosmosClient, err := cosmos.NewClient(cfg.ChainConfig)
			if err != nil {
				return err
			}

			// Create the bot
			hephaestus, err := bot.Create(cfg.BotConfig, cfg.ThemisConfig, cfg.VerificationConfig, cosmosClient)
			if err != nil {
				return err
			}

			// Start the bot
			hephaestus.Start()

			return nil
		},
	}
}
