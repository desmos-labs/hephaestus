package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/hephaestus/types"
)

// NewCmdHandler returns a new command handler for the command that has the given name
func (bot *Bot) NewCmdHandler(cmdName string, handler types.CmdHandler) disgord.HandlerMessageCreate {
	return func(s disgord.Session, data *disgord.MessageCreate) {
		// Consider only those messages starting with "connect"
		msg := data.Message
		if !strings.HasPrefix(msg.Content, cmdName) {
			return
		}

		log.Debug().Str(types.LogCommand, cmdName).Msg("received command")

		// Merge the handler with the limit check
		mergedHandlers := types.MergeHandlers(bot.checkCmdLimit(cmdName), handler)

		// Handle the message
		err := mergedHandlers(s, data)
		if err != nil {
			bot.handleError(msg, s, err)
		} else {
			bot.React(msg, s, "✅")
		}

		// Get the channel details
		channel, err := s.Channel(msg.ChannelID).Get()
		if err != nil {
			bot.handleError(msg, s, err)
			return
		}

		// Delete the message if it's not a DM
		if channel.Type != disgord.ChannelTypeDM {
			err = s.Channel(msg.ChannelID).Message(msg.ID).Delete()
			if err != nil {
				bot.handleError(msg, s, err)
			}
		}
	}
}

func (bot *Bot) handleError(msg *disgord.Message, s disgord.Session, err error) {
	log.Warn().Err(err).Str(types.LogUser, msg.Author.Username).Msg("error while handling command")

	customErr, ok := err.(*types.Error)
	if ok {
		bot.Reply(msg, s, customErr.Description)
		bot.React(msg, s, customErr.Reaction)
	} else {
		bot.Reply(msg, s, err.Error())
		bot.React(msg, s, "🚨")
	}
}

func (bot *Bot) checkCmdLimit(cmdName string) types.CmdHandler {
	return func(s disgord.Session, data *disgord.MessageCreate) error {
		// Check the command limitation
		msg := data.Message
		if expirationDate := bot.CheckCommandLimit(msg.Author.ID, cmdName); expirationDate != nil {
			return types.NewTimeErr(fmt.Sprintf(
				"Command limit reached. You will be able to use this command again on: %s",
				expirationDate.Format(time.RFC822),
			))
		}
		return nil
	}
}
