package states

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type DownloadCheckState struct {
	blockMachine  *BlockMachine
	lock          *sync.Mutex
	reqBlockId    uint32
	startBlockId  uint32
	targetBlockId uint32
}

func (s *DownloadCheckState) Initialize() {
	s.reqBlockId = 0
	s.startBlockId = 0
	s.targetBlockId = 0
}

func (s *DownloadCheckState) Rebuild() {

}

func (s *DownloadCheckState) StartMining() {

}

func (s *DownloadCheckState) RecvBlockHeight(height uint32, pubKey string) {

}

func (s *DownloadCheckState) RecvBlockHash(from string, masterHash []byte, blockIdx uint32) {
	s.lock.Lock()
	if s.reqBlockId == 0 {
		s.reqBlockId = blockIdx
	} else {
		return
	}
	glogger.DebugOutput(s, fmt.Sprint("- recv verification Hash ", blockIdx), glogger.BlockConsensus)
	if s.startBlockId == 0 {
		s.startBlockId = blockIdx
	}
	hash := s.blockMachine.LoadHashFromTempDb(blockIdx)
	if hash == nil || bytes.Compare(hash, masterHash) != 0 {
		s.blockMachine.blockServer.RequestGetBlock(from, blockIdx)
	} else {
		s.reqBlockId = blockIdx + 1
		s.blockMachine.blockServer.RequestGetBlockHash(from, s.reqBlockId)
	}
	s.lock.Unlock()
}

func (s *DownloadCheckState) RecvBlock(pairedBlock *types.PairedBlock, pubKey string) {
	s.lock.Lock()
	glogger.DebugOutput(s, fmt.Sprint("- recv GetBlock ", pairedBlock.BlockId()), glogger.BlockConsensus)

	result := s.blockMachine.CheckAndSave(pairedBlock)
	if result == false {
		s.blockMachine.blockServer.MergeErrorNotification(pubKey, result)
		glogger.DebugOutput(s, fmt.Sprint("-- Merge Error", result), glogger.BlockConsensus)
		s.blockMachine.setState(s.blockMachine.getHeightestState)
		s.blockMachine.blockServer.BroadcastBlockChainNotification()
	} else if pairedBlock.BlockId() == s.targetBlockId {
		s.blockMachine.SetTargetHeight(s.targetBlockId)
		s.blockMachine.SetStartBlockId(s.startBlockId)
		s.blockMachine.setState(s.blockMachine.blockMergeState)
	} else {
		s.reqBlockId = pairedBlock.BlockId() + 1
		s.blockMachine.blockServer.RequestGetBlockHash(pubKey, s.reqBlockId)
	}

	s.lock.Unlock()
}

func (s *DownloadCheckState) TimerExpired(context interface{}) bool {
	return false
}

func (s *DownloadCheckState) Exit() {

}
