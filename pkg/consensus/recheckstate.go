package consensus

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type RecheckState struct {
	blockMachine *BlockMachine
}

func (s *RecheckState) Inititalize() {
}

func (s *RecheckState) Rebuild() {

}

func (s *RecheckState) StartMining() {

}

func (s *RecheckState) RecvBlockHeight(height uint32, pubKey string) {

}

func (s *RecheckState) RecvBlockHash(from string, masterHash string, blockIdx uint32) {

}

func (s *RecheckState) RecvBlock(pairedBlock *types.PairedBlock, pubKey string) {

}

func (s *RecheckState) TimerExpired(context interface{}) bool {
	return false
}

func (s *RecheckState) Exit() {

}
