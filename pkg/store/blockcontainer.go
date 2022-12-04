package store

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gsql"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type BlockContainer struct {
	gSql gsql.GSql
}

func NewBlockContainer() *BlockContainer {
	return &BlockContainer{
		gSql: gsql.NewGSql("sqlite"),
	}
}

func (blockContainer *BlockContainer) BlockContainerOpen(schemeSqlFilePath string, dbFilePath string) {
	blockContainer.gSql.OpenSQL(dbFilePath)
	blockContainer.gSql.CreateTable(schemeSqlFilePath)
}

func (blockContainer *BlockContainer) SaveTransaction(blockId uint32, tx types.GhostTransaction, txIndexInBlock uint32) {
	blockContainer.gSql.InsertTx(blockId, tx, types.NormalTx, txIndexInBlock)
}

func (blockContainer *BlockContainer) GetUnusedOutputList(txType uint32, toAddr []byte) []types.PrevOutputParam {
	return blockContainer.gSql.SelectUnusedOutputs(txType, toAddr)
}

func (blockContainer *BlockContainer) CheckExistTxId(txId []byte) bool {
	return blockContainer.gSql.CheckExistTxId(txId)
}

func (blockContainer *BlockContainer) GetTx(txId []byte) *types.GhostTransaction {
	return blockContainer.gSql.SelectTx(txId)
}
