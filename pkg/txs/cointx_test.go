package txs

import (
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"github.com/stretchr/testify/assert"
)

var (
	Sender = gcrypto.GenerateKeyPair()
	Broker = gcrypto.GenerateKeyPair()
	Recver = gcrypto.GenerateKeyPair()

	gScript        = gvm.NewGScript()
	gVm            = gvm.NewGVM()
	blockContainer = store.NewBlockContainer()
	txs            = NewTXs(gScript, blockContainer, gVm)
)

func init() {
	blockContainer.BlockContainerOpen("../../db.sqlite3.sql", "./")
	txInfo := TransferCoinInfo{
		ToAddr:       Recver.Get160PubKey(),
		Broker:       Broker.Get160PubKey(),
		FeeAddr:      Broker.Get160PubKey(),
		FeeBroker:    Broker.Get160PubKey(),
		Prevs:        nil,
		TransferCoin: 9999,
	}
	tx := txs.TransferCoin(txInfo)
	tx = txs.InkTheContract(tx, Recver)
	blockContainer.TxContainer.SaveTransaction(0, tx, 0)
}

func TestSaveCoinTx(t *testing.T) {
	txInfo := TransferCoinInfo{
		ToAddr:       Recver.Get160PubKey(),
		Broker:       Broker.Get160PubKey(),
		FeeAddr:      Broker.Get160PubKey(),
		FeeBroker:    Broker.Get160PubKey(),
		Prevs:        nil,
		TransferCoin: 9999,
	}
	tx := txs.TransferCoin(txInfo)
	tx = txs.InkTheContract(tx, Recver)
	blockContainer.TxContainer.SaveTransaction(0, tx, 0)
}

func TestMakeCoinTx(t *testing.T) {
	transferCoin := uint64(10)
	outputParams, ok := txs.CandidateUTXO(transferCoin, Recver.Get160PubKey())

	assert.Equal(t, true, ok, "output이 없습니다. test를 다시 검토하세요")

	prevMap := map[uint32][]types.PrevOutputParam{}
	prevMap[types.TxTypeCoinTransfer] = outputParams

	txInfo := TransferCoinInfo{
		ToAddr:       Recver.Get160PubKey(),
		Broker:       Broker.Get160PubKey(),
		FeeAddr:      Broker.Get160PubKey(),
		FeeBroker:    Broker.Get160PubKey(),
		Prevs:        prevMap,
		TransferCoin: transferCoin,
	}
	tx := txs.TransferCoin(txInfo)
	tx = txs.InkTheContract(tx, Recver)

	err := txs.TransactionChecker(tx, nil, blockContainer.TxContainer)
	assert.Equal(t, true, err == nil, "tx validate error: "+err.Error())
}
