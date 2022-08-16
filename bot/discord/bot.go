package bot

import (
	"context"
	"time"

	"github.com/jasonlvhit/gocron"

	"github.com/desmos-labs/hephaestus/network"

	"github.com/desmos-labs/hephaestus/limitations"

	"github.com/andersfylling/disgord"
	"github.com/andersfylling/disgord/std"
	"github.com/rs/zerolog/log"

	"github.com/desmos-labs/hephaestus/types"
)

// Bot represents the object that should be used to interact with Discord
type Bot struct {
	cfg     *types.BotConfig
	discord *disgord.Client

	testnet *network.Client
	mainnet *network.Client
}

// Create allows to build a new Bot instance
func Create(cfg *types.BotConfig, testnet *network.Client, mainnet *network.Client) (*Bot, error) {
	// Set the default prefix if empty
	if cfg.Prefix == "" {
		cfg.Prefix = "!"
	}

	discordClient := disgord.New(disgord.Config{
		ProjectName: types.AppName,
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
		DMIntents: disgord.IntentDirectMessages |
			disgord.IntentDirectMessageReactions |
			disgord.IntentDirectMessageTyping,
		Presence: &disgord.UpdateStatusPayload{
			Game: &disgord.Activity{
				Name: "Welcome users!",
			},
		},
	})

	return &Bot{
		cfg:     cfg,
		discord: discordClient,
		testnet: testnet,
		mainnet: mainnet,
	}, nil
}

// Start starts the bot so that it can listen to events properly
func (bot *Bot) Start() {
	//nolint:errcheck
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
	handler.MessageCreate(
		bot.NewCmdHandler(types.CmdHelp, bot.HandleHelp),
		bot.NewCmdHandler(types.CmdDocs, bot.HandleDocs),
		bot.NewCmdHandler(types.CmdSend, bot.HandleSendTokens),
		bot.NewCmdHandler(types.CmdConnect, bot.HandleConnect),
		bot.NewCmdHandler(types.CmdVerify, bot.HandleVerify),
		bot.NewCmdHandler(types.CmdCheck, bot.HandleCheck),
	)

	// Setup periodic tasks
	log.Debug().Msg("setting up periodic tasks...")
	scheduler := gocron.NewScheduler()

	//nolint:errcheck
	scheduler.Every(5).Minutes().Do(bot.CleanRoles)

	scheduler.Start()

	log.Debug().Msg("listening for messages...")
}

// Reply sends a Discord message as a reply to the given msg
func (bot *Bot) Reply(msg *disgord.Message, s disgord.Session, message string) {
	_, _, err := msg.Author.SendMsg(context.Background(), s, &disgord.Message{
		Type:    disgord.MessageTypeDefault,
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
		log.Error().Err(err).Msg("failed to react to message")
	}
}

// Unreact allows to unreact with the provided emoji to the given message
func (bot *Bot) Unreact(msg *disgord.Message, s disgord.Session, emoji interface{}, flags ...disgord.Flag) {
	err := msg.Unreact(context.Background(), s, emoji, flags...)
	if err != nil {
		log.Error().Err(err).Msg("failed to unreact to message")
	}
}

// CheckCommandLimit returns the date on which the given user will be able to run the command again
func (bot *Bot) CheckCommandLimit(userID disgord.Snowflake, command string) *time.Time {
	// Try getting the expiration date for the command
	expirationDate, err := limitations.GetLimitationExpiration(userID.String(), command)
	if err != nil {
		panic(err)
	}

	// Check if the user is blocked
	if expirationDate != nil && time.Now().Before(*expirationDate) {
		log.Debug().Str(types.LogCommand, command).Time(types.LogExpirationEnd, *expirationDate).Msg("user is limited")
		return expirationDate
	}

	return nil
}

// SetCommandLimitation sets the limitation for the given user for the provided command
func (bot *Bot) SetCommandLimitation(userID disgord.Snowflake, cmd string) {
	// Set the expiration
	commandLimitation := bot.cfg.FindLimitationByCommand(cmd)
	if commandLimitation != nil {
		err := limitations.SetLimitationExpiration(userID, cmd, time.Now().Add(commandLimitation.Duration))
		if err != nil {
			log.Error().Err(err).Str(types.LogCommand, cmd).Msg("error while setting limitation expiration")
		}
	}
}
