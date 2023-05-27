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

	dbFilePath           string
	newBlockEventHandler []func(*types.PairedBlock)
	delBlockEventHandler []func(*types.PairedBlock)
}

func NewBlockContainer(dbType string) *BlockContainer {
	g := gsql.NewGSql(dbType)
	gm := gsql.NewMergeGSql(dbType)
	gCandidate := g.(gsql.GCandidateSql) // gsql.NewGCandidateSql("sqlite3")
	bc := &BlockContainer{
		gSql:                 g,
		gMergeSql:            gm,
		TxContainer:          NewTxContainer(g, gCandidate),
		MergeTxContainer:     NewTxContainer(gm, nil),
		newBlockEventHandler: make([]func(*types.PairedBlock), 0),
		delBlockEventHandler: make([]func(*types.PairedBlock), 0),
	}
	bc.CandidateBlk = NewCandidateBlock(g, gm, bc)

	return bc
}

func (blockContainer *BlockContainer) RegisterBlockEvent(newBlockEvent func(*types.PairedBlock),
	delBlockEvent func(*types.PairedBlock)) {
	blockContainer.newBlockEventHandler = append(blockContainer.newBlockEventHandler, newBlockEvent)
	blockContainer.delBlockEventHandler = append(blockContainer.delBlockEventHandler, delBlockEvent)
}

func (blockContainer *BlockContainer) BlockContainerOpen(dbFilePath string) {
	sqlOpen(blockContainer.gSql, dbFilePath, DbFilename)
	sqlOpen(blockContainer.gMergeSql, dbFilePath, MergeDbFilename)

	blockContainer.dbFilePath = dbFilePath
	blockContainer.TxContainer.Initialize()
}

func sqlOpen(gSql gsql.GSql, dbFilePath string, dbFilename string) {
	if err := gSql.OpenSQL(dbFilePath, dbFilename); err != nil {
		log.Fatal(err)
	}
	if err := gSql.CreateTable(); err != nil {
		log.Fatal(err)
	}
}

func (blockContainer *BlockContainer) Reopen() {
	blockContainer.BlockContainerOpen(blockContainer.dbFilePath)
}

func (blockContainer *BlockContainer) Close() {
	blockContainer.gSql.CloseSQL()
}

func (blockContainer *BlockContainer) Reset() {
	blockContainer.gSql.DropTable()
	blockContainer.CandidateBlk.DropTable()
	blockContainer.BlockContainerOpen(blockContainer.dbFilePath)
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

func (blockContainer *BlockContainer) GetIssuedCoinOnBlock(blockId uint32) uint64 {
	return blockContainer.gSql.GetIssuedCoin(blockId)
}

func (blockContainer *BlockContainer) InsertBlock(pairedBlock *types.PairedBlock) {
	for _, handler := range blockContainer.newBlockEventHandler {
		handler(pairedBlock)
	}
	blockContainer.gSql.InsertBlock(pairedBlock)
}

func (blockContainer *BlockContainer) DeleteAfterTargetId(blockId uint32) {
	if len(blockContainer.delBlockEventHandler) > 0 {
		height := blockContainer.BlockHeight()
		for i := blockId; i < height; i++ {
			pairedBlock := blockContainer.GetBlock(blockId)
			for _, handler := range blockContainer.delBlockEventHandler {
				handler(pairedBlock)
			}
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
