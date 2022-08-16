package bot

import (
	"time"

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
	bot.telegram.Handle(bot.getPrefixedCmd(types.CmdDocs), bot.HandleDocs)
	bot.telegram.Handle(bot.getPrefixedCmd(types.CmdHelp), bot.HandleHelp)
	bot.telegram.Handle(bot.getPrefixedCmd(types.CmdConnect), bot.HandleConnect)
	log.Debug().Msg("listening for messages...")
	bot.telegram.Start()
}

func (bot *Bot) getPrefixedCmd(cmd string) string {
	return bot.cfg.Prefix + cmd
}
