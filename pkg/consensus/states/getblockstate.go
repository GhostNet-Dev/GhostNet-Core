package states

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type GetBlockState struct {
	blockMachine *BlockMachine
}

func (s *GetBlockState) Inititalize() {
}

func (s *GetBlockState) Rebuild() {

}

func (s *GetBlockState) StartMining() {

}

func (s *GetBlockState) RecvBlockHeight(height uint32, pubKey string) {

}

func (s *GetBlockState) RecvBlockHash(from string, masterHash string, blockIdx uint32) {

}

func (s *GetBlockState) RecvBlock(pairedBlock *types.PairedBlock, pubKey string) {

}

func (s *GetBlockState) TimerExpired(context interface{}) bool {
	return false
}

func (s *GetBlockState) Exit() {

}
