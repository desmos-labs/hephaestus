package bot

import (
	"context"
	"time"

	"github.com/desmos-labs/discord-bot/config"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/discord-bot/cosmos"
	"github.com/desmos-labs/discord-bot/keys"
)

// Bot represents the object that should be used to interact with Discord
type Bot struct {
	cfg *config.BotConfig

	discord      *disgord.Client
	cosmosClient *cosmos.Client
}

// Create allows to build a new Bot instance
func Create(
	cfg *config.BotConfig, cosmosClient *cosmos.Client,
) (*Bot, error) {
	// Set the default prefix if empty
	if cfg.Prefix == "" {
		cfg.Prefix = "!"
	}

	discordClient := disgord.New(disgord.Config{
		ProjectName: keys.AppName,
		BotToken:    cfg.Token,
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
		cfg:          cfg,
		discord:      discordClient,
		cosmosClient: cosmosClient,
	}, nil
}

// Start starts the bot so that it can listen to events properly
func (bot *Bot) Start() {
	// nolint:errcheck
	defer bot.discord.Gateway().StayConnectedUntilInterrupted()

	log.Debug().Msg("starting bot")

	// Create a middleware that only accepts messages with a "ping" prefix
	// tip: use this to identify bot commands
	filter, _ := std.NewMsgFilter(context.Background(), bot.discord)
	filter.SetPrefix(bot.cfg.Prefix)

	handler := bot.discord.Gateway().
		WithMiddleware(
			filter.NotByBot,    // Ignore hephaestus messages
			filter.HasPrefix,   // Message must have the given prefix
			filter.StripPrefix, // Remove the command prefix from the message
		)
	handler.MessageCreate(bot.HandleSendTokens)

	log.Debug().Msg("listening for messages...")
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

func (bot *Bot) CheckCommandLimit(userID disgord.Snowflake, command string) *time.Time {
	// Try getting the expiration date for the command
	expirationDate, err := GetLimitationExpiration(userID, command)
	if err != nil {
		panic(err)
	}

	// Check if the user is blocked
	if expirationDate != nil && time.Now().Before(*expirationDate) {
		log.Debug().Str(keys.LogCommand, command).Time(LogExpirationEnd, *expirationDate).Msg("user is limited")
		return expirationDate
	}

	return nil
}

func (bot *Bot) SetCommandLimitation(userID disgord.Snowflake, cmd string) {
	// Set the expiration
	commandLimitation := bot.cfg.FindLimitationByCommand(cmd)
	if commandLimitation != nil {
		err := SetLimitationExpiration(userID, cmd, time.Now().Add(commandLimitation.Duration))
		if err != nil {
			log.Error().Err(err).Str(keys.LogCommand, cmd).Msg("error while setting limitation expiration")
		}
	}
}
