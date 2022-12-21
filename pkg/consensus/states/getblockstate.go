package states

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type GetBlockState struct {
	blockMachine *consensus.BlockMachine
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

func (s *GetBlockState) Exit() {

}
