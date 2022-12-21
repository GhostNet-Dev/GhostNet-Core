package states

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type MiningState struct {
	blockMachine *consensus.BlockMachine
}

func (s *MiningState) Inititalize() {
}

func (s *MiningState) Rebuild() {

}

func (s *MiningState) StartMining() {

}

func (s *MiningState) RecvBlockHeight(height uint32, pubKey string) {

}

func (s *MiningState) RecvBlockHash(from string, masterHash string, blockIdx uint32) {

}
func (s *MiningState) RecvBlock(pairedBlock *types.PairedBlock, pubKey string) {

}

func (s *MiningState) Exit() {

}
