package types

import (
	"bytes"
	"crypto/sha256"
	"testing"

	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	"github.com/stretchr/testify/assert"
	mems "github.com/traherom/memstream"
)

func TestTxSerilalize(t *testing.T) {
	tx := MakeTx()
	buf := tx.SerializeToByte()
	byteBuf := bytes.NewBuffer(buf)

	newTx := GhostTransaction{}
	newTx.Deserialize(byteBuf)
	newBuf := newTx.SerializeToByte()
	result := bytes.Compare(buf, newBuf)
	assert.Equal(t, 0, result, "bytes are different.")
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
	assert.Equal(t, 0, result, "binary is different.")
	assert.Equal(t, int(size), len(stream.Bytes()), "Size is different.")
	assert.Equal(t, output.Value, newOutput.Value, "Value is different.")
}

func TestTxInputSerializeDeserialize(t *testing.T) {
	input := MakeTxInput()
	size := input.Size()
	stream := mems.NewCapacity(int(size))
	input.Serialize(stream)
	byteBuf := bytes.NewBuffer(stream.Bytes())

	newInput := TxInput{}
	newInput.Deserialize(byteBuf)
	assert.Equal(t, int(size), len(stream.Bytes()), "Size is different.")
	assert.Equal(t, input.Sequence, newInput.Sequence, "Value is different.")
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
	assert.Equal(t, 0, result, "binary is different.")
	assert.Equal(t, int(size), len(stream.Bytes()), "Size is different.")
	assert.Equal(t, body.Nonce, newBody.Nonce, "Value is different.")
}

func MakeTxOutput() TxOutput {
	dummy := make([]byte, ghostBytes.PubKeySize)

	output := TxOutput{
		Addr:         dummy,
		BrokerAddr:   dummy,
		Value:        1212,
		ScriptSize:   ghostBytes.PubKeySize,
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
	txId := hash.Sum(nil)
	return GhostTransaction{txId, txBody}
}
