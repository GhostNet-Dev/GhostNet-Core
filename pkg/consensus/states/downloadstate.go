package states

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type DownloadCheckState struct {
	blockMachine *consensus.BlockMachine
}

func (s *DownloadCheckState) Inititalize() {
}

func (s *DownloadCheckState) Rebuild() {

}

func (s *DownloadCheckState) StartMining() {

}

func (s *DownloadCheckState) RecvBlockHeight(height uint32, pubKey string) {

}

func (s *DownloadCheckState) RecvBlockHash(from string, masterHash string, blockIdx uint32) {

}
func (s *DownloadCheckState) RecvBlock(pairedBlock *types.PairedBlock, pubKey string) {

}

func (s *DownloadCheckState) Exit() {

}
