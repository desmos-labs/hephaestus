package bot

import (
	"encoding/hex"
	"fmt"
	"strings"

	desmoscmd "github.com/desmos-labs/desmos/app/desmos/cmd"
	"github.com/tendermint/tendermint/crypto/secp256k1"
	"github.com/tendermint/tendermint/libs/json"

	"github.com/desmos-labs/hephaestus/utils"

	"github.com/andersfylling/disgord"

	"github.com/desmos-labs/hephaestus/themis"
	"github.com/desmos-labs/hephaestus/types"
)

type CallData struct {
	Username string `json:"username"`
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
	content := strings.TrimSpace(strings.TrimPrefix(msg.Content, types.CmdConnect))
	if content == "" {
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
`+"`!%[1]s <JSON>`"+`

Eg. `+"`!%[1]s connect {...}`"+`
`, types.CmdConnect))
		return nil
	}

	var signatureData desmoscmd.SignatureData
	err := json.Unmarshal([]byte(content), &signatureData)
	if err != nil {
		return types.NewWarnErr("Invalid data provided: %s", err)
	}

	// Verify the username
	username := utils.GetUsername(msg)
	if signatureData.Value != username {
		return types.NewWarnErr("Invalid signed value. Make sure you sign your username (%s)", username)
	}

	// Verify the signature
	pubKeyBz, err := hex.DecodeString(signatureData.PubKey)
	if err != nil {
		return types.NewWarnErr("Error while reading public key: %s", err)
	}

	sigBz, err := hex.DecodeString(signatureData.Signature)
	if err != nil {
		return types.NewWarnErr("Error while reading signature: %s", err)
	}

	pubKey := secp256k1.PubKey(pubKeyBz)
	if !pubKey.VerifySignature([]byte(signatureData.Value), sigBz) {
		return types.NewWarnErr("Invalid signature. Make sure you have signed your username (%s)", username)
	}

	// Upload the data to Themis
	connectionData := types.NewConnectionData(signatureData.Address, signatureData.PubKey, signatureData.Value, signatureData.Signature)
	err = connectionData.Validate()
	if err != nil {
		return err
	}

	err = themis.UploadData(connectionData, bot.themisCfg.Host, bot.privKey)
	if err != nil {
		return err
	}

	callData := CallData{Username: signatureData.Value}
	callDataBz, err := json.Marshal(&callData)
	if err != nil {
		return types.NewWarnErr("Error while serializing call data: %s", err)
	}

	bot.Reply(msg, s, fmt.Sprintf("Your verification data has been stored successfully. "+
		"All you have to do now is execute the following command:\n"+
		"```"+
		"desmos tx profiles link-app ibc-profiles [channel] discord \"%[1]s\" %[2]s --from <key_name>"+
		"```",
		connectionData.Username,
		hex.EncodeToString(callDataBz),
	))

	return nil
}
