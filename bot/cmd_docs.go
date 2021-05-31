package bot

import (
	"fmt"

	"github.com/andersfylling/disgord"
)

const (
	DocsCmd = "docs"
)

// HandleDocs handles the the request for docs by the user
func (bot *Bot) HandleDocs(s disgord.Session, data *disgord.MessageCreate) error {
	// Answer to the command
	msg := data.Message
	bot.Reply(msg, s, fmt.Sprintf(
		"Here are a series of useful links:\n"+
			"- General documentation: %s\n"+
			"- Become a validator: %s",
		"https://docs.desmos.network/",
		"https://docs.desmos.network/validators/setup.html",
	))
	return nil
}
