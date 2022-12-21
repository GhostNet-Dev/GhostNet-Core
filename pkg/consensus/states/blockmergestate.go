package states

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type BlockMergeState struct {
	blockMachine *consensus.BlockMachine
}

func (s *BlockMergeState) Inititalize() {
}

func (s *BlockMergeState) Rebuild() {

}

func (s *BlockMergeState) StartMining() {

}

func (s *BlockMergeState) RecvBlockHeight(height uint32, pubKey string) {

}

func (s *BlockMergeState) RecvBlockHash(from string, masterHash string, blockIdx uint32) {

}
func (s *BlockMergeState) RecvBlock(pairedBlock *types.PairedBlock, pubKey string) {

}

func (s *BlockMergeState) Exit() {

}
