package states

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type MiningState struct {
	blockMachine *BlockMachine
	glog         *glogger.GLogger
}

func (s *MiningState) Initialize() {
}

func (s *MiningState) Rebuild() {
	s.blockMachine.BlockServer.MiningStop()
	s.blockMachine.setState(s.blockMachine.getHeightestState)
	s.blockMachine.BlockServer.BroadcastBlockChainNotification()
	s.glog.DebugOutput(s, "MiningState: Rebuild", glogger.BlockConsensus)
}

func (s *MiningState) StartMining() {

}

func (s *MiningState) RecvBlockHeight(height uint32, pubKey string) {
	if s.blockMachine.BlockServer.CheckHeightForRebuild(height) {
		s.Rebuild()
	}
}

func (s *MiningState) RecvBlockHash(from string, masterHash []byte, blockIdx uint32) {

}

func (s *MiningState) RecvBlock(pairedBlock *types.PairedBlock, pubKey string) {

}

func (s *MiningState) TimerExpired(context interface{}) bool {
	return false
}

func (s *MiningState) Exit() {

}
