package bot

import (
	"context"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/discord-bot/consts"
	"github.com/desmos-labs/discord-bot/cosmos"
)

// Bot represents the object that should be used to interact with Discord
type Bot struct {
	prefix string

	discord      *disgord.Client
	cosmosClient *cosmos.Client
}

// Create allows to build a new Bot instance
func Create(prefix string, token string, cosmosClient *cosmos.Client) (*Bot, error) {
	discordClient := disgord.New(disgord.Config{
		ProjectName: consts.AppName,
		BotToken:    token,
		RejectEvents: []string{
			// Rarely used, and causes unnecessary spam
			disgord.EvtTypingStart,

			// These require special privilege
			// https://discord.com/developers/docs/topics/gateway#privileged-intents
			disgord.EvtPresenceUpdate,
			disgord.EvtGuildMemberAdd,
			disgord.EvtGuildMemberUpdate,
			disgord.EvtGuildMemberRemove,
		},
		Presence: &disgord.UpdateStatusPayload{
			Game: &disgord.Activity{
				Name: "Welcome users!",
			},
		},
	})

	return &Bot{
		prefix:       prefix,
		discord:      discordClient,
		cosmosClient: cosmosClient,
	}, nil
}

// Start starts the bot so that it can listen to events properly
func (bot *Bot) Start() {
	defer bot.discord.Gateway().StayConnectedUntilInterrupted()

	// Create a middleware that only accepts messages with a "ping" prefix
	// tip: use this to identify bot commands
	filter, _ := std.NewMsgFilter(context.Background(), bot.discord)
	filter.SetPrefix(bot.prefix)

	handler := bot.discord.Gateway().
		WithMiddleware(
			filter.NotByBot,    // Ignore hephaestus messages
			filter.HasPrefix,   // Message must have the given prefix
			filter.StripPrefix, // Remove the command prefix from the message
		)
	handler.MessageCreate(bot.replyPongToPing, bot.handleSendTokens)
}

// Reply sends a Discord message as a reply to the given msg
func (bot *Bot) Reply(msg *disgord.Message, s disgord.Session, message string) {
	_, err := msg.Reply(context.Background(), s, &disgord.CreateMessageParams{
		MessageReference: &disgord.MessageReference{
			MessageID: msg.ID,
			ChannelID: msg.ChannelID,
			GuildID:   msg.GuildID,
		},
		Content: message,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to reply to message")
	}
}

// React allows to react with the provided emoji to the given message
func (bot *Bot) React(msg *disgord.Message, s disgord.Session, emoji interface{}, flags ...disgord.Flag) {
	err := msg.React(context.Background(), s, emoji, flags...)
	if err != nil {
		log.Error().Err(err).Msg("failed to reply to message")
	}
}
