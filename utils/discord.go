package utils

import (
	"fmt"

	"github.com/andersfylling/disgord"
)

// GetUsername returns the username of the sender of the given message
func GetUsername(msg *disgord.Message) string {
	return fmt.Sprintf("%s#%d", msg.Author.Username, msg.Author.Discriminator)
}
