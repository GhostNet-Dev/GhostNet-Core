package states

import (
	"fmt"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type BlockMergeState struct {
	blockMachine *BlockMachine
	blockServer  IBlockServer
	glog         *glogger.GLogger
}

func (s *BlockMergeState) Initialize() {
	go s.MergeTask()
}

func (s *BlockMergeState) MergeTask() {
	s.blockMachine.MergeExecute()
	s.blockMachine.BlockServer.BlockServerInitStart()
	s.blockMachine.setState(s.blockMachine.miningState)
}

func (s *BlockMergeState) Rebuild() {

}

func (s *BlockMergeState) StartMining() {

}

func (s *BlockMergeState) RecvBlockHeight(height uint32, pubKey string) {

}

func (s *BlockMergeState) RecvBlockHash(from string, masterHash []byte, blockIdx uint32) {

}

func (s *BlockMergeState) RecvBlock(pairedBlock *types.PairedBlock, pubKey string) {
}

func (s *BlockMergeState) TimerExpired(context interface{}) bool {
	s.glog.DebugOutput(s, fmt.Sprint("Timeout - ", s), glogger.BlockConsensus)
	return false
}

func (s *BlockMergeState) Exit() {
}
