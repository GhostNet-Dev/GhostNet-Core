package types

import (
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
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
	buf := tx.Serialize()
	assert.Equal(t, tx.Size(), len(buf), "Size가 다릅니다.")
}
func TestTxOutputSerializeDeserialize(t *testing.T) {
	output := makeTxOutput()
	newOutput := TxOutput{}
	buf := output.Serialize()
	newOutput.DeserializeTxOutput(buf)
	assert.Equal(t, newOutput.Value, output.Value, "Value가 다릅니다.")
}
