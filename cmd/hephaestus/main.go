package main

import (
	"github.com/desmos-labs/discord-bot/cmd"
	"github.com/rs/zerolog/log"
	"os"
)

func main() {
	executor := cmd.RootCmd()
	if err := executor.Execute(); err != nil {
		log.Error().Err(err).Send()
		os.Exit(1)
	}
}
