package bot

import (
	"strconv"
	"time"

	"github.com/desmos-labs/hephaestus/limitations"
	"github.com/desmos-labs/hephaestus/network"
	"github.com/desmos-labs/hephaestus/types"
	"github.com/rs/zerolog/log"
	telebot "gopkg.in/telebot.v3"
)

// Bot represents the object that should be used to interact with Discord
type Bot struct {
	cfg      *types.BotConfig
	telegram *telebot.Bot

	testnet *network.Client
	mainnet *network.Client
}

// Create allows to build a new Bot instance
func Create(cfg *types.BotConfig, testnet *network.Client, mainnet *network.Client) (*Bot, error) {
	// Set the default prefix if empty
	if cfg.Prefix == "" {
		cfg.Prefix = "/"
	}
	bot, err := telebot.NewBot(telebot.Settings{
		Token:     cfg.Token,
		Poller:    &telebot.LongPoller{Timeout: 10 * time.Second},
		ParseMode: telebot.ModeMarkdown,
	})
	return &Bot{
		cfg:      cfg,
		telegram: bot,
		testnet:  testnet,
		mainnet:  mainnet,
	}, err
}

// Start starts the bot so that it can listen to events properly
func (bot *Bot) Start() {
	log.Debug().Msg("starting bot")
	bot.Handle(types.CmdDocs, bot.HandleDocs)
	bot.Handle(types.CmdHelp, bot.HandleHelp)
	bot.Handle(types.CmdConnect, bot.HandleConnect)
	bot.Handle(types.CmdSend, bot.HandleSendTokens)
	log.Debug().Msg("listening for messages...")
	bot.telegram.Start()
}

// CheckCommandLimit returns the date on which the given user will be able to run the command again
func (bot *Bot) CheckCommandLimit(userID int64, command string) *time.Time {
	// Try getting the expiration date for the command
	expirationDate, err := limitations.GetLimitationExpiration(strconv.FormatInt(userID, 10), command)
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
func (bot *Bot) SetCommandLimitation(userID int64, cmd string) {
	// Set the expiration
	commandLimitation := bot.cfg.FindLimitationByCommand(cmd)
	if commandLimitation != nil {
		err := limitations.SetLimitationExpiration(strconv.FormatInt(userID, 10), cmd, time.Now().Add(commandLimitation.Duration))
		if err != nil {
			log.Error().Err(err).Str(types.LogCommand, cmd).Msg("error while setting limitation expiration")
		}
	}
}
