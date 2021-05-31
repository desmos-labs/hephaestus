package types

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

	LogCommand       = "command"
	LogUser          = "user"
	LogExpirationEnd = "expiration_end"
)
