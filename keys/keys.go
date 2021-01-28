package keys

import (
	"os"
	"path"
)

var (
	HomeDir, _ = os.UserHomeDir()
	DataDir    = path.Join(HomeDir, ".hephaestus")
)

const (
	AppName = "hephaestus"

	ReactionWarning = "⚠️"
	ReactionDone    = "✅"
	ReactionTime    = "⌛"

	LogCommand = "command"
	LogUser    = "user"
)
