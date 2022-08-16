package bot

import (
	"fmt"

	telebot "gopkg.in/telebot.v3"
)

// HandleDocs handles the the request for docs by the user
func (bot *Bot) HandleDocs(ctx telebot.Context) error {
	// Answer to the command
	ctx.Reply(fmt.Sprintf(
		"Here are a series of useful links:\n"+
			"- General documentation: %s\n"+
			"- Become a validator: %s",
		"https://docs.desmos.network/",
		"https://docs.desmos.network/validators/setup.html"))
	return nil
}
