package store

import "github.com/GhostNet-Dev/GhostNet-Core/pkg/types"

type CandidateTxPool struct {
	BlockId         uint32
	PoolId          uint32
	TxCandidate     []types.GhostTransaction
	DataTxCandidate []types.GhostDataTransaction
}

func (txContainer *TxContainer) SaveCandidateTx(tx *types.GhostTransaction) {
	txContainer.gCandidateSql.InsertCandidateTx(tx, txContainer.CurrentPoolId)
}

func (txContainer *TxContainer) SaveCandidateDataTx(tx *types.GhostDataTransaction) {
	txContainer.gCandidateSql.InsertCandidateDataTx(tx, txContainer.CurrentPoolId)
}

func (txContainer *TxContainer) DeleteCandidatePool(poolId uint32) {
	txContainer.gCandidateSql.DeleteCandidatePool(poolId)
}

func (txContainer *TxContainer) GetCandidateTxCount() uint32 {
	return txContainer.gCandidateSql.SelectCandidateTxCount()
}

func (txContainer *TxContainer) GetMaxPoolId() uint32 {
	return txContainer.gCandidateSql.GetMaxPoolId()
}

func (txContainer *TxContainer) GetCandidateTxPool() *CandidateTxPool {
	return txContainer.CandidateTxPools.Pop().(*CandidateTxPool)
}

func (txContainer *TxContainer) MakeCandidateTrPool(blockId uint32, minTxCount uint32) *CandidateTxPool {
	minPoolId := txContainer.gCandidateSql.GetMinPoolId()
	if txContainer.CurrentPoolId-minPoolId > 5 {
		txContainer.gCandidateSql.UpdatePoolId(minPoolId, txContainer.CurrentPoolId)
	}
	txList := txContainer.gSql.SelectTxsPool(txContainer.CurrentPoolId)
	dataTxList := txContainer.gSql.SelectDataTxsPool(txContainer.CurrentPoolId)
	// validation은 여기서 하지 않는다. 밖에서 하고 들어온다.
	if len(txList) == 0 || len(txList)+len(dataTxList) < int(minTxCount) {
		return nil
	}
	poolId := txContainer.CurrentPoolId
	txContainer.CurrentPoolId++
	return &CandidateTxPool{
		BlockId:         blockId,
		TxCandidate:     txList,
		PoolId:          poolId,
		DataTxCandidate: dataTxList,
	}
}
