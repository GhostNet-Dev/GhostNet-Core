package commands

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/spf13/cobra"
)

var (
	port     string
	host     string
	username string
	password string
)

// StartNodeCmd root command binding
func StartNodeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "startnode",
		Short: `Run a GhostNet Core Server for Distributed BlockChain Network`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			ExecuteNode()
		},
	}
	cmd.Flags().StringVarP(&port, "port", "", "50129", "Port Number")
	cmd.Flags().StringVarP(&host, "ip", "i", "", "Host Address")
	cmd.Flags().StringVarP(&username, "username", "u", "", "Ghost Account Nickname")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Ghost Account Password")
	return cmd
}

func ExecuteNode() {
	fmt.Printf("execute node port = %s\n", port)
	args := []string{
		"-port=" + port,
		"-ip=" + host,
		"-username=" + username,
		"-password=" + password,
	}
	execCmd := exec.Command(ghostDeamonName, args...)
	if err := execCmd.Start(); err != nil {
		log.Fatal(err)
	}
}
