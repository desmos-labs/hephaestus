package utils

import (
	"fmt"

	"github.com/andersfylling/disgord"
)

// GetMsgAuthorUsername returns the username of the sender of the given message
func GetMsgAuthorUsername(msg *disgord.Message) string {
	return GetUserUsername(msg.Author)
}

// GetUserUsername returns the username of the given user
func GetUserUsername(user *disgord.User) string {
	return fmt.Sprintf("%s#%s", user.Username, user.Discriminator.String())
}
