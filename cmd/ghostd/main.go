package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	cmd "github.com/GhostNet-Dev/GhostNet-Core/cmd/ghostd/commands"
	//"github.com/sirupsen/logrus"
)

func init() {
	path, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatal(err)
	}

	if err := os.Setenv("PATH", appendPaths(os.Getenv("PATH"), path)); err != nil {
		log.Fatal(err)
	}
	/*
		// Log as JSON instead of the default ASCII formatter.
		logrus.SetFormatter(&logrus.JSONFormatter{})

		// Output to stdout instead of the default stderr
		// Can be any io.Writer, see below for File example
		logrus.SetOutput(os.Stdout)

		// Only log the warning severity or above.
		logrus.SetLevel(logrus.WarnLevel)
	*/
}

func appendPaths(envPath string, path string) string {
	if envPath == "" {
		return path
	}
	return strings.Join([]string{envPath, path}, string(os.PathListSeparator))
}

func main() {
	startCmd := cmd.RootCommand()
	if err := startCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
