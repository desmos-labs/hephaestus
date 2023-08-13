package bot

import (
	"strings"

	"github.com/andersfylling/disgord"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/hephaestus/types"
)

// AutoBanMessage bans the user if the message contains banned words
func (b *Bot) AutoBanMessage(s disgord.Session, data *disgord.MessageCreate) {
	if !b.containsBannedWords(data.Message.Content) {
		return
	}

	// Ban the user
	err := s.Guild(data.Message.GuildID).Member(data.Message.Author.ID).Ban(&disgord.BanMember{
		DeleteMessageDays: 7,
		Reason:            "Scam",
	})
	if err != nil {
		log.Error().Err(err).Str(types.LogUser, data.Message.Author.Username).Msg("error while banning user")
		return
	}
}

func (b *Bot) containsBannedWords(msg string) bool {
	for _, bannedWord := range b.cfg.BannedWords {
		if strings.Contains(strings.ToLower(msg), strings.ToLower(bannedWord)) {
			return true
		}
	}
	return false
}
