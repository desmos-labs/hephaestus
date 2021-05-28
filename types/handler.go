package types

import (
	"github.com/andersfylling/disgord"
)

// CmdHandler represents a function that extends a disgord.HandlerMessageCreate to allow it to return an error
type CmdHandler = func(s disgord.Session, h *disgord.MessageCreate) error

// MergeHandlers merges all the given handlers into a single one
func MergeHandlers(handlers ...CmdHandler) CmdHandler {
	return func(s disgord.Session, data *disgord.MessageCreate) error {
		for _, h := range handlers {
			err := h(s, data)
			if err != nil {
				return err
			}
		}
		return nil
	}
}
