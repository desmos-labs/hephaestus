package main

import (
	"os"

	"github.com/desmos-labs/hephaestus/cmd"
)

func main() {
	executor := cmd.RootCmd()
	if err := executor.Execute(); err != nil {
		os.Exit(1)
	}
}
