package bot

import (
	"fmt"

	"github.com/andersfylling/disgord"
)

const (
	HelpCmd = "help"
)

// HandleHelp handles the the request for help by the user
func (bot *Bot) HandleHelp(s disgord.Session, data *disgord.MessageCreate) error {
	// Answer to the command
	msg := data.Message
	bot.Reply(msg, s, fmt.Sprintf(
		"Here are the available commands:\n"+
			"- `!%s` - to get help\n"+
			"- `!%s` - to read the documentation\n"+
			"- `!%s <address>` - to ask for testnet tokens\n",
		HelpCmd,
		DocsCmd,
		SendCmd,
	))

	return nil
}
