package store

import (
	"github.com/GhostNet-Dev/GhostNet-Core/libs/container"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gsql"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type BlockContainer struct {
	gSql                 gsql.GSql
	TxContainer          *TxContainer
	CandidateTxContainer *TxContainer

	CandidateTxPools *container.Queue
}

func NewBlockContainer() *BlockContainer {
	g := gsql.NewGSql("sqlite")
	return &BlockContainer{
		gSql:                 g,
		TxContainer:          NewTxContainer(g, "transactions"),
		CandidateTxContainer: NewTxContainer(g, "c_transactions"),
		CandidateTxPools:     container.NewQueue(),
	}
}

func (blockContainer *BlockContainer) BlockContainerOpen(schemeSqlFilePath string, dbFilePath string) {
	blockContainer.gSql.OpenSQL(dbFilePath)
	blockContainer.gSql.CreateTable(schemeSqlFilePath)
}

func (blockContainer *BlockContainer) GetBlock(blockId uint32) *types.PairedBlock {
	return blockContainer.gSql.SelectBlock(blockId)
}

func (blockContainer *BlockContainer) BlockHeight() uint32 {
	return blockContainer.gSql.GetBlockHeight()
}

type CandidateTxPool struct {
	BlockId         uint32
	PoolId          uint32
	TxCandidate     []types.GhostTransaction
	DataTxCandidate []types.GhostDataTransaction
}

func (blockContainer *BlockContainer) GetCandidateTxPool() *CandidateTxPool {
	return blockContainer.CandidateTxPools.Pop().(*CandidateTxPool)
}
