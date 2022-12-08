package bot

import (
	"fmt"

	"github.com/desmos-labs/hephaestus/types"
	telebot "gopkg.in/telebot.v3"
)

// HandleHelp handles the the request for help by the user
func (bot *Bot) HandleHelp(ctx telebot.Context) error {
	// Answer to the command
	//TODO: implement verify command if it is needed
	ctx.Reply(fmt.Sprintf(
		"Here are the available commands:\n"+
			"- `/%s`, to get help\n"+
			"- `/%s`, to read the documentation\n"+
			"- `/%s <address>`, to ask for testnet tokens\n"+
			"- `/%s`, to connect your Desmos profile to Telegram\n",
		types.CmdHelp,
		types.CmdDocs,
		types.CmdSend,
		types.CmdConnect,
	))

	return nil
}
