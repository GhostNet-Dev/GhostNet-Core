package store

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gsql"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type BlockContainer struct {
	gSql        gsql.GSql
	TxContainer *TxContainer
}

func NewBlockContainer() *BlockContainer {
	g := gsql.NewGSql("sqlite3")
	gCandidate := gsql.NewGCandidateSql("sqlite3")
	return &BlockContainer{
		gSql:        g,
		TxContainer: NewTxContainer(g, gCandidate),
	}
}

func (blockContainer *BlockContainer) BlockContainerOpen(schemeSqlFilePath string, dbFilePath string) {
	blockContainer.gSql.OpenSQL(dbFilePath)
	blockContainer.gSql.CreateTable(schemeSqlFilePath)
	blockContainer.TxContainer.Initialize()
}

func (blockContainer *BlockContainer) GetBlock(blockId uint32) *types.PairedBlock {
	return blockContainer.gSql.SelectBlock(blockId)
}

func (blockContainer *BlockContainer) BlockHeight() uint32 {
	return blockContainer.gSql.GetBlockHeight()
}

func (blockContainer *BlockContainer) InsertBlock(pairedBlock *types.PairedBlock) {
	blockContainer.gSql.InsertBlock(pairedBlock)
}
