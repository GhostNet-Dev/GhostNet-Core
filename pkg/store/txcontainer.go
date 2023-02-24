package store

import (
	"github.com/GhostNet-Dev/GhostNet-Core/libs/container"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gsql"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type TxContainer struct {
	gSql             gsql.GSql
	gCandidateSql    gsql.GCandidateSql
	CandidateTxPools *container.Queue
	CurrentPoolId    uint32
}

func NewTxContainer(g gsql.GSql, gCandidateSql gsql.GCandidateSql) *TxContainer {
	return &TxContainer{
		gSql:             g,
		gCandidateSql:    gCandidateSql,
		CandidateTxPools: container.NewQueue(),
	}
}

func (txContainer *TxContainer) Initialize() {
	txContainer.CurrentPoolId = txContainer.gCandidateSql.GetMaxPoolId()
}

func (txContainer *TxContainer) SaveTransaction(blockId uint32, tx *types.GhostTransaction, txIndexInBlock uint32) {
	txContainer.gSql.InsertTx(blockId, tx, types.NormalTx, txIndexInBlock)
}

func (txContainer *TxContainer) GetUnusedOutputList(txType types.TxOutputType, toAddr []byte) []types.PrevOutputParam {
	return txContainer.gSql.SelectUnusedOutputs(txType, toAddr)
}

func (txContainer *TxContainer) CheckExistTxId(txId []byte) bool {
	return txContainer.gSql.CheckExistTxId(txId)
}

func (txContainer *TxContainer) GetTx(txId []byte) *types.GhostTransaction {
	return txContainer.gSql.SelectTx(txId)
}

func (txContainer *TxContainer) CheckRefExist(refTxId []byte, outIndex uint32, notTxId []byte) bool {
	return txContainer.gSql.CheckExistRefOutout(refTxId, outIndex, notTxId)
}

func (txContainer *TxContainer) GetMaxLogicalAddress(owner []byte) (uint64, error) {
	return txContainer.gSql.GetMaxLogicalAddress(owner)
}

func (txContainer *TxContainer) GetNicknameToAddress(nickname string) []byte {
	return txContainer.gSql.GetNicknameToAddress([]byte(nickname))
}
