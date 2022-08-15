package bot

import (
	"fmt"
	"regexp"

	"github.com/desmos-labs/hephaestus/types"

	"github.com/andersfylling/disgord"
)

var (
	addressRegex = regexp.MustCompile("desmos[0-9a-zA-Z]+")
)

// HandleCheck handles the check request from a user
func (bot *Bot) HandleCheck(s disgord.Session, data *disgord.MessageCreate) error {
	msg := data.Message

	// Get the address to check
	var address string
	if msg.ReferencedMessage != nil {
		address = addressRegex.FindString(msg.ReferencedMessage.Content)
	} else {
		address = addressRegex.FindString(msg.Content)
	}

	// Check the address validity
	if address == "" {
		return types.NewWarnErr("Missing address")
	}

	// Check the balance amount
	balance, err := bot.mainnet.GetBalance(address)
	if err != nil {
		return err
	}

	if balance != nil && !balance.IsZero() {
		bot.Reply(msg, s, fmt.Sprintf("User with addresss %s already has %s", address, balance.String()))
	} else {
		bot.Reply(msg, s, fmt.Sprintf("User with address %s has no DSM yet", address))
	}

	return nil
}
