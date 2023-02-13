package commands

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/grpc"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/rpc"
	"github.com/spf13/cobra"
)

func ResumeCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "resume",
		Short: `stop a GhostNet Server`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			ctrlExecuteCommand(id, "resume")
		},
	}
	cmd.Flags().StringVarP(&host, "ip", "i", "", "Host Address")
	cmd.Flags().StringVarP(&port, "port", "", "50129", "Port Number")
	cmd.Flags().StringVarP(&rpcPort, "rpc", "r", "50229", "GRPC Port Number")
	cmd.Flags().Uint32VarP(&id, "id", "", 0, "Container Id, if not select, show all container")
	return cmd
}

func SuspendCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "suspend",
		Short: `suspend a GhostNet Server`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			ctrlExecuteCommand(id, "suspend")
		},
	}
	cmd.Flags().StringVarP(&host, "ip", "i", "", "Host Address")
	cmd.Flags().StringVarP(&port, "port", "", "50129", "Port Number")
	cmd.Flags().StringVarP(&rpcPort, "rpc", "r", "50229", "GRPC Port Number")
	cmd.Flags().Uint32VarP(&id, "id", "", 0, "Container Id, if not select, show all container")
	cmd.Flags().Uint32VarP(&timeout, "timeout", "t", 3, "rpc connection timeout")
	return cmd
}

func ctrlExecuteCommand(id uint32, operation string) {
	grpcClient := grpc.NewGrpcClient(host, rpcPort, timeout)
	grpcClient.ConnectServer()
	defer grpcClient.CloseServer()
	result := false
	switch operation {
	case "resume":
		result = grpcClient.ControlContainer(id, rpc.ContainerControlType_StartResume)
	case "suspend":
		result = grpcClient.ControlContainer(id, rpc.ContainerControlType_StopSuspend)
	}
	log.Printf("[%d] %s = %t", id, operation, result)
	// todo get all infomation
}
