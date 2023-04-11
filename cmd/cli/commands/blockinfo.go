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
		Short: `Ghost block information`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			grpcClient := grpc.NewGrpcClient(host, rpcPort, timeout)
			grpcClient.ConnectServer()
			defer grpcClient.CloseServer()
			paired := grpcClient.GetBlockInfo(id, blockId)
			if paired != nil {
				log.Printf("[%d] blockId = %d, prev hash = %s, Alice Count = %d, Tx Count = %d\n",
					id, blockId, base58.CheckEncode(paired.Block.Header.PreviousBlockHeaderHash, 0),
					paired.Block.Header.AliceCount, paired.Block.Header.TransactionCount)
			}
		},
	}
	cmd.Flags().StringVarP(&rpcPort, "rpc", "r", "50229", "GRPC Port Number")
	cmd.Flags().Uint32VarP(&id, "id", "", 0, "Container Id, if not select, show all container")
	cmd.Flags().Uint32VarP(&blockId, "bid", "", 1, "block id for information")
	return cmd
}

func GetBlockListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "blocklist",
		Short: `Ghost block list`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			grpcClient := grpc.NewGrpcClient(host, rpcPort, timeout)
			grpcClient.ConnectServer()
			defer grpcClient.CloseServer()
			for blockId = 1; ; blockId++ {
				paired := grpcClient.GetBlockInfo(id, blockId)
				if paired != nil {
					log.Printf("[%d] blockId = %d, prev hash = %s, Alice Count = %d, Tx Count = %d\n",
					id, blockId, base58.CheckEncode(paired.Block.Header.PreviousBlockHeaderHash, 0),
					paired.Block.Header.AliceCount, paired.Block.Header.TransactionCount)
				} else {
					break
				}
			}
		},
	}
	cmd.Flags().StringVarP(&rpcPort, "rpc", "r", "50229", "GRPC Port Number")
	cmd.Flags().Uint32VarP(&id, "id", "", 0, "Container Id, if not select, show all container")
	return cmd
}
