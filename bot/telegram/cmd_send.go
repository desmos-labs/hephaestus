package bot

import (
	"fmt"

	"github.com/desmos-labs/hephaestus/types"
	"github.com/rs/zerolog/log"
	telebot "gopkg.in/telebot.v3"
)

const (
	LogTxHash = "tx_hash"
)

// HandleSendTokens handles the sending of tokens to a user that asks them
func (bot *Bot) HandleSendTokens(ctx telebot.Context) error {
	// Get the network client to be used
	var networkClient = bot.testnet

	// Parse the address to make sure it's valid
	recipient := ctx.Message().Payload
	_, err := networkClient.ParseAddress(recipient)
	if err != nil {
		return types.NewWarnErr(fmt.Sprintf("Invalid recipient address: %s", recipient))
	}

	// Create the message
	res, err := networkClient.SendTokens(recipient, 2000000)
	if err != nil {
		fmt.Println(err)

		return fmt.Errorf("error while sending transaction: %s", err)
	}

	log.Debug().Str(types.LogRecipient, recipient).Str(LogTxHash, res.TxHash).Msg("tokens sent successfully")
	//TODO: limitation implementation
	ctx.Reply(fmt.Sprintf(
		"Your tokens have been sent successfully. You can see it by running `desmos q tx %s`."+
			"If your balance does not update in the next seconds, make sure your node is synced.", res.TxHash,
	))

	return nil
}
