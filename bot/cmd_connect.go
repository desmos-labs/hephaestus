package bot

import (
	"encoding/hex"
	"fmt"
	"strings"

	signcmd "github.com/desmos-labs/desmos/v2/app/desmos/cmd/sign"
	"github.com/tendermint/tendermint/crypto/secp256k1"
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
	parts := strings.Split(msg.Content, " ")[1:]
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

	// Get the signature data
	username := utils.GetUsername(msg)
	signatureData, err := bot.getSignatureData(parts[1], username)
	if err != nil {
		return err
	}

	// Upload the data to Themis
	err = networkClient.UploadDataToThemis(username, signatureData)
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
		"desmos tx profiles link-app ibc-profiles [channel] discord \"%[1]s\" %[2]s --from <key_name>"+
		"```",
		username,
		hex.EncodeToString(callDataBz),
	))

	return nil
}

func (bot *Bot) getSignatureData(jsonData string, username string) (*signcmd.SignatureData, error) {
	var signatureData signcmd.SignatureData
	err := json.Unmarshal([]byte(jsonData), &signatureData)
	if err != nil {
		return nil, types.NewWarnErr("Invalid data provided: %s", err)
	}

	// Verify the signature
	pubKeyBz, err := hex.DecodeString(signatureData.PubKey)
	if err != nil {
		return nil, types.NewWarnErr("Error while reading public key: %s", err)
	}

	valueBz, err := hex.DecodeString(signatureData.Value)
	if err != nil {
		return nil, types.NewWarnErr("Error while reading value: %s", err)
	}

	sigBz, err := hex.DecodeString(signatureData.Signature)
	if err != nil {
		return nil, types.NewWarnErr("Error while reading signature: %s", err)
	}

	pubKey := secp256k1.PubKey(pubKeyBz)
	if !pubKey.VerifySignature(valueBz, sigBz) {
		return nil, types.NewWarnErr("Invalid signature. Make sure you have signed the message using the correct account")
	}

	return &signatureData, nil
}
