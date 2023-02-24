package txs

import (
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestSaveDataTx(t *testing.T) {
	rootTx := MakeRootFsTx()
	blockContainer.TxContainer.SaveTransaction(0, rootTx, 0)
	prev := map[types.TxOutputType][]types.PrevOutputParam{
		types.TxTypeDataTransfer: {
			{
				TxType: types.TxTypeFSRoot,
				VOutPoint: types.TxOutPoint{
					TxId:       rootTx.TxId,
					TxOutIndex: 0,
				},
				Vout: rootTx.Body.Vout[0],
			},
		},
	}
	txInfo := TransferTxInfo{
		MyWallet:     *MyWallet,
		ToAddr:       Recver.Get160PubKey(),
		Broker:       Broker.Get160PubKey(),
		FeeAddr:      Broker.Get160PubKey(),
		FeeBroker:    Broker.Get160PubKey(),
		Prevs:        prev,
		TransferCoin: 9999,
	}
	data := []byte("hello blockchain")
	tx, dataTx := txs.CreateDataTx(txInfo, 0, data)
	tx = txs.InkTheContract(tx, Recver)
	err := txs.TransactionValidation(tx, dataTx, blockContainer.TxContainer)
	assert.Equal(t, true, err.Result(), "tx validate error: "+err.Error())
}
