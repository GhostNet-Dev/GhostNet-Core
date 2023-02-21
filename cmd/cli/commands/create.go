package commands

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/grpc"
	"github.com/spf13/cobra"
)

func CreateContainerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: `suspend a GhostNet Server`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			createExecuteCommand(username, password, host, port)
		},
	}
	cmd.Flags().StringVarP(&host, "ip", "i", "", "Host Address")
	cmd.Flags().StringVarP(&port, "port", "", "50129", "Port Number")
	cmd.Flags().StringVarP(&rpcPort, "rpc", "r", "50229", "GRPC Port Number")
	cmd.Flags().StringVarP(&username, "username", "u", "", "Ghost Account Nickname")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Ghost Account Password")
	cmd.Flags().Uint32VarP(&timeout, "timeout", "t", 8, "rpc connection timeout")
	cmd.MarkFlagRequired("username")
	return cmd
}

func createExecuteCommand(username, password, host, port string) {
	if username == "" {
		log.Println("ghostnet need username to login")
		return
	}
	grpcClient := grpc.NewGrpcClient(host, rpcPort, timeout)
	grpcClient.ConnectServer()
	defer grpcClient.CloseServer()
	log.Printf("Create Container user = %s, host = %s, port = %s", username, host, port)
	ret := grpcClient.CreateContainer(username, password, host, port)
	log.Printf("Result = %t", ret)
}
