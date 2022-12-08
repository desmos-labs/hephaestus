package bot

import (
	"encoding/hex"
	"fmt"
	"strings"
	"time"

	"github.com/tendermint/tendermint/libs/json"

	"github.com/desmos-labs/hephaestus/utils"

	"github.com/andersfylling/disgord"

	"github.com/desmos-labs/hephaestus/types"
)

type CallData struct {
	// Username is the plain-text Discord username of the user that wants to be verified
	Username string `json:"username"`
}

// NewCallData returns a new CallData instance
func NewCallData(username string) *CallData {
	return &CallData{
		Username: username,
	}
}

// HandleConnect handle a connection request. This request is done by the users when they want to connect their Desmos
// profile with their Discord account.
// The command expects one single argument which must be the JSON object returned from the "desmos sign" command.
//
// The handling of the command will fail in the following occasions:
// 1. The signed value does not correspond to the username of the user sending the message
// 2. Any of the values are badly encoded
func (bot *Bot) HandleConnect(s disgord.Session, data *disgord.MessageCreate) error {
	// Get the arguments
	msg := data.Message
	parts := strings.SplitN(msg.Content, " ", 3)[1:]
	if len(parts) == 0 {
		bot.Reply(msg, s, fmt.Sprintf(`**Connect**
This command allows you to connect your Discord account to your Desmos profile.
To do this, you have to: 

1. Sign your Discord username using the Desmos CLI or any Desmos-compatible application.
2. Use the %[1]s command to send the signature result. 

__Signing your Discord username__
1. Copy your Discord username by clicking on it in the bottom part of your Discord client. 
   Your full username should be in the form of <username>#<identifier> (eg. Foo#123).

2. Open your Desmos CLI or application, and sign your username. 
   If you use the Desmos CLI, you can do this by using the following command:
   `+"`desmos sign <Discord username> --from <your-key>`"+`
   
   Eg. `+"`desmos sign \"Foo#123\" --from foo`"+`

__Sending the signed value__
The sign command should return a JSON object. The last thing you have to do is now send it to me using the %[1]s command. To do this, simply send me a message as the following: 
`+"`!%[1]s <%[2]s/%[3]s> <JSON>`"+`

Eg. `+"`!%[1]s %[2]s {...}`"+`
`, types.CmdConnect, types.NetworkTestnet, types.NetworkMainnet))
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

	// Get the signature data
	username := utils.GetMsgAuthorUsername(msg)
	signatureData, err := utils.GetSignatureData(parts[1])
	if err != nil {
		return err
	}

	// Upload the data to Themis
	err = networkClient.UploadDataToThemis(username, bot.cfg.Name, signatureData)
	if err != nil {
		return err
	}

	// Return to the user the call data for the Desmos command
	callDataBz, err := json.Marshal(NewCallData(username))
	if err != nil {
		return types.NewWarnErr("Error while serializing call data: %s", err)
	}

	bot.Reply(msg, s, fmt.Sprintf("Your verification data has been stored successfully. "+
		"All you have to do now is execute the following command:\n"+
		"```"+
		"desmos tx profiles link-app ibc-profiles [channel] discord \"%[1]s\" %[2]s --packet-timeout-height 0-0 --packet-timeout-timestamp %[3]d --from <key_name>"+
		"```",
		username,
		hex.EncodeToString(callDataBz),
		time.Now().Add(time.Hour).UnixNano(),
	))

	return nil
}
