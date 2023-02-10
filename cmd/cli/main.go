package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	cmd "github.com/GhostNet-Dev/GhostNet-Core/cmd/cli/commands"
)

func init() {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	if err := os.Setenv("PATH", appendPaths(os.Getenv("PATH"), path)); err != nil {
		log.Fatal(err)
	}
}

func appendPaths(envPath string, path string) string {
	if envPath == "" {
		return path
	}
	return strings.Join([]string{envPath, path}, string(os.PathListSeparator))
}

func main() {
	startCmd := cmd.NewRootCommand()
	startCmd.AddCommand(
		cmd.NewGenesisCommand(),
		cmd.PsCommand(),
		cmd.ResumeCommand(),
		cmd.StartNodeCommand(),
		cmd.SuspendCommand(),
		cmd.UserAddCommand(),
		cmd.GetPrivateKeyCommand(),
		cmd.CreateContainerCommand(),
		cmd.GetBlockInfoCommand(),
		cmd.GetAccountListCommand(),
	)

	if err := startCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
