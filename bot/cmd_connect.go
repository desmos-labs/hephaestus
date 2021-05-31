package bot

import (
	"fmt"
	"strings"

	"github.com/andersfylling/disgord"

	"github.com/desmos-labs/hephaestus/themis"
	"github.com/desmos-labs/hephaestus/types"
)

const (
	ConnectCmd = "connect"
)

// HandleConnect handle a connection request. This request is done by the users when they want to connect their Desmos
// profile with their Discord account.
// The command expects four arguments:
// 1. The hex-encoded address that the user has on Desmos
// 2. The hex-encoded public key associated with the Desmos address
// 3. The plain text value that has been signed using the private key associated with the public key
// 4. The hex-encoded signature obtained by signing the plain text value with the user private key
//
// The handling of the command will fail in the following occasions:
// 1. The signed value does not correspond to the username of the user sending the message
// 2. Any of the values are badly encoded
func (bot *Bot) HandleConnect(s disgord.Session, data *disgord.MessageCreate) error {
	// Get the arguments
	msg := data.Message
	parts := strings.Split(msg.Content, " ")
	if len(parts) != 4 {
		return types.NewWarnErr("Invalid number of arguments")
	}

	username := fmt.Sprintf("%s#%d", msg.Author.Username, msg.Author.Discriminator)

	connectionData := types.NewConnectionData(parts[1], parts[2], username, parts[3])
	err := connectionData.Validate()
	if err != nil {
		return err
	}

	err = themis.UploadData(connectionData, bot.themisCfg.Host, bot.privKey)
	if err != nil {
		return err
	}

	bot.Reply(msg, s, fmt.Sprintf("Your verification data has been stored successfully. "+
		"All you have to do now is execute the following command:\n"+
		"```desmos tx profiles connect discord \"%s\" --from <key_name>```",
		connectionData.Username,
	))

	return nil
}
