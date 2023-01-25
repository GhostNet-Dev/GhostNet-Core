package blockmanager

func (blockMgr *BlockManager) BroadcastBlockChainNotification() {

}

func (blockMgr *BlockManager) MiningStart() {

}

func (blockMgr *BlockManager) MiningStop() {

}

func (blockMgr *BlockManager) SetHeighestCandidatePool(candidateList []string) {

}

func (blockMgr *BlockManager) RequestGetBlock(pubKey string, blockIdx uint32) {

}

func (blockMgr *BlockManager) RequestGetBlockHash(pubKey string, currBlockHeight uint32) {

}

func (blockMgr *BlockManager) MergeErrorNotification(pubKey string, result bool) {

}

func (blockMgr *BlockManager) LocalBlockDbValidation() bool {
	return false
}

func (blockMgr *BlockManager) BlockServerInitStart() {
	blockMgr.consensus.Clear()
	//TODO it needs to more clear!

}

func (blockMgr *BlockManager) CheckHeightForRebuild(neighborHeight uint32) bool {
	currHeight := blockMgr.blockContainer.BlockHeight()

	if currHeight < neighborHeight {
		return true
	}
	return false
}
