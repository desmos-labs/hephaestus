package bot

import (
	"fmt"
	"strings"

	"github.com/rs/zerolog/log"

	"github.com/andersfylling/disgord"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/desmos-labs/discord-bot/consts"
)

const (
	SendCmd = "send"

	LogRecipient = "recipient"
	LogTxHash    = "tx_hash"
)

// handleSendTokens handles the sending of tokens to a user that asks them
func (bot *Bot) handleSendTokens(s disgord.Session, data *disgord.MessageCreate) {
	// Consider only those messages starting with "send"
	msg := data.Message
	if !strings.HasPrefix(msg.Content, "send") {
		return
	}

	log.Debug().Str(consts.LogCommand, SendCmd).Msg("received command")

	// Get the recipient
	parts := strings.Split(msg.Content, " ")
	if len(parts) < 2 {
		bot.Reply(msg, s, "Missing recipient")
		bot.React(msg, s, consts.ReactionWarning)
	}

	log.Debug().Str(consts.LogCommand, SendCmd).Str(LogRecipient, parts[1]).Msg("sending tokens")

	// Get the sender
	txMsg := &banktypes.MsgSend{
		FromAddress: bot.cosmosClient.AccAddress(),
		ToAddress:   parts[1],
		Amount:      sdk.NewCoins(sdk.NewCoin("udaric", sdk.NewInt(100000))),
	}

	res, err := bot.cosmosClient.BroadcastTx(txMsg)
	if err != nil {
		log.Error().Err(err).Str(LogRecipient, parts[1]).Msg("error while sending tokens")
		bot.React(msg, s, consts.ReactionWarning)
		bot.Reply(msg, s, err.Error())
	} else {
		log.Debug().Str(LogRecipient, parts[1]).Str(LogTxHash, res.TxHash).Msg("tokens sent successfully")
		bot.React(msg, s, consts.ReactionDone)
		bot.Reply(msg, s, fmt.Sprintf("Tokens sent with transaction `%s`", res.TxHash))
	}
}
