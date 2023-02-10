package commands

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/grpc"
	"github.com/spf13/cobra"
)

func NewGenesisCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "newgen",
		Short: `suspend a GhostNet Server`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			newGenExecuteCommand(id, password)
		},
	}
	cmd.Flags().StringVarP(&host, "ip", "i", "", "Host Address")
	cmd.Flags().StringVarP(&port, "port", "", "50129", "Port Number")
	cmd.Flags().StringVarP(&rpcPort, "rpc", "r", "50229", "GRPC Port Number")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Ghost Account Password")
	cmd.Flags().Uint32VarP(&id, "id", "", 0, "Container Id, if not select, show all container")
	return cmd
}

func newGenExecuteCommand(id uint32, password string) {
	grpcClient := grpc.NewGrpcClient(host, rpcPort)
	grpcClient.ConnectServer()
	defer grpcClient.CloseServer()
	ret := grpcClient.CreateGenesisEncrypt(id, password)
	log.Printf("Create Genesis Block = %t", ret)
}
