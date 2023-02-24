package store

import (
	"bytes"
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gsql"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type BlockContainer struct {
	gSql         gsql.GSql
	TxContainer  *TxContainer
	CandidateBlk *CandidateBlock

	schemeSqlFilePath string
	dbFilePath        string
}

func NewBlockContainer(dbname string) *BlockContainer {
	g := gsql.NewGSql(dbname)
	gCandidate := gsql.NewGCandidateSql("sqlite3")
	bc := &BlockContainer{
		gSql:        g,
		TxContainer: NewTxContainer(g, gCandidate),
	}
	bc.CandidateBlk = NewCandidateBlock(bc)

	return bc
}

func (blockContainer *BlockContainer) BlockContainerOpen(schemeSqlFilePath string, dbFilePath string) {
	if err := blockContainer.gSql.OpenSQL(dbFilePath); err != nil {
		log.Fatal(err)
	}
	if err := blockContainer.gSql.CreateTable(schemeSqlFilePath); err != nil {
		log.Fatal(err)
	}
	blockContainer.dbFilePath = dbFilePath
	blockContainer.schemeSqlFilePath = schemeSqlFilePath
	blockContainer.TxContainer.Initialize()
}

func (blockContainer *BlockContainer) Reopen() {
	blockContainer.BlockContainerOpen(blockContainer.schemeSqlFilePath, blockContainer.dbFilePath)
}

func (blockContainer *BlockContainer) Close() {
	blockContainer.gSql.CloseSQL()
}

func (blockContainer *BlockContainer) Reset() {
	blockContainer.gSql.DropTable()
	blockContainer.CandidateBlk.Reset()
	blockContainer.BlockContainerOpen(blockContainer.schemeSqlFilePath, blockContainer.dbFilePath)
}

func (blockContainer *BlockContainer) GetBlock(blockId uint32) *types.PairedBlock {
	return blockContainer.gSql.SelectBlock(blockId)
}

func (blockContainer *BlockContainer) GetBlockHeader(blockId uint32) (*types.GhostNetBlockHeader, *types.GhostNetDataBlockHeader) {
	return blockContainer.gSql.SelectBlockHeader(blockId)
}

func (blockContainer *BlockContainer) BlockHeight() uint32 {
	return blockContainer.gSql.GetBlockHeight()
}

func (blockContainer *BlockContainer) InsertBlock(pairedBlock *types.PairedBlock) {
	blockContainer.gSql.InsertBlock(pairedBlock)
}

func (blockContainer *BlockContainer) DeleteAfterTargetId(blockId uint32) {
	blockContainer.gSql.DeleteAfterTargetId(blockId)
}

func (blockContainer *BlockContainer) GenesisBlockChecker(genesisBlock *types.PairedBlock) {
	pairedBlock := blockContainer.gSql.SelectBlock(1)
	if pairedBlock == nil {
		blockContainer.gSql.InsertBlock(genesisBlock)
		return
	}
	if bytes.Equal(pairedBlock.Block.GetHashKey(), genesisBlock.Block.GetHashKey()) {
		return
	}
	log.Print("conflict genesis block between db and embed block")
	blockContainer.Reset()
	blockContainer.gSql.InsertBlock(genesisBlock)
}
