package main

import (
	"os"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/desmos-labs/desmos/v4/app"

	"github.com/desmos-labs/hephaestus/cmd"
)

func main() {
	// Setup the Cosmos SDK config
	app.SetupConfig(sdk.GetConfig())

	executor := cmd.RootCmd()
	if err := executor.Execute(); err != nil {
		os.Exit(1)
	}
}
