package states

import (
	"fmt"
	"sync"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type GetHeigtestState struct {
	blockMachine  *BlockMachine
	glog          *glogger.GLogger
	maxHeight     uint32
	selectNode    string
	candidatePool map[uint32][]string
}

func (s *GetHeigtestState) Initialize() {
	s.candidatePool = make(map[uint32][]string)
	go s.TimerExpired(nil)
}

func (s *GetHeigtestState) Rebuild() {

}

func (s *GetHeigtestState) StartMining() {

}

func (s *GetHeigtestState) RecvBlockHeight(height uint32, pubKey string) {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()

	if s.maxHeight < height {
		s.maxHeight = height
		s.selectNode = pubKey
		if node, exist := s.candidatePool[height]; !exist {
			s.candidatePool[height] = []string{pubKey}
		} else {
			s.candidatePool[height] = append(node, pubKey)
		}
	}
	s.glog.DebugOutput(s, fmt.Sprint(height), glogger.BlockConsensus)
}

func (s *GetHeigtestState) RecvBlockHash(from string, masterHash []byte, blockIdx uint32) {

}

func (s *GetHeigtestState) RecvBlock(pairedBlock *types.PairedBlock, pubKey string) {

}

func (s *GetHeigtestState) TimerExpired(context interface{}) bool {
	<-time.After(time.Second * 8)

	curBlockHeight := s.blockMachine.blockContainer.BlockHeight()
	if s.maxHeight > curBlockHeight && s.maxHeight != 0 {
		candiList := s.candidatePool[s.maxHeight]
		s.blockMachine.SetHeighestCandidatePool(candiList)
		s.blockMachine.SetTargetHeight(s.maxHeight)
		s.blockMachine.setState(s.blockMachine.verificationState)
		s.blockMachine.BlockServer.RequestGetBlockHash(candiList[0], curBlockHeight)
		s.glog.DebugOutput(s, fmt.Sprint("-> GetBlockState maxHeight = ", s.maxHeight, " / Request to ", s.selectNode), glogger.BlockConsensus)
	} else {
		s.blockMachine.BlockServer.BlockServerInitStart()
		s.blockMachine.setState(s.blockMachine.miningState)
		s.glog.DebugOutput(s, fmt.Sprint("-> MiningState maxHeight = ", s.maxHeight), glogger.BlockConsensus)
	}
	return false
}

func (s *GetHeigtestState) Exit() {

}
