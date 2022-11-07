package bot

import (
	"fmt"
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/hephaestus/types"

	"github.com/desmos-labs/hephaestus/utils"
)

// HandleVerify handles a verification request. This request is done by the user who already has a Desmos profile
// that is connected to their Discord profile. With this request they can verify everything has been completed
// successfully and get the role they deserve inside the Discord channel.
// This command has no arguments.
func (bot *Bot) HandleVerify(s disgord.Session, data *disgord.MessageCreate) error {
	// Get the arguments
	msg := data.Message
	parts := strings.Split(msg.Content, " ")[1:]
	if len(parts) == 0 {
		bot.Reply(msg, s, fmt.Sprintf(`**Verify**
This command allows you to verify your Discord account on this server.
To do this, you have to: 

1. Connect your Desmos Profile with your Discord account using the !%[4]s command.
2. Use the %[1]s command to get your Discord role.

The !%[1]s command should be used as follow:
`+"`!%[1]s <%[2]s/%[3]s>`"+`

Eg. `+"`!%[1]s %[2]s`"+`
`, types.CmdVerify, types.NetworkTestnet, types.NetworkMainnet, types.CmdConnect))
		return nil
	}

	// Get the network client to be used
	var networkClient = bot.testnet
	if parts[0] == types.NetworkMainnet {
		networkClient = bot.mainnet
	}

	if networkClient == nil {
		return types.NewWarnErr("Network not enabled for this operation")
	}

	// Get the role to assign
	role, err := networkClient.GetDiscordRole(utils.GetMsgAuthorUsername(data.Message))
	if err != nil {
		return err
	}

	if role.IsZero() {
		log.Debug().Msg("no role found to give to verified users")
		return nil
	}

	err = s.Guild(data.Message.GuildID).Member(data.Message.Author.ID).AddRole(role)
	if err != nil {
		return err
	}

	bot.Reply(data.Message, s, "Congratulations on verifying your Desmos account with Discord. "+
		"Your role has been updated successfully 🎉")

	return nil
}
