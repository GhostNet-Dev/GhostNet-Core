package consensus

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type IdleState struct {
	blockMachine *BlockMachine
}

func (idle *IdleState) Inititalize() {
}

func (idle *IdleState) Rebuild() {

}

func (idle *IdleState) StartMining() {

}

func (idle *IdleState) RecvBlockHeight(height uint32, pubKey string) {

}

func (idle *IdleState) RecvBlockHash(from string, masterHash string, blockIdx uint32) {

}

func (idle *IdleState) RecvBlock(pairedBlock *types.PairedBlock, pubKey string) {

}

func (idle *IdleState) TimerExpired(context interface{}) bool {
	return false
}

func (idle *IdleState) Exit() {

}
