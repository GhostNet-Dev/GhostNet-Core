package store

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type CandidateBlock struct {
	blockPool      map[uint32]*types.PairedBlock
	blockContainer *BlockContainer
	maxBlockId     uint32
}

func NewCandidateBlock(blockContainer *BlockContainer) *CandidateBlock {
	return &CandidateBlock{
		blockPool:      make(map[uint32]*types.PairedBlock),
		blockContainer: blockContainer,
		maxBlockId:     0,
	}
}

func (candidateBlock *CandidateBlock) Reset() {
	candidateBlock.blockPool = make(map[uint32]*types.PairedBlock)
	candidateBlock.maxBlockId = 0
}

func (candidateBlock *CandidateBlock) GetBlock(blockId uint32) (paired *types.PairedBlock) {
	if candiPaired, exist := candidateBlock.blockPool[blockId]; exist == false {
		paired = candidateBlock.blockContainer.GetBlock(blockId)

	} else {
		paired = candiPaired
	}
	return paired
}

func (candidateBlock *CandidateBlock) DeleteBlock(blockId uint32) {
	delete(candidateBlock.blockPool, blockId)
}

func (candidateBlock *CandidateBlock) AddBlock(pairedBlock *types.PairedBlock) {
	if candidateBlock.maxBlockId < pairedBlock.BlockId() {
		candidateBlock.maxBlockId = pairedBlock.BlockId()
	}
	candidateBlock.blockPool[pairedBlock.BlockId()] = pairedBlock
}

func (candidateBlock *CandidateBlock) GetMaxBlockId() uint32 {
	return candidateBlock.maxBlockId
}
