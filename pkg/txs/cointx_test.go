package txs

import (
	"log"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"github.com/stretchr/testify/assert"
)

var (
	Sender   = gcrypto.GenerateKeyPair()
	Broker   = gcrypto.GenerateKeyPair()
	Recver   = gcrypto.GenerateKeyPair()
	MyWallet = gcrypto.NewWallet("", Sender, nil, nil)

	gScript        = gvm.NewGCompiler()
	gVm            = gvm.NewGVM()
	blockContainer = store.NewBlockContainer("sqlite3")
	txs            = NewTXs(gScript, blockContainer, gVm)
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	blockContainer.BlockContainerOpen("./")
	txInfo := TransferTxInfo{
		FromAddr:     MyWallet.MyPubKey(),
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
	txInfo := TransferTxInfo{
		FromAddr:     MyWallet.MyPubKey(),
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

	prevMap := map[types.TxOutputType][]types.PrevOutputParam{}
	prevMap[types.TxTypeCoinTransfer] = outputParams

	txInfo := TransferTxInfo{
		FromAddr:     MyWallet.MyPubKey(),
		ToAddr:       Recver.Get160PubKey(),
		Broker:       Broker.Get160PubKey(),
		FeeAddr:      Broker.Get160PubKey(),
		FeeBroker:    Broker.Get160PubKey(),
		Prevs:        prevMap,
		TransferCoin: transferCoin,
	}
	tx := txs.TransferCoin(txInfo)
	tx = txs.InkTheContract(tx, Recver)

	err := txs.TransactionValidation(tx, nil, blockContainer.TxContainer, 0)
	assert.Equal(t, true, err.Result(), "tx validate error: "+err.Error())
}

//TODO Execute Tx Script Test
