package states

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type DownloadCheckState struct {
	blockMachine *BlockMachine
	glog         *glogger.GLogger
	reqBlockId   uint32
	startBlockId uint32
	// Set from Verification State
	targetBlockId uint32
}

func (s *DownloadCheckState) Initialize() {
	s.reqBlockId = 0
	s.startBlockId = 0
	s.targetBlockId = s.blockMachine.GetTargetHeight()
}

func (s *DownloadCheckState) Rebuild() {

}

func (s *DownloadCheckState) StartMining() {

}

func (s *DownloadCheckState) RecvBlockHeight(height uint32, pubKey string) {

}

func (s *DownloadCheckState) RecvBlockHash(from string, masterHash []byte, blockIdx uint32) {
	mutex := sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	if s.reqBlockId == 0 {
		s.reqBlockId = blockIdx
	} else if s.reqBlockId != blockIdx {
		return
	}
	s.glog.DebugOutput(s, fmt.Sprint("- recv verification Hash ", blockIdx), glogger.BlockConsensus)
	if s.startBlockId == 0 {
		s.startBlockId = blockIdx
	}
	hash := s.blockMachine.LoadHashFromTempDb(blockIdx)
	if hash == nil || !bytes.Equal(hash, masterHash) {
		s.blockMachine.BlockServer.RequestGetBlock(from, blockIdx)
	} else if blockIdx == s.targetBlockId {
		s.blockMachine.SetTargetHeight(s.targetBlockId)
		s.blockMachine.SetStartBlockId(s.startBlockId)
		s.blockMachine.setState(s.blockMachine.blockMergeState)
	} else {
		s.reqBlockId = blockIdx + 1
		s.blockMachine.BlockServer.RequestGetBlockHash(from, s.reqBlockId)
	}
}

func (s *DownloadCheckState) RecvBlock(pairedBlock *types.PairedBlock, pubKey string) {
	mutex := sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	s.glog.DebugOutput(s, fmt.Sprint("- recv GetBlock ", pairedBlock.BlockId(),
		" target block id = ", s.targetBlockId), glogger.BlockConsensus)

	result := s.blockMachine.CheckAndSave(pairedBlock)
	if !result {
		s.blockMachine.BlockServer.MergeErrorNotification(pubKey, result)
		s.glog.DebugOutput(s, fmt.Sprint("-- Merge Error", result), glogger.BlockConsensus)
		s.blockMachine.setState(s.blockMachine.getHeightestState)
		s.blockMachine.BlockServer.BroadcastBlockChainNotification()
	} else if pairedBlock.BlockId() == s.targetBlockId {
		s.blockMachine.SetTargetHeight(s.targetBlockId)
		s.blockMachine.SetStartBlockId(s.startBlockId)
		s.blockMachine.setState(s.blockMachine.blockMergeState)
	} else {
		s.reqBlockId = pairedBlock.BlockId() + 1
		s.blockMachine.BlockServer.RequestGetBlockHash(pubKey, s.reqBlockId)
	}
}

func (s *DownloadCheckState) TimerExpired(context interface{}) bool {
	return false
}

func (s *DownloadCheckState) Exit() {

}
