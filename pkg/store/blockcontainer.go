package store

import (
	"bytes"
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gsql"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

const (
	DbFilename      = "block.db"
	MergeDbFilename = "merge.db"
)

type BlockContainer struct {
	gSql             gsql.GSql
	gMergeSql        gsql.GSql
	TxContainer      *TxContainer
	MergeTxContainer *TxContainer
	CandidateBlk     *CandidateBlock

	schemeSqlFilePath    string
	dbFilePath           string
	newBlockEventHandler func(*types.PairedBlock)
	delBlockEventHandler func(*types.PairedBlock)
}

func NewBlockContainer(dbType string) *BlockContainer {
	g := gsql.NewGSql(dbType)
	gm := gsql.NewMergeGSql(dbType)
	gCandidate := g.(gsql.GCandidateSql) // gsql.NewGCandidateSql("sqlite3")
	bc := &BlockContainer{
		gSql:             g,
		gMergeSql:        gm,
		TxContainer:      NewTxContainer(g, gCandidate),
		MergeTxContainer: NewTxContainer(gm, nil),
	}
	bc.CandidateBlk = NewCandidateBlock(g, gm, bc)

	return bc
}

func (blockContainer *BlockContainer) RegisterBlockEvent(newBlockEvent func(*types.PairedBlock),
	delBlockEvent func(*types.PairedBlock)) {
	blockContainer.newBlockEventHandler = newBlockEvent
	blockContainer.delBlockEventHandler = delBlockEvent
}

func (blockContainer *BlockContainer) BlockContainerOpen(schemeSqlFilePath string, dbFilePath string) {
	sqlOpen(blockContainer.gSql, schemeSqlFilePath, dbFilePath, DbFilename)
	sqlOpen(blockContainer.gMergeSql, schemeSqlFilePath, dbFilePath, MergeDbFilename)

	blockContainer.dbFilePath = dbFilePath
	blockContainer.schemeSqlFilePath = schemeSqlFilePath
	blockContainer.TxContainer.Initialize()
}

func sqlOpen(gSql gsql.GSql, schemeSqlFilePath string, dbFilePath string, dbFilename string) {
	if err := gSql.OpenSQL(dbFilePath, dbFilename); err != nil {
		log.Fatal(err)
	}
	if err := gSql.CreateTable(schemeSqlFilePath); err != nil {
		log.Fatal(err)
	}
}

func (blockContainer *BlockContainer) Reopen() {
	blockContainer.BlockContainerOpen(blockContainer.schemeSqlFilePath, blockContainer.dbFilePath)
}

func (blockContainer *BlockContainer) Close() {
	blockContainer.gSql.CloseSQL()
}

func (blockContainer *BlockContainer) Reset() {
	blockContainer.gSql.DropTable()
	blockContainer.CandidateBlk.DropTable()
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
	if blockContainer.newBlockEventHandler != nil {
		blockContainer.newBlockEventHandler(pairedBlock)
	}
	blockContainer.gSql.InsertBlock(pairedBlock)
}

func (blockContainer *BlockContainer) DeleteAfterTargetId(blockId uint32) {
	if blockContainer.delBlockEventHandler != nil {
		height := blockContainer.BlockHeight()
		for i := blockId; i < height; i++ {
			pairedBlock := blockContainer.GetBlock(blockId)
			blockContainer.delBlockEventHandler(pairedBlock)
		}
	}
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
