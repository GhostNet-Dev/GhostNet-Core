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
	s.maxHeight = s.blockMachine.blockContainer.BlockHeight()
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
func (s *GetHeigtestState) GetValidNode(candidatePool map[uint32][]string) {

}

func (s *GetHeigtestState) TimerExpired(context interface{}) bool {
	<-time.After(time.Second * 8)

	curBlockHeight := s.blockMachine.blockContainer.BlockHeight()
	pubKey, candiList, maxHeight := s.blockMachine.BlockServer.CheckValidNode(s.candidatePool, s.maxHeight)
	if maxHeight > curBlockHeight && maxHeight != 0 {
		s.blockMachine.SetHeighestCandidatePool(candiList)
		s.blockMachine.SetTargetHeight(maxHeight)
		s.blockMachine.setState(s.blockMachine.verificationState)
		s.blockMachine.BlockServer.RequestGetBlockHash(pubKey, curBlockHeight)
		s.glog.DebugOutput(s, fmt.Sprint("-> GetBlockState maxHeight = ", s.maxHeight, "(", maxHeight, ") / Request to ", s.selectNode), glogger.BlockConsensus)
	} else {
		s.blockMachine.BlockServer.BlockServerInitStart()
		s.blockMachine.setState(s.blockMachine.miningState)
		s.glog.DebugOutput(s, fmt.Sprint("-> MiningState maxHeight = ", s.maxHeight), glogger.BlockConsensus)
	}
	return false
}

func (s *GetHeigtestState) Exit() {

}
