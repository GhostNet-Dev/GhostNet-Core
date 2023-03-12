package states

import (
	"bytes"
	"fmt"
	"sync"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type VerificationState struct {
	blockMachine      *BlockMachine
	lock              *sync.Mutex
	glog              *glogger.GLogger
	lastestReqBlockId uint32
}

func (s *VerificationState) Initialize() {
}

func (s *VerificationState) Rebuild() {

}

func (s *VerificationState) StartMining() {

}

func (s *VerificationState) RecvBlockHeight(height uint32, pubKey string) {

}

func (s *VerificationState) RecvBlockHash(from string, masterHash []byte, blockIdx uint32) {
	s.lock.Lock()
	masterList := s.blockMachine.GetHeighestCandidatePool()
	exist := false
	for _, master := range masterList {
		if master == from {
			exist = true
		}
	}
	if exist == false {
		return
	}
	header, _ := s.blockMachine.blockContainer.GetBlockHeader(blockIdx)
	if header == nil {
		return
	}
	if bytes.Compare(header.GetHashKey(), masterHash) != 0 {
		s.glog.DebugOutput(s, fmt.Sprint("- recv verification Hash ", blockIdx), glogger.BlockConsensus)
		s.lastestReqBlockId = blockIdx - 1
		s.blockMachine.blockServer.RequestGetBlockHash(from, s.lastestReqBlockId)
	} else {
		s.glog.DebugOutput(s, fmt.Sprint("- recv verification Find ", blockIdx), glogger.BlockConsensus)
		s.blockMachine.setState(s.blockMachine.downloadCheckState)
		s.blockMachine.blockServer.RequestGetBlockHash(from, blockIdx+1)
	}

	s.lock.Unlock()
}
func (s *VerificationState) RecvBlock(pairedBlock *types.PairedBlock, pubKey string) {

}

func (s *VerificationState) TimerExpired(context interface{}) bool {
	return false
}

func (s *VerificationState) Exit() {

}
