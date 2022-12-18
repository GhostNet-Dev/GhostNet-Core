package txs

import (
	"bytes"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestMakeRootFsTxSerDesrialize(t *testing.T) {
	tx := txs.CreateRootFsTx(TransferCoinInfo{
		ToAddr:    Recver.Get160PubKey(),
		Broker:    Broker.Get160PubKey(),
		FeeAddr:   Sender.Get160PubKey(),
		FeeBroker: Broker.Get160PubKey(),
	}, "test")
	tx = txs.InkTheContract(tx, Sender)

	buf := tx.SerializeToByte()
	byteBuf := bytes.NewBuffer(buf)

	newTx := types.GhostTransaction{}
	newTx.Deserialize(byteBuf)
	newBuf := newTx.SerializeToByte()
	result := bytes.Compare(buf, newBuf)
	assert.Equal(t, 0, result, "bytes are different.")
}
