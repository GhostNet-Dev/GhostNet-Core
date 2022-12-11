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

	CurrentPoolId uint32
}

func NewBlockContainer() *BlockContainer {
	g := gsql.NewGSql("sqlite3")
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

	blockContainer.CurrentPoolId = blockContainer.gSql.GetMaxPoolId()
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
