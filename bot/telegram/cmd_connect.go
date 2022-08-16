package bot

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/desmos-labs/hephaestus/types"
	"github.com/desmos-labs/hephaestus/utils"
	telebot "gopkg.in/telebot.v3"
)

type CallData struct {
	// Username is the plain-text Telegram username of the user that wants to be verified
	Username string `json:"username"`
}

// NewCallData returns a new CallData instance
func NewCallData(username string) *CallData {
	return &CallData{
		Username: username,
	}
}

// HandleConnect handle a connection request. This request is done by the users when they want to connect their Desmos
// profile with their Telegram account.
// The command expects one single argument which must be the JSON object returned from the "desmos sign" command.
//
// The handling of the command will fail in the following occasions:
// 1. The signed value does not correspond to the username of the user sending the message
// 2. Any of the values are badly encoded
func (bot *Bot) HandleConnect(ctx telebot.Context) error {
	parts := strings.SplitN(ctx.Message().Payload, " ", 2)
	if len(parts) != 2 {
		ctx.Reply(fmt.Sprintf(`**Connect**
This command allows you to connect your Telegram account to your Desmos profile.
To do this, you have to: 

1. Sign your Telegram username using the Desmos CLI or any Desmos-compatible application.
2. Use the %[1]s command to send the signature result. 

__Signing your Telegram username__
1. Copy your Telegram username by clicking on it in the bottom part of your Telegram client. 

2. Open your Desmos CLI or application, and sign your username. 
   If you use the Desmos CLI, you can do this by using the following command:
   `+"`desmos sign <Telegram username> --from <your-key>`"+`
   
   Eg. `+"`desmos sign \"foo_123\" --from foo`"+`

__Sending the signed value__
The sign command should return a JSON object. The last thing you have to do is now send it to me using the %[1]s command. To do this, simply send me a message as the following: 
`+"`/%[1]s <%[2]s/%[3]s> <JSON>`"+`

Eg. `+"`/%[1]s %[2]s {...}`"+`
`, types.CmdConnect, types.NetworkTestnet, types.NetworkMainnet))
		return nil
	}

	// Get the network client to be used
	var networkClient = bot.testnet
	if parts[0] == types.NetworkMainnet {
		networkClient = bot.mainnet
	}

	// Get the signature data
	username := ctx.Sender().Username
	signatureData, err := utils.GetSignatureData(parts[1])
	if err != nil {
		return err
	}

	// Upload the data to Themis
	err = networkClient.UploadDataToThemis(username, bot.cfg.Name, signatureData)
	// if err != nil {
	// 	return err
	// }

	// Return to the user the call data for the Desmos command
	callDataBz, err := json.Marshal(NewCallData(username))
	if err != nil {
		return types.NewWarnErr("Error while serializing call data: %s", err)
	}
	ctx.Reply(fmt.Sprintf("Your verification data has been stored successfully. "+
		"All you have to do now is execute the following command:\n"+
		"```"+
		"desmos tx profiles link-app ibc-profiles [channel] Telegram \"%[1]s\" %[2]s --packet-timeout-height 0-0 --packet-timeout-timestamp %[3]d --from <key_name>"+
		"```",
		username,
		hex.EncodeToString(callDataBz),
		time.Now().Add(time.Hour).UnixNano(),
	))
	return nil
}
