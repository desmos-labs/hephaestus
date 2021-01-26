package main

import (
	"os"

	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/discord-bot/cmd"
)

func main() {
	executor := cmd.RootCmd()
	if err := executor.Execute(); err != nil {
		log.Error().Err(err).Send()
		os.Exit(1)
	}
}
