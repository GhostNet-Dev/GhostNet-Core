package states

import (
	"fmt"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type RecheckState struct {
	blockMachine *BlockMachine
	glog         *glogger.GLogger
}

func (s *RecheckState) Initialize() {
	go s.BlockCheckTask()
}

func (s *RecheckState) BlockCheckTask() {
	if s.blockMachine.LocalBlockCheckProcess() {
		s.blockMachine.BlockServer.BlockServerInitStart()
		s.blockMachine.setState(s.blockMachine.miningState)
		s.glog.DebugOutput(s, fmt.Sprint("-- "), glogger.BlockConsensus)
	} else {
		s.blockMachine.setState(s.blockMachine.getHeightestState)
		s.blockMachine.BlockServer.BroadcastBlockChainNotification()
	}
}

func (s *RecheckState) Rebuild() {

}

func (s *RecheckState) StartMining() {

}

func (s *RecheckState) RecvBlockHeight(height uint32, pubKey string) {

}

func (s *RecheckState) RecvBlockHash(from string, masterHash []byte, blockIdx uint32) {

}

func (s *RecheckState) RecvBlock(pairedBlock *types.PairedBlock, pubKey string) {

}

func (s *RecheckState) TimerExpired(context interface{}) bool {
	return false
}

func (s *RecheckState) Exit() {

}
