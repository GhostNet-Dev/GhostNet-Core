package txs

import (
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/crypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"github.com/stretchr/testify/assert"
)

var (
	Sender = crypto.GenerateKeyPair()
	Broker = crypto.GenerateKeyPair()
	Recver = crypto.GenerateKeyPair()

	gScript        = gvm.NewGScript()
	blockContainer = store.NewBlockContainer()
)

func TestMakeCoinTx(t *testing.T) {
	transferCoin := uint64(10)
	txs := NewTXs(gScript, blockContainer)
	outputParams, ok := txs.CandidateUTXO(transferCoin, Sender.PubKey)

	assert.Equal(t, true, ok, "output이 없습니다. test를 다시 검토하세요")

	var prevMap map[uint32][]types.PrevOutputParam
	prevMap[types.TxTypeCoinTransfer] = outputParams

	txInfo := TransferCoinInfo{
		ToAddr:       Sender.PubKey,
		Broker:       Broker.PubKey,
		FeeAddr:      Broker.PubKey,
		FeeBroker:    Broker.PubKey,
		Prevs:        prevMap,
		TransferCoin: transferCoin,
	}
	tx := txs.TransferCoin(txInfo)
}
