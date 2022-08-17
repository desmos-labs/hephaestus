package bot

import (
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
	"gopkg.in/telebot.v3"

	"github.com/desmos-labs/hephaestus/types"
)

func (bot *Bot) Handle(cmd string, h telebot.HandlerFunc, m ...telebot.MiddlewareFunc) {
	bot.telegram.Handle(bot.getPrefixedCmd(cmd), bot.NewCmdHandler(cmd, h), m...)
}

func (bot *Bot) getPrefixedCmd(cmd string) string {
	return bot.cfg.Prefix + cmd
}

// MergeHandlers merges all the given handlers into a single one
func MergeHandlers(handlers ...telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		for _, h := range handlers {
			err := h(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// NewCmdHandler returns a new command handler for the command that has the given name
func (bot *Bot) NewCmdHandler(cmdName string, handler telebot.HandlerFunc) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		log.Debug().Str(types.LogCommand, cmdName).Msg("received command")

		// Merge the handler with the limit check
		mergedHandlers := MergeHandlers(bot.checkCmdLimit(cmdName), handler)

		// Handle the message
		err := mergedHandlers(ctx)
		if err != nil {
			bot.handleError(ctx, err)
		}
		//TODO: react with message with ✅ when it is supported
		bot.SetCommandLimitation(ctx.Sender().ID, cmdName)
		return nil
	}
}

func (bot *Bot) handleError(ctx telebot.Context, err error) {
	log.Warn().Err(err).Str(types.LogUser, ctx.Sender().Username).Msg("error while handling command")

	customErr, ok := err.(*types.Error)
	if ok {
		ctx.Reply(customErr.Description)
		//TODO: using ctx.React(customErr.Reaction) when it is supported
	} else {
		ctx.Reply(err.Error())
		//TODO: using ctx.React("🚨") when it is supported
	}
}

func (bot *Bot) checkCmdLimit(cmdName string) telebot.HandlerFunc {
	return func(ctx telebot.Context) error {
		// Check the command limitation
		if expirationDate := bot.CheckCommandLimit(ctx.Sender().ID, cmdName); expirationDate != nil {
			return types.NewTimeErr(fmt.Sprintf(
				"Command limit reached. You will be able to use this command again on: %s",
				expirationDate.Format(time.RFC822),
			))
		}
		return nil
	}
}
