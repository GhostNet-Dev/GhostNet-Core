package main

import (
	"log"

	cmd "github.com/GhostNet-Dev/GhostNet-Core/cmd/container/commands"
)

func main() {
	startCmd := cmd.RootCmd()
	if err := startCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
