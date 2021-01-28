package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/andersfylling/disgord"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/discord-bot/keys"
)

const (
	HelpCmd = "help"
)

// HandleHelp handles the the request for help by the user
func (bot *Bot) HandleHelp(s disgord.Session, data *disgord.MessageCreate) {
	// Consider only those messages starting with "send"
	msg := data.Message
	if !strings.HasPrefix(msg.Content, HelpCmd) {
		return
	}

	log.Debug().Str(keys.LogCommand, HelpCmd).Msgf("received %s command", HelpCmd)

	// Check the command limitation
	if expirationDate := bot.CheckCommandLimit(msg.Author.ID, HelpCmd); expirationDate != nil {
		bot.Reply(msg, s, fmt.Sprintf("Cannot do this now. You will be able to ask for help again on %s",
			expirationDate.Format(time.RFC822)))
		bot.React(msg, s, keys.ReactionTime)
		return
	}

	// Answer to the command
	bot.React(msg, s, keys.ReactionDone)
	bot.Reply(msg, s, fmt.Sprintf(
		"Here are the available commands:\n"+
			"- `!%s` - to get help\n"+
			"- `!%s` - to read the documentation\n"+
			"- `!%s <address>` - to ask for testnet tokens\n",
		HelpCmd,
		DocsCmd,
		SendCmd,
	))
}
