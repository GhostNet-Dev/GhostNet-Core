package states

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type IdleState struct {
	blockMachine *BlockMachine
	glog         *glogger.GLogger
}

func (s *IdleState) Initialize() {
}

func (s *IdleState) Rebuild() {
	s.blockMachine.blockServer.MiningStop()
	s.blockMachine.setState(s.blockMachine.getHeightestState)
	s.blockMachine.blockServer.BroadcastBlockChainNotification()
	s.glog.DebugOutput(s, "Rebuild", glogger.BlockConsensus)
}

func (s *IdleState) StartMining() {
	s.blockMachine.blockServer.BlockServerInitStart()
	s.blockMachine.setState(s.blockMachine.getHeightestState)
	s.blockMachine.blockServer.BroadcastBlockChainNotification()
	s.glog.DebugOutput(s, "StartMining", glogger.BlockConsensus)
}

func (s *IdleState) RecvBlockHeight(height uint32, pubKey string) {
	if s.blockMachine.blockServer.CheckHeightForRebuild(height) == true {
		s.Rebuild()
	}
}

func (s *IdleState) RecvBlockHash(from string, masterHash []byte, blockIdx uint32) {

}

func (s *IdleState) RecvBlock(pairedBlock *types.PairedBlock, pubKey string) {

}

func (s *IdleState) TimerExpired(context interface{}) bool {
	return false
}

func (s *IdleState) Exit() {

}
