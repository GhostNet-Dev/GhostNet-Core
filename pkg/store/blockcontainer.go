package store

import (
	gsql "github.com/GhostNet-Dev/GhostNet-Core/pkg/gsql"
	types "github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
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
