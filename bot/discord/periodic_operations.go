package bot

import "github.com/rs/zerolog/log"

// RefreshRoles iterates over all the members and cleans their roles
func (bot *Bot) RefreshRoles() {
	log.Info().Timestamp().Msg("refreshing roles")

	if bot.testnet != nil {
		err := bot.testnet.RefreshRoles(bot.discord)
		if err != nil {
			log.Error().Err(err).Msg("error while cleaning testnet roles")
		}
	}

	if bot.mainnet != nil {
		err := bot.mainnet.RefreshRoles(bot.discord)
		if err != nil {
			log.Error().Err(err).Msg("error while cleaning mainnet roles")
		}
	}

	log.Info().Timestamp().Msg("roles refreshed")
}
