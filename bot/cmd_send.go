package bot

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/andersfylling/disgord"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/desmos-labs/hephaestus/types"
)

const (
	SendCmd = "send"

	LogRecipient = "recipient"
	LogTxHash    = "tx_hash"
)

// HandleSendTokens handles the sending of tokens to a user that asks them
func (bot *Bot) HandleSendTokens(s disgord.Session, data *disgord.MessageCreate) error {
	// Get the recipient
	msg := data.Message
	parts := strings.Split(msg.Content, " ")
	if len(parts) < 2 {
		return types.NewWarnErr("Missing recipient")
	}

	// Parse the address to make sure it's valid
	recipient := parts[1]
	addr, err := sdk.AccAddressFromBech32(recipient)
	if err != nil {
		return types.NewWarnErr(fmt.Sprintf("Invalid recipient address: %s", recipient))
	}

	// Create the message
	txMsg := &banktypes.MsgSend{
		FromAddress: bot.cosmosClient.AccAddress(),
		ToAddress:   addr.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin("udaric", sdk.NewInt(2000000))),
	}

	// Send the transaction
	log.Debug().Str(types.LogCommand, SendCmd).Str(LogRecipient, addr.String()).Msg("sending tokens")
	res, err := bot.cosmosClient.BroadcastTx(txMsg)
	if err != nil {
		return fmt.Errorf("error while sending transaction: %s", err)
	}

	log.Debug().Str(LogRecipient, addr.String()).Str(LogTxHash, res.TxHash).Msg("tokens sent successfully")
	bot.SetCommandLimitation(msg.Author.ID, SendCmd)
	bot.Reply(msg, s, fmt.Sprintf(
		"Your tokens have been sent successfully. You can see it by running `desmos q tx %s`."+
			"If your balance does not update in the next seconds, make sure your node is synced.", res.TxHash,
	))

	return nil
}
