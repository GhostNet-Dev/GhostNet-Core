package consensus

import "github.com/GhostNet-Dev/GhostNet-Core/pkg/store"

const (
	INIT_MAX_TRANSACTION_COUNT = 1
	CREATE_BLOCK_INTERVAL      = 60
)

var (
	maxTxCountDic = make(map[uint32]uint32)
)

func InitInterval() {
	maxTxCountDic[0] = INIT_MAX_TRANSACTION_COUNT
}

func (con *Consensus) GetMinimumReqTrCount() uint32 {
	height := con.blockContainer.BlockHeight()
	if height < 2 {
		return 0
	}

	maxTransactionCnt := uint32(INIT_MAX_TRANSACTION_COUNT)
	if height > CREATE_BLOCK_INTERVAL {
		maxTransactionCnt = con.GetMaxTransactionCount(height)
	}

	return maxTransactionCnt
}

func coreCalculator(currTime uint64, prevTime uint64, prevTxCount uint32) uint32 {
	span := float64(currTime - prevTime)
	result := uint32((1 / (span / CREATE_BLOCK_INTERVAL)) * float64(prevTxCount))
	if result >= INIT_MAX_TRANSACTION_COUNT {
		return result
	}
	return INIT_MAX_TRANSACTION_COUNT
}

func calculateMaxTransactionCount(blockContainer *store.BlockContainer,
	height uint32) uint32 {
	slot := height / CREATE_BLOCK_INTERVAL
	startBlock := slot * CREATE_BLOCK_INTERVAL
	endBlock := startBlock + CREATE_BLOCK_INTERVAL - 1

	if startBlock == 0 {
		startBlock = 2
	}

	prevPairedBlock := blockContainer.GetBlock(startBlock)
	pairedBlock := blockContainer.GetBlock(endBlock)
	if prevPairedBlock == nil || pairedBlock == nil {
		return INIT_MAX_TRANSACTION_COUNT
	}

	currTimestamp := pairedBlock.Block.Header.TimeStamp
	prevTimestamp := prevPairedBlock.Block.Header.TimeStamp
	return coreCalculator(currTimestamp, prevTimestamp, pairedBlock.Block.Header.TransactionCount)
}

func (con *Consensus) GetMaxTransactionCount(height uint32) uint32 {
	maxTransactionCnt := uint32(INIT_MAX_TRANSACTION_COUNT)
	slot := height / CREATE_BLOCK_INTERVAL

	if _, ok := maxTxCountDic[slot]; ok {
		maxTransactionCnt = calculateMaxTransactionCount(con.blockContainer, height)
		maxTxCountDic[slot] = maxTransactionCnt
	} else {
		maxTransactionCnt = maxTxCountDic[slot]
	}
	return maxTransactionCnt
}
