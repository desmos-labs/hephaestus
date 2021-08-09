package bot

import (
	"github.com/andersfylling/disgord"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/hephaestus/gql"
	"github.com/desmos-labs/hephaestus/utils"
)

// HandleVerify handles a verification request. This request is done by the user who already has a Desmos profile
// that is connected to their Discord profile. With this request they can verify everything has been completed
// successfully and get the role they deserve inside the Discord channel.
// This command has no arguments.
func (bot *Bot) HandleVerify(s disgord.Session, data *disgord.MessageCreate) error {
	username := utils.GetUsername(data.Message)

	// Check to see if the user is a validator
	isValidator, err := gql.CheckIsValidator(bot.verificationCfg.GraphQLEndpoint, username)
	if err != nil {
		return err
	}

	server := s.Guild(data.Message.GuildID)

	var role = disgord.Snowflake(bot.verificationCfg.VerifiedUserRole)
	if isValidator {
		role = disgord.Snowflake(bot.verificationCfg.VerifiedValidatorRole)
	}

	if role.IsZero() {
		log.Debug().Msg("no role found to give to verified users")
		return nil
	}

	err = server.Member(data.Message.Author.ID).AddRole(role)
	if err != nil {
		return err
	}

	bot.Reply(data.Message, s, "Congratulations on verifying your Desmos account with Discord. "+
		"Your role has been updated successfully 🎉")

	return nil
}
