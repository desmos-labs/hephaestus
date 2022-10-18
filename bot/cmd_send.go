package bot

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/andersfylling/disgord"

	"github.com/desmos-labs/hephaestus/types"
)

const (
	LogTxHash = "tx_hash"
)

// HandleSendTokens handles the sending of tokens to a user that asks them
func (bot *Bot) HandleSendTokens(s disgord.Session, data *disgord.MessageCreate) error {
	// Get the recipient
	msg := data.Message
	parts := strings.Split(msg.Content, " ")
	if len(parts) < 2 {
		return types.NewWarnErr("Missing recipient")
	}

	// Get the network client to be used
	var networkClient = bot.testnet

	if networkClient == nil {
		return types.NewWarnErr("Network not enabled for this operation")
	}

	// Parse the address to make sure it's valid
	recipient := parts[1]
	_, err := networkClient.ParseAddress(recipient)
	if err != nil {
		return types.NewWarnErr(fmt.Sprintf("Invalid recipient address: %s", recipient))
	}

	// Send the tokens
	res, err := networkClient.SendTokens(recipient, 2000000)
	if err != nil {
		return fmt.Errorf("error while sending transaction: %s", err)
	}

	log.Debug().Str(types.LogRecipient, recipient).Str(LogTxHash, res.TxHash).Msg("tokens sent successfully")
	bot.SetCommandLimitation(msg.Author.ID, types.CmdSend)
	bot.Reply(msg, s, fmt.Sprintf(
		"Your tokens have been sent successfully. You can see it by running `desmos q tx %s`."+
			"If your balance does not update in the next seconds, make sure your node is synced.", res.TxHash,
	))

	return nil
}
