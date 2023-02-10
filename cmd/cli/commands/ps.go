package commands

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/grpc"
	"github.com/spf13/cobra"
)

// StartNodeCmd root command binding
func PsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ps",
		Short: `Get List of Running a GhostNet Server`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			psExecuteCommand()
		},
	}
	cmd.Flags().StringVarP(&host, "ip", "i", "", "Host Address")
	cmd.Flags().StringVarP(&port, "port", "", "50129", "Port Number")
	cmd.Flags().StringVarP(&rpcPort, "rpc", "r", "50229", "GRPC Port Number")
	cmd.Flags().Uint32VarP(&id, "id", "", 0, "Container Id, if not select, show all container")
	return cmd
}

func psExecuteCommand() {
	grpcClient := grpc.NewGrpcClient(host, rpcPort)
	grpcClient.ConnectServer()
	defer grpcClient.CloseServer()
	// todo get all infomation
	if id == 0 {
		info := grpcClient.GetInfo()
		if info == nil {
			return
		}
		for i := uint32(0); i < info.TotalContainer; i++ {
			response := grpcClient.GetContainerList(i)
			log.Printf("[%d] Container User = %s, Port = %s, PID = %d",
				response.Id, response.Username, response.Port, response.PID)
		}
	} else {
		response := grpcClient.GetContainerList(id)
		log.Printf("[%d] Container User = %s, Port = %s, PID = %d",
			response.Id, response.Username, response.Port, response.PID)
	}
}
