package store

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gsql"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type CandidateBlock struct {
	gSql           gsql.GSql
	gMergeSql      gsql.GSql
	blockContainer *BlockContainer
	maxBlockId     uint32
}

func NewCandidateBlock(gSql gsql.GSql, gMergeSql gsql.GSql, blockContainer *BlockContainer) *CandidateBlock {
	return &CandidateBlock{
		gSql:           gSql,
		gMergeSql:      gMergeSql,
		blockContainer: blockContainer,
		maxBlockId:     0,
	}
}

func (candidateBlock *CandidateBlock) DropTable() {
	candidateBlock.gMergeSql.DropTable()
	candidateBlock.maxBlockId = 0
}

func (candidateBlock *CandidateBlock) Reset() {
	candidateBlock.gMergeSql.DropTable()
	sqlOpen(candidateBlock.gMergeSql, candidateBlock.blockContainer.dbFilePath, MergeDbFilename)
	candidateBlock.maxBlockId = 0
}

func (candidateBlock *CandidateBlock) GetBlock(blockId uint32) (paired *types.PairedBlock) {
	if candiPaired := candidateBlock.gMergeSql.SelectBlock(blockId); candiPaired == nil {
		paired = candidateBlock.gSql.SelectBlock(blockId)
	} else {
		paired = candiPaired
	}
	return paired
}

func (candidateBlock *CandidateBlock) DeleteBlock(blockId uint32) {
	candidateBlock.gMergeSql.DeleteBlock(blockId)
}

func (candidateBlock *CandidateBlock) AddBlock(pairedBlock *types.PairedBlock) {
	if candidateBlock.maxBlockId < pairedBlock.BlockId() {
		candidateBlock.maxBlockId = pairedBlock.BlockId()
	}
	candidateBlock.gMergeSql.InsertBlock(pairedBlock)
}

func (candidateBlock *CandidateBlock) GetMaxBlockId() uint32 {
	return candidateBlock.maxBlockId
}
