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
	DocsCmd = "docs"
)

// HandleDocs handles the the request for docs by the user
func (bot *Bot) HandleDocs(s disgord.Session, data *disgord.MessageCreate) {
	// Consider only those messages starting with "send"
	msg := data.Message
	if !strings.HasPrefix(msg.Content, DocsCmd) {
		return
	}

	log.Debug().Str(keys.LogCommand, DocsCmd).Msgf("received %s command", DocsCmd)

	// Check the command limitation
	if expirationDate := bot.CheckCommandLimit(msg.Author.ID, DocsCmd); expirationDate != nil {
		bot.Reply(msg, s, fmt.Sprintf("Cannot do this now. You will be able to ask for docs again on %s",
			expirationDate.Format(time.RFC822)))
		bot.React(msg, s, keys.ReactionTime)
		return
	}

	// Answer to the command
	bot.React(msg, s, keys.ReactionDone)
	bot.Reply(msg, s, fmt.Sprintf(
		"Here are a series of useful links:\n"+
			"- General documentation: %s\n"+
			"- Become a validator: %s",
		"https://docs.desmos.network/",
		"https://docs.desmos.network/validators/setup.html",
	))
}
