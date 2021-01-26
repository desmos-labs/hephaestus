package bot

import (
	"github.com/andersfylling/disgord"
)

// replyPongToPing is a handler that replies pong to ping messages
func (bot *Bot) replyPongToPing(s disgord.Session, data *disgord.MessageCreate) {
	msg := data.Message
	if msg.Content != "ping" {
		return
	}

	// whenever the message written is "ping", the hephaestus replies "pong"
	bot.Reply(msg, s, "pong")
}
