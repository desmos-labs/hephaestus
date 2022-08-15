package bot

import "github.com/rs/zerolog/log"

// CleanRoles iterates over all the members and cleans their roles
func (bot *Bot) CleanRoles() {
	log.Debug().Msg("cleaning roles")

	err := bot.testnet.CleanRoles(bot.discord)
	if err != nil {
		log.Error().Err(err).Msg("error while cleaning testnet roles")
	}

	err = bot.mainnet.CleanRoles(bot.discord)
	if err != nil {
		log.Error().Err(err).Msg("error while cleaning mainnet roles")
	}
}
