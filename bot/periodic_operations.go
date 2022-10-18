package bot

import "github.com/rs/zerolog/log"

// CleanRoles iterates over all the members and cleans their roles
func (bot *Bot) CleanRoles() {
	log.Debug().Msg("cleaning roles")

	if bot.testnet != nil {
		err := bot.testnet.CleanRoles(bot.discord)
		if err != nil {
			log.Error().Err(err).Msg("error while cleaning testnet roles")
		}
	}

	if bot.mainnet != nil {
		err := bot.mainnet.CleanRoles(bot.discord)
		if err != nil {
			log.Error().Err(err).Msg("error while cleaning mainnet roles")
		}
	}
}
