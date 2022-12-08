package bot

import (
	"fmt"

	"github.com/desmos-labs/hephaestus/types"

	"github.com/andersfylling/disgord"
)

// HandleHelp handles the the request for help by the user
func (bot *Bot) HandleHelp(s disgord.Session, data *disgord.MessageCreate) error {
	// Answer to the command
	msg := data.Message
	bot.Reply(msg, s, fmt.Sprintf(
		"Here are the available commands:\n"+
			"- `!%s`, to get help\n"+
			"- `!%s`, to read the documentation\n"+
			"- `!%s <address>`, to ask for testnet tokens\n"+
			"- `!%s`, to connect your Desmos profile to Discord\n"+
			"- `!%s`, to verify the connection between Discord and your Desmos profile\n",
		types.CmdHelp,
		types.CmdDocs,
		types.CmdSend,
		types.CmdConnect,
		types.CmdVerify,
	))

	return nil
}
