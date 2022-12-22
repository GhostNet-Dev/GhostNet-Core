package mainserver

import "github.com/GhostNet-Dev/GhostNet-Core/pkg/types"

func (server *BlockServer) BroadcastBlockChainNotification() {

}

func (server *BlockServer) MiningStart() {

}

func (server *BlockServer) MiningStop() {

}

func (server *BlockServer) CheckBlockHeight() {

}

func (server *BlockServer) SetHeighestCandidatePool() {

}

func (server *BlockServer) LoadHashFromTempDb(blockIdx uint32) string {
	return ""
}

func (server *BlockServer) RequestGetBlock(blockIdx uint32) {

}

func (server *BlockServer) RequestGetBlockHash(currBlockHeight uint32) {

}

func (server *BlockServer) CheckAndSave(pairedBlock *types.PairedBlock) bool {
	return false
}

func (server *BlockServer) MergeErrorNotification() {

}

func (server *BlockServer) LocalBlockDbValidation() bool {
	return false
}
