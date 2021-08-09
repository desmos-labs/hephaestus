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

	CmdHelp    = "help"
	CmdConnect = "connect"
	CmdDocs    = "docs"
	CmdSend    = "send"
	CmdVerify  = "verify"

	LogCommand       = "command"
	LogUser          = "user"
	LogExpirationEnd = "expiration_end"
)
