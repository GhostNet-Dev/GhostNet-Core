package store

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gsql"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type TxContainer struct {
	TableName string
	gSql      gsql.GSql
}

func NewTxContainer(g gsql.GSql, tableName string) *TxContainer {
	return &TxContainer{
		TableName: tableName,
		gSql:      g,
	}
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
