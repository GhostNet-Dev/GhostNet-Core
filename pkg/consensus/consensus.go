package consensus

import (
	"bytes"
	"fmt"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blocks"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type Consensus struct {
	blockContainer *store.BlockContainer
	block          *blocks.Blocks
	glog           *glogger.GLogger
}

func NewConsensus(bc *store.BlockContainer, block *blocks.Blocks, glog *glogger.GLogger) *Consensus {
	InitInterval()
	return &Consensus{
		blockContainer: bc,
		block:          block,
		glog:           glog,
	}
}

func (con *Consensus) Consensus(pairedBlock *types.PairedBlock) (merge bool, rebuild bool) {
	merge, rebuild = false, false
	block := pairedBlock.Block
	height := con.blockContainer.BlockHeight()
	requiredTxCount := con.GetMaxTransactionCount(height)

	if requiredTxCount < 0 || pairedBlock.TxCount() < requiredTxCount {
		return false, false
	}

	if block.Header.Id == height+1 {
		currentPair := con.blockContainer.GetBlock(height)
		currentBlockHash := currentPair.Block.GetHashKey()
		prevBlockHashInCandidateBlk := block.Header.PreviousBlockHeaderHash
		if bytes.Compare(currentBlockHash, prevBlockHashInCandidateBlk) != 0 {
			rebuild = true
		} else {
			con.blockContainer.InsertBlock(pairedBlock)
			merge = true
		}
	} else if block.Header.Id > height+1 {
		rebuild = true
	}

	return merge, rebuild
}

func (con *Consensus) CheckMinimumTxCount(pairedBlock *types.PairedBlock) bool {
	block := pairedBlock.Block
	height := con.blockContainer.BlockHeight()
	maxTxCount := con.GetMaxTransactionCount(height)

	if maxTxCount < 0 || pairedBlock.TxCount() < maxTxCount {
		return false
	}

	if pairedBlock.BlockId() == height+1 {
		currentLocalPairedBlock := con.blockContainer.GetBlock(height)
		currentBlockHashKey := currentLocalPairedBlock.Block.GetHashKey()
		preBlockhashKeyFromNewBlock := block.Header.PreviousBlockHeaderHash
		if bytes.Compare(currentBlockHashKey, preBlockhashKeyFromNewBlock) == 0 {
			return true
		} else {
			// TODO: trigger rebuild
			return false
		}

	} else if pairedBlock.BlockId() > height+1 {
		// TODO: trigger rebuild
	}
	return false
}

func (con *Consensus) CheckTriggerNewBlock() (bool, uint32) {
	height := con.blockContainer.BlockHeight()
	triggerTxCount := con.GetMaxTransactionCount(height)
	txCount := con.blockContainer.TxContainer.GetCandidateTxCount()
	return txCount >= triggerTxCount, triggerTxCount
}

func (con *Consensus) Clear() {
	con.blockContainer.CandidateBlk.Reset()
}

func (con *Consensus) LoadHashFromTempDb(blockId uint32) []byte {
	pairedBlock := con.blockContainer.CandidateBlk.GetBlock(blockId)
	if pairedBlock == nil {
		return nil
	}

	return pairedBlock.Block.GetHashKey()
}

func (con *Consensus) CheckAndSave(startBlockId uint32, candidatePair *types.PairedBlock) bool {
	blockId := candidatePair.BlockId()
	if checkSaveRetryBlock := con.blockContainer.CandidateBlk.GetBlock(blockId); checkSaveRetryBlock != nil {
		if bytes.Compare(checkSaveRetryBlock.Block.GetHashKey(),
			candidatePair.Block.GetHashKey()) == 0 {
			return true
		} else {
			con.blockContainer.CandidateBlk.DeleteBlock(blockId)
		}
	}

	prevPairedBlock := con.blockContainer.CandidateBlk.GetBlock(blockId - 1)
	if prevPairedBlock == nil {
		return false
	}
	if con.block.BlockMergeCheck(startBlockId, candidatePair, prevPairedBlock) == false {
		return false
	}
	con.blockContainer.CandidateBlk.AddBlock(candidatePair)
	return true
}

func (con *Consensus) MergeExecute(startBlockId uint32, endBlockId uint32) {
	con.blockContainer.DeleteAfterTargetId(startBlockId)
	for blockId := startBlockId; blockId <= endBlockId; blockId++ {
		pairedBlock := con.blockContainer.CandidateBlk.GetBlock(blockId)
		if pairedBlock == nil {
			break
		}
		con.blockContainer.InsertBlock(pairedBlock)
		con.glog.DebugOutput(con, fmt.Sprint("mergeBlockId = ", blockId), glogger.BlockConsensus)
	}
}

func (con *Consensus) CheckIntegrityBlockChainList(startBlockId uint32, endBlockId uint32) uint32 {
	if startBlockId < 2 {
		return endBlockId
	}

	for blockId := startBlockId; blockId <= endBlockId; blockId++ {
		pair := con.blockContainer.GetBlock(blockId)
		prevPair := con.blockContainer.GetBlock(blockId - 1)
		if pair == nil || prevPair == nil {
			return blockId
		}
		if con.block.BlockValidation(pair, prevPair) == false {
			return blockId
		}
	}
	return endBlockId
}

func (con *Consensus) LocalBlockChainValidation() bool {
	height := con.blockContainer.BlockHeight()
	startBlockId := con.CheckIntegrityBlockChainList(2, height)
	if startBlockId != height {
		con.glog.DebugOutput(con, fmt.Sprint("Delete After BlockId = ", startBlockId), glogger.BlockConsensus)
		con.blockContainer.DeleteAfterTargetId(startBlockId)
		return false
	}

	return true
}
