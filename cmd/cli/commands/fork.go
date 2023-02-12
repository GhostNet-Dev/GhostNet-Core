package commands

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/grpc"
	"github.com/spf13/cobra"
)

func ForkContainerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fork",
		Short: `fork container a GhostNet Server`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			forkExecuteCommand(username, password, host, port)
		},
	}
	cmd.Flags().StringVarP(&host, "ip", "i", "", "Host Address")
	cmd.Flags().StringVarP(&port, "port", "", "50129", "Port Number")
	cmd.Flags().StringVarP(&rpcPort, "rpc", "r", "50229", "GRPC Port Number")
	cmd.Flags().StringVarP(&username, "username", "u", "", "Ghost Account Nickname")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Ghost Account Password")
	cmd.MarkFlagRequired("username")
	return cmd
}

func forkExecuteCommand(username, password, host, port string) {
	if username == "" {
		log.Println("ghostnet need username to login")
		return
	}
	grpcClient := grpc.NewGrpcClient(host, rpcPort)
	grpcClient.ConnectServer()
	defer grpcClient.CloseServer()
	log.Printf("Fork Container user = %s, host = %s, port = %s", username, host, port)
	ret := grpcClient.ForkContainer(username, password, host, port)
	log.Printf("Result = %t", ret)
}
