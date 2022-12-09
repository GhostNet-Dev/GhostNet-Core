package blocks

import (
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

func (blocks *Blocks) MakeNewBlock(ghostAddrss *gcrypto.GhostAddress) *types.PairedBlock {
	height := blocks.blockContainer.BlockHeight()
	if height < 2 {
		return nil
	}

	pairedBlock := blocks.blockContainer.GetBlock(height)
	if pairedBlock == nil {
		return nil
	}

	newId := pairedBlock.Block.Header.Id + 1
	prevHash := pairedBlock.Block.GetHashKey()
	prevDataHash := pairedBlock.DataBlock.GetHashKey()
	newUsedTxPool := blocks.blockContainer.MakeCandidateTrPool(newId, blocks.GetMinimumReqTrCount())
	if newUsedTxPool == nil {
		return nil
	}

	dataBlock := blocks.CreateGhostNetDataBlock(newId, prevDataHash, newUsedTxPool.DataTxCandidate)
	block := blocks.CreateGhostNetBlock(newId, prevHash, dataBlock.GetHashKey(), ghostAddrss,
		newUsedTxPool.TxCandidate)

	return &types.PairedBlock{
		Block:     *block,
		DataBlock: *dataBlock,
	}
}

func (blocks *Blocks) CreateGhostNetBlock(newBlockId uint32, prevBlockHash []byte, dataBlockhash []byte,
	address *gcrypto.GhostAddress, newTxList []types.GhostTransaction) *types.GhostNetBlock {
	return &types.GhostNetBlock{
		Header: types.GhostNetBlockHeader{
			Id:                      newBlockId,
			Version:                 blocks.Version,
			PreviousBlockHeaderHash: prevBlockHash,
			MerkleRoot:              blocks.CreateMerkleRoot(newTxList),
			DataBlockHeaderHash:     dataBlockhash,
			TransactionCount:        uint32(len(newTxList)),
			TimeStamp:               blocks.DateTimeToUnixTimeNow(),
		},
	}
}

func (blocks *Blocks) CreateGhostNetDataBlock(newBlockId uint32, prevBlockHash []byte,
	newTxList []types.GhostDataTransaction) *types.GhostNetDataBlock {
	return &types.GhostNetDataBlock{
		Header: types.GhostNetDataBlockHeader{
			Id:                      newBlockId,
			Version:                 blocks.Version,
			PreviousBlockHeaderHash: prevBlockHash,
			MerkleRoot:              blocks.CreateMerkleDataRoot(newTxList),
			TransactionCount:        uint32(len(newTxList)),
		},
	}
}

func (blocks *Blocks) DateTimeToUnixTimeNow() uint64 {
	return uint64(time.Now().Unix())
}
