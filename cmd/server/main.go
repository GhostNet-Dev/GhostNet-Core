package main

import (
	"os"

	cmd "github.com/GhostNet-Dev/GhostNet-Core/cmd/server/commands"
)

func main() {
	startCmd := cmd.RootCmd
	startCmd.AddCommand(cmd.StartCmdTest)
	if err := startCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
