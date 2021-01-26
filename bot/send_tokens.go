package bot

import (
	"strings"

	"github.com/andersfylling/disgord"
	sdk "github.com/cosmos/cosmos-sdk/types"
	banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/desmos-labs/discord-bot/consts"
)

func (bot *Bot) handleSendTokens(s disgord.Session, data *disgord.MessageCreate) {
	// Consider only those messages starting with "send"
	msg := data.Message
	if !strings.HasPrefix(msg.Content, "send") {
		return
	}

	// Get the recipient
	parts := strings.Split(msg.Content, " ")
	if len(parts) < 2 {
		bot.Reply(msg, s, "Test")
		bot.React(msg, s, consts.ReactionWarning)
	}

	// Get the sender
	txMsg := &banktypes.MsgSend{
		FromAddress: bot.cosmosClient.AccAddress(),
		ToAddress:   parts[1],
		Amount:      sdk.NewCoins(sdk.NewCoin("udaric", sdk.NewInt(100000))),
	}
	err := bot.cosmosClient.BroadcastTx(txMsg)
	if err != nil {
		bot.Reply(msg, s, err.Error())
	}
}
