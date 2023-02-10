package commands

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/grpc"
	"github.com/btcsuite/btcutil/base58"
	"github.com/spf13/cobra"
)

var blockId uint32 = 0

func GetBlockInfoCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "block",
		Short: `Ghost Account List`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			getBlockInfoExecute(id, blockId)
		},
	}
	cmd.Flags().StringVarP(&host, "ip", "i", "", "Host Address")
	cmd.Flags().StringVarP(&port, "port", "", "50129", "Port Number")
	cmd.Flags().StringVarP(&rpcPort, "rpc", "r", "50229", "GRPC Port Number")
	cmd.Flags().Uint32VarP(&id, "id", "", 0, "Container Id, if not select, show all container")
	cmd.Flags().Uint32VarP(&blockId, "bid", "", 0, "block id for information")
	return cmd
}

func getBlockInfoExecute(id, blockId uint32) {
	grpcClient := grpc.NewGrpcClient(host, rpcPort)
	grpcClient.ConnectServer()
	defer grpcClient.CloseServer()
	paired := grpcClient.GetBlockInfo(id, blockId)
	if paired != nil {
		log.Printf("[%d] blockId = %d, prev hash = %s\n",
			id, blockId, base58.CheckEncode(paired.Block.Header.PreviousBlockHeaderHash, 0))
	}
}
