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
		mergedHandlers := types.MergeHandlers(
			bot.checkCmdLimit(cmdName),
			handler,
		)

		err := mergedHandlers(s, data)
		if err != nil {
			log.Warn().Err(err).Str(types.LogCommand, cmdName).Msg("error while handling command")

			customErr, ok := err.(*types.Error)
			if ok {
				bot.Reply(msg, s, customErr.Description)
				bot.React(msg, s, customErr.Reaction)
			} else {
				bot.Reply(msg, s, err.Error())
				bot.React(msg, s, "🚨")
			}

			return
		}

		bot.React(msg, s, "✅")
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
