package states

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type GetHeigtestState struct {
	blockMachine *BlockMachine
}

func (s *GetHeigtestState) Inititalize() {
}

func (s *GetHeigtestState) Rebuild() {

}

func (s *GetHeigtestState) StartMining() {

}

func (s *GetHeigtestState) RecvBlockHeight(height uint32, pubKey string) {

}

func (s *GetHeigtestState) RecvBlockHash(from string, masterHash string, blockIdx uint32) {

}

func (s *GetHeigtestState) RecvBlock(pairedBlock *types.PairedBlock, pubKey string) {

}

func (s *GetHeigtestState) TimerExpired(context interface{}) bool {
	return false
}

func (s *GetHeigtestState) Exit() {

}
