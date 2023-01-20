package states

import "github.com/GhostNet-Dev/GhostNet-Core/pkg/types"

type BlockState interface {
	Rebuild()
	StartMining()
	RecvBlockHeight(height uint32, pubKey string)
	RecvBlockHash(from string, masterHash string, blockIdx uint32)
	RecvBlock(pairedBlock *types.PairedBlock, pubKey string)
	TimerExpired(context interface{}) bool
	Exit()
}
