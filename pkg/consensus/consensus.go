package consensus

import (
	"bytes"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type Consensus struct {
	blockContainer *store.BlockContainer
}

func NewConsensus(bc *store.BlockContainer) *Consensus {
	InitInterval()
	return &Consensus{
		blockContainer: bc,
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
