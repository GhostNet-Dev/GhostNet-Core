package main

import (
	"os"

	cmd "github.com/GhostNet-Dev/GhostNet-Core/cmd/container/commands"
)

func main() {
	startCmd := cmd.RootCmd()
	if err := startCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
