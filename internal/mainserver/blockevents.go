package mainserver

import "github.com/GhostNet-Dev/GhostNet-Core/pkg/types"

func (blockMgr *BlockManager) BroadcastBlockChainNotification() {

}

func (blockMgr *BlockManager) MiningStart() {

}

func (blockMgr *BlockManager) MiningStop() {

}

func (blockMgr *BlockManager) CheckBlockHeight() {

}

func (blockMgr *BlockManager) SetHeighestCandidatePool() {

}

func (blockMgr *BlockManager) LoadHashFromTempDb(blockIdx uint32) string {
	return ""
}

func (blockMgr *BlockManager) RequestGetBlock(blockIdx uint32) {

}

func (blockMgr *BlockManager) RequestGetBlockHash(currBlockHeight uint32) {

}

func (blockMgr *BlockManager) CheckAndSave(pairedBlock *types.PairedBlock) bool {
	return false
}

func (blockMgr *BlockManager) MergeErrorNotification() {

}

func (blockMgr *BlockManager) LocalBlockDbValidation() bool {
	return false
}
