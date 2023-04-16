package blocks

import (
	"log"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

const (
	CoinBase uint64 = 1_000_000
)

func (blocks *Blocks) MinerStart() {
	blocks.miningFlag = true
}

func (blocks *Blocks) MinerStop() {
	blocks.miningFlag = false
}

func (blocks *Blocks) MakeNewBlock(miner *gcrypto.GhostAddress, creator []byte,
	minimumRequiredTxCount uint32) *types.PairedBlock {
	height := blocks.blockContainer.BlockHeight()
	if height < 1 {
		return nil
	}

	pairedBlock := blocks.blockContainer.GetBlock(height)
	if pairedBlock == nil {
		return nil
	}

	newId := pairedBlock.Block.Header.Id + 1
	prevHash := pairedBlock.Block.GetHashKey()
	prevDataHash := pairedBlock.DataBlock.GetHashKey()
	newUsedTxPool := blocks.blockContainer.TxContainer.MakeCandidateTrPool(newId, minimumRequiredTxCount)
	if newUsedTxPool == nil {
		return nil
	}

	dataBlock := blocks.CreateGhostNetDataBlock(newId, prevDataHash, newUsedTxPool.DataTxCandidate)
	block := blocks.CreateGhostNetBlock(newId, prevHash, dataBlock.GetHashKey(), miner, creator,
		newUsedTxPool.TxCandidate)

	newPairedBlock := &types.PairedBlock{
		Block:     *block,
		DataBlock: *dataBlock,
	}

	if blocks.BlockValidation(newPairedBlock, pairedBlock) {
		blocks.blockContainer.InsertBlock(newPairedBlock)
		blocks.blockContainer.TxContainer.DeleteCandidatePool(newUsedTxPool.PoolId)
		return newPairedBlock
	}

	return nil
}

func (blocks *Blocks) CreateGhostNetBlock(newBlockId uint32, prevBlockHash []byte, dataBlockhash []byte,
	miner *gcrypto.GhostAddress, broker []byte, newTxList []types.GhostTransaction) *types.GhostNetBlock {
	hashs := make([][]byte, len(newTxList))
	for i, tx := range newTxList {
		hashs[i] = tx.GetHashKey()
	}
	block := &types.GhostNetBlock{
		Header: types.GhostNetBlockHeader{
			Id:                      newBlockId,
			Version:                 blocks.Version,
			PreviousBlockHeaderHash: prevBlockHash,
			MerkleRoot:              CreateMerkleRoot(hashs),
			DataBlockHeaderHash:     dataBlockhash,
			AliceCount:              1,
			TransactionCount:        uint32(len(newTxList)),
			TimeStamp:               blocks.DateTimeToUnixTimeNow(),
		},
		Alice:       []types.GhostTransaction{*blocks.MakeAliceCoin(newBlockId, broker, newTxList)},
		Transaction: newTxList,
	}

	return blocks.InkTheBlock(block, miner)
}

func (blocks *Blocks) InkTheBlock(block *types.GhostNetBlock, ghostAddr *gcrypto.GhostAddress) *types.GhostNetBlock {
	blocks.gScript.MakeBlockSignature(block, ghostAddr)
	return block
}

func (blocks *Blocks) CreateGhostNetDataBlock(newBlockId uint32, prevBlockHash []byte,
	newTxList []types.GhostDataTransaction) *types.GhostNetDataBlock {
	hashs := make([][]byte, len(newTxList))
	for i, tx := range newTxList {
		hashs[i] = tx.GetHashKey()
	}
	return &types.GhostNetDataBlock{
		Header: types.GhostNetDataBlockHeader{
			Id:                      newBlockId,
			Version:                 blocks.Version,
			PreviousBlockHeaderHash: prevBlockHash,
			MerkleRoot:              CreateMerkleRoot(hashs),
			TransactionCount:        uint32(len(newTxList)),
		},
		Transaction: newTxList,
	}
}

func (blocks *Blocks) DateTimeToUnixTimeNow() uint64 {
	return uint64(time.Now().Unix())
}

func (blocks *Blocks) MakeAliceCoin(blockId uint32, adamsAddr []byte,
	txs []types.GhostTransaction) *types.GhostTransaction {
	brokerGather := map[string]uint64{}
	if len(txs) < 1 {
		log.Fatal("not enough tx")
	}
	txCoin := CoinBase / uint64(len(txs))
	var totalRealSum uint64
	for _, tx := range txs {
		outCoin := txCoin / uint64(tx.Body.OutputCounter)

		for _, output := range tx.Body.Vout {
			broker := string(output.BrokerAddr)
			brokerGather[broker] += outCoin
			totalRealSum += outCoin
		}
	}

	if totalRealSum != CoinBase {
		remain := CoinBase - totalRealSum
		broker := string(adamsAddr)
		brokerGather[broker] += remain
	}

	outputs := []types.TxOutput{}
	for key, coin := range brokerGather {
		// TODO: string -> byte -> string이 같은지 확인이 필요하다.
		pubKey := []byte(key)
		baseScript := gvm.MakeLockScriptOut(pubKey)
		outputs = append(outputs, types.TxOutput{
			Addr:         pubKey,
			BrokerAddr:   pubKey,
			Value:        coin,
			ScriptPubKey: baseScript,
			ScriptSize:   uint32(len(baseScript)),
		})
	}

	tx := &types.GhostTransaction{
		TxId: make([]byte, gbytes.HashSize),
		Body: types.TxBody{
			InputCounter:  0,
			Vout:          outputs,
			OutputCounter: uint32(len(outputs)),
			Nonce:         blockId,
		},
	}
	tx.TxId = tx.GetHashKey()
	return tx
}
