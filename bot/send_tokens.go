package bot

import (
	"fmt"
	"strings"
	"time"

	"github.com/rs/zerolog/log"

	"github.com/andersfylling/disgord"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/desmos-labs/discord-bot/keys"
)

const (
	SendCmd = "send"

	LogRecipient = "recipient"
	LogTxHash    = "tx_hash"
)

// HandleSendTokens handles the sending of tokens to a user that asks them
func (bot *Bot) HandleSendTokens(s disgord.Session, data *disgord.MessageCreate) {
	// Consider only those messages starting with "send"
	msg := data.Message
	if !strings.HasPrefix(msg.Content, SendCmd) {
		return
	}

	log.Debug().Str(keys.LogCommand, SendCmd).Msgf("received %s command", SendCmd)

	// Check the command limitation
	if expirationDate := bot.CheckCommandLimit(msg.Author.ID, SendCmd); expirationDate != nil {
		bot.Reply(msg, s, fmt.Sprintf("Cannot do this now. You will be able to ask tokens once again on %s",
			expirationDate.Format(time.RFC822)))
		bot.React(msg, s, keys.ReactionTime)
		return
	}

	// Get the recipient
	parts := strings.Split(msg.Content, " ")
	if len(parts) < 2 {
		bot.Reply(msg, s, "Missing recipient")
		bot.React(msg, s, keys.ReactionWarning)
		return
	}

	// Parse the address to make sure it's valid
	addr, err := sdk.AccAddressFromBech32(parts[1])
	if err != nil {
		log.Error().Err(err).Str(LogRecipient, parts[1]).Msg("invalid address")
		bot.React(msg, s, keys.ReactionWarning)
		bot.Reply(msg, s, "invalid address provided")
	}

	// Create the message
	txMsg := &banktypes.MsgSend{
		FromAddress: bot.cosmosClient.AccAddress(),
		ToAddress:   addr.String(),
		Amount:      sdk.NewCoins(sdk.NewCoin("udaric", sdk.NewInt(100000))),
	}

	// Send the transaction
	log.Debug().Str(keys.LogCommand, SendCmd).Str(LogRecipient, addr.String()).Msg("sending tokens")
	res, err := bot.cosmosClient.BroadcastTx(txMsg)

	if err != nil {
		log.Error().Err(err).Str(LogRecipient, addr.String()).Msg("error while sending tokens")
		bot.React(msg, s, keys.ReactionWarning)
		bot.Reply(msg, s, err.Error())
	} else {
		log.Debug().Str(LogRecipient, addr.String()).Str(LogTxHash, res.TxHash).Msg("tokens sent successfully")
		bot.SetCommandLimitation(msg.Author.ID, SendCmd)
		bot.React(msg, s, keys.ReactionDone)
		bot.Reply(msg, s, fmt.Sprintf(
			"Your tokens have been sent successfully. You can see it by running `desmos q tx %s`."+
				"If you balance is not updated in the next seconds, make sure your node is synced.", res.TxHash,
		))
	}
}
