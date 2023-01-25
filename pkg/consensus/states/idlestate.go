package states

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type IdleState struct {
	blockMachine *BlockMachine
}

func (idle *IdleState) Initialize() {
}

func (idle *IdleState) Rebuild() {
	idle.blockMachine.blockServer.MiningStop()
	idle.blockMachine.setState(idle.blockMachine.getHeightestState)
	idle.blockMachine.blockServer.BroadcastBlockChainNotification()
	glogger.DebugOutput(idle, "Rebuild", glogger.BlockConsensus)
}

func (idle *IdleState) StartMining() {
	idle.blockMachine.blockServer.BlockServerInitStart()
	idle.blockMachine.setState(idle.blockMachine.getHeightestState)
	idle.blockMachine.blockServer.BroadcastBlockChainNotification()
	glogger.DebugOutput(idle, "StartMining", glogger.BlockConsensus)
}

func (idle *IdleState) RecvBlockHeight(height uint32, pubKey string) {
	if idle.blockMachine.blockServer.CheckHeightForRebuild(height) == true {
		idle.Rebuild()
	}
}

func (idle *IdleState) RecvBlockHash(from string, masterHash []byte, blockIdx uint32) {

}

func (idle *IdleState) RecvBlock(pairedBlock *types.PairedBlock, pubKey string) {

}

func (idle *IdleState) TimerExpired(context interface{}) bool {
	return false
}

func (idle *IdleState) Exit() {

}
