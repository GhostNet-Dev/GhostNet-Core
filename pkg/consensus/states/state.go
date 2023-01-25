package states

import "github.com/GhostNet-Dev/GhostNet-Core/pkg/types"

type IBlockState interface {
	Initialize()
	Rebuild()
	StartMining()
	RecvBlockHeight(height uint32, pubKey string)
	RecvBlockHash(from string, masterHash []byte, blockIdx uint32)
	RecvBlock(pairedBlock *types.PairedBlock, pubKey string)
	TimerExpired(context interface{}) bool
	Exit()
}
