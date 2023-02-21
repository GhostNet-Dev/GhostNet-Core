package txs

import "testing"

func TestSaveDataTx(t *testing.T) {
	txInfo := TransferTxInfo{
		MyWallet:     *MyWallet,
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
