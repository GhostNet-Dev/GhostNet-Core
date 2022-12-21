package store

import "github.com/GhostNet-Dev/GhostNet-Core/pkg/types"

type CandidateTxPool struct {
	BlockId         uint32
	PoolId          uint32
	TxCandidate     []types.GhostTransaction
	DataTxCandidate []types.GhostDataTransaction
}

func (blockContainer *BlockContainer) GetCandidateTxPool(blockId uint32) *CandidateTxPool {
	return blockContainer.CandidateTxPools.Pop().(*CandidateTxPool)
}

func (blockContainer *BlockContainer) MakeCandidateTrPool(blockId uint32, minTxCount uint32) *CandidateTxPool {
	minPoolId := blockContainer.gSql.GetMinPoolId()
	if blockContainer.CurrentPoolId-minPoolId > 5 {
		blockContainer.gSql.UpdatePoolId(minPoolId, blockContainer.CurrentPoolId)
	}
	txList := blockContainer.gSql.SelectTxsPool(blockContainer.CurrentPoolId)
	dataTxList := blockContainer.gSql.SelectDataTxsPool(blockContainer.CurrentPoolId)
	// validation은 여기서 하지 않는다. 밖에서 하고 들어온다.
	if minTxCount < 2 || len(txList) == 0 || len(txList)+len(dataTxList) < int(minTxCount) {
		return nil
	}
	blockContainer.CurrentPoolId++
	return &CandidateTxPool{
		BlockId:         blockId,
		TxCandidate:     txList,
		DataTxCandidate: dataTxList,
	}
}
