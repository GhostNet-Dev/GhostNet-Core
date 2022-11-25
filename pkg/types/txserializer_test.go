package types

import (
	"bytes"
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
	mems "github.com/traherom/memstream"
)

func makeTxOutput() TxOutput {
	dummy := make([]byte, 4)
	hash := sha256.New()
	hash.Write(dummy)
	key := hash.Sum((nil))

	output := TxOutput{
		Addr:         key,
		BrokerAddr:   key,
		Value:        1212,
		ScriptSize:   4,
		ScriptPubKey: dummy,
	}
	return output
}

func TestTxSerilalize(t *testing.T) {
	tx := GhostTransaction{}
	size := tx.Size()
	stream := mems.NewCapacity(int(size))
	tx.Serialize(stream)
	assert.Equal(t, int(tx.Size()), len(stream.Bytes()), "Size가 다릅니다.")
}
func TestTxOutputSerializeDeserialize(t *testing.T) {
	output := makeTxOutput()
	size := output.Size()
	stream := mems.NewCapacity(int(size))
	output.Serialize(stream)
	byteBuf := bytes.NewBuffer(stream.Bytes())

	newOutput := TxOutput{}
	newOutput.DeserializeTxOutput(byteBuf)
	assert.Equal(t, newOutput.Value, output.Value, "Value가 다릅니다.")
}

func TestTxOutputSerializeDeserialize2(t *testing.T) {
	output := makeTxOutput()
	size := output.Size()
	seriBuf := bytes.NewBuffer(make([]byte, size))
	output.Serialize2(seriBuf)

	newOutput := TxOutput{}
	newOutput.DeserializeTxOutput(seriBuf)
	assert.Equal(t, newOutput.Value, output.Value, "Value가 다릅니다.")
}
