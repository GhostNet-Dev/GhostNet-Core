package blocks

import (
	"bytes"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

func (blocks *Blocks) BlockMergeCheck(pairedBlock *types.PairedBlock, prevPairedBlock *types.PairedBlock) bool {
	return blocks.BlockValidation(pairedBlock, prevPairedBlock)
}

func (blocks *Blocks) BlockValidation(pairedBlock *types.PairedBlock, prevPairedBlock *types.PairedBlock) bool {
	header := pairedBlock.Block.Header
	txs := pairedBlock.Block.Transaction
	alice := pairedBlock.Block.Alice
	prevBlockId := header.Id - 1

	// 이전 Block과 hash일치 여부
	if prevBlockId > 0 && prevPairedBlock == nil { // for genesis block
		prevPairedBlock = blocks.blockContainer.GetBlock(prevBlockId)
		prevHash := prevPairedBlock.Block.GetHashKey()
		if bytes.Compare(header.PreviousBlockHeaderHash, prevHash) != 0 {
			return false
		}
	}

	// merkle tree check
	hashs := make([][]byte, len(txs))
	for i, tx := range txs {
		hashs[i] = tx.GetHashKey()
	}
	merkleRoot := CreateMerkleRoot(hashs)
	if bytes.Compare(header.MerkleRoot, merkleRoot) != 0 {
		return false
	}

	if header.AliceCount != uint32(len(alice)) || header.TransactionCount != uint32(len(txs)) {
		return false
	}

	if !blocks.AliceTransactionValidation(alice, txs) {
		return false
	}

	for _, tx := range txs {
		txChkResult := blocks.txs.TransactionValidation(&tx, nil, blocks.blockContainer.TxContainer)
		if !txChkResult.Result() {
			return false
		}
	}

	return true
}

func (blocks *Blocks) AliceTransactionValidation(alice []types.GhostTransaction,
	txs []types.GhostTransaction) bool {
	brokerGather := map[string]uint64{}
	txCoin := CoinBase / uint64(len(txs))
	totalRealSum := uint64(0)

	for _, tx := range txs {
		outputCoin := txCoin / uint64(tx.Body.OutputCounter)
		for _, output := range tx.Body.Vout {
			broker := string(output.BrokerAddr)
			brokerGather[broker] += outputCoin
			totalRealSum += outputCoin
		}
	}

	if totalRealSum != CoinBase {
		remainCoin := CoinBase - totalRealSum
		broker := string(store.AdamsAddress())
		brokerGather[broker] += remainCoin
	}

	aliceOutput := alice[0].Body.Vout
	for _, output := range aliceOutput {
		broker := string(output.Addr)
		if _, err := brokerGather[broker]; !err {
			return false
		}
		if brokerGather[broker] != output.Value {
			return false
		}
	}

	return true
}
