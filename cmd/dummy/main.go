package main

import (
	"log"

	cmd "github.com/GhostNet-Dev/GhostNet-Core/cmd/dummy/commands"
)

func main() {
	startCmd := cmd.RootCommand()
	if err := startCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
