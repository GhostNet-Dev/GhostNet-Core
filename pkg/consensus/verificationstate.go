package consensus

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type VerificationState struct {
	blockMachine *BlockMachine
}

func (s *VerificationState) Inititalize() {
}

func (s *VerificationState) Rebuild() {

}

func (s *VerificationState) StartMining() {

}

func (s *VerificationState) RecvBlockHeight(height uint32, pubKey string) {

}

func (s *VerificationState) RecvBlockHash(from string, masterHash string, blockIdx uint32) {

}
func (s *VerificationState) RecvBlock(pairedBlock *types.PairedBlock, pubKey string) {

}

func (s *VerificationState) TimerExpired(context interface{}) bool {
	return false
}

func (s *VerificationState) Exit() {

}
