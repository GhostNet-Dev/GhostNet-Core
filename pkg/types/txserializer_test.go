package types

import (
	"bytes"
	"crypto/sha256"
	"testing"

	"github.com/stretchr/testify/assert"
	mems "github.com/traherom/memstream"
)

func TestTxSerilalize(t *testing.T) {
	tx := GhostTransaction{
		TxId: make([]byte, 32),
	}
	size := tx.Size()
	stream := mems.NewCapacity(int(size))
	tx.Serialize(stream)
	assert.Equal(t, int(tx.Size()), len(stream.Bytes()), "Size가 다릅니다.")
}

func TestTxOutputSerializeDeserialize(t *testing.T) {
	output := MakeTxOutput()
	size := output.Size()
	stream := mems.NewCapacity(int(size))
	output.Serialize(stream)
	byteBuf := bytes.NewBuffer(stream.Bytes())

	newOutput := TxOutput{}
	newOutput.Deserialize(byteBuf)
	result := bytes.Compare(output.Addr, newOutput.Addr)
	assert.Equal(t, 0, result, "binary가 다릅니다.")
	assert.Equal(t, int(size), len(stream.Bytes()), "Size가 다릅니다.")
	assert.Equal(t, output.Value, newOutput.Value, "Value가 다릅니다.")
}

func TestTxInputSerializeDeserialize(t *testing.T) {
	input := MakeTxInput()
	size := input.Size()
	stream := mems.NewCapacity(int(size))
	input.Serialize(stream)
	byteBuf := bytes.NewBuffer(stream.Bytes())

	newInput := TxInput{}
	newInput.Deserialize(byteBuf)
	assert.Equal(t, int(size), len(stream.Bytes()), "Size가 다릅니다.")
	assert.Equal(t, input.Sequence, newInput.Sequence, "Value가 다릅니다.")
}

func TestTxBodySerializeDeserialize(t *testing.T) {
	body := MakeTxBody()
	size := body.Size()
	stream := mems.NewCapacity(int(size))
	body.Serialize(stream)
	byteBuf := bytes.NewBuffer(stream.Bytes())

	newBody := TxBody{}
	newBody.Deserialize(byteBuf)
	result := bytes.Compare(body.Vout[0].Addr, newBody.Vout[0].Addr)
	assert.Equal(t, 0, result, "binary가 다릅니다.")
	assert.Equal(t, int(size), len(stream.Bytes()), "Size가 다릅니다.")
	assert.Equal(t, body.Nonce, newBody.Nonce, "Value가 다릅니다.")
}

func MakeTxOutput() TxOutput {
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

func MakeTxInput() TxInput {
	dummy := make([]byte, 4)
	hash := sha256.New()
	hash.Write(dummy)
	key := hash.Sum((nil))

	input := TxInput{
		PrevOut: TxOutPoint{
			TxId:       key,
			TxOutIndex: 0,
		},
		Sequence:   3232,
		ScriptSize: 4,
		ScriptSig:  dummy,
	}
	return input
}

func MakeTxBody() TxBody {
	return TxBody{
		InputCounter: 2,
		Vin: []TxInput{
			MakeTxInput(),
			MakeTxInput(),
		},
		OutputCounter: 1,
		Vout: []TxOutput{
			MakeTxOutput(),
		},
		Nonce:    2233,
		LockTime: 1234,
	}
}

func MakeTx() GhostTransaction {
	txBody := MakeTxBody()
	stream := mems.NewCapacity(int(txBody.Size()))
	txBody.Serialize(stream)
	hash := sha256.New()
	hash.Write(stream.Bytes())
	txId := hash.Sum((nil))
	return GhostTransaction{txId, txBody}
}
