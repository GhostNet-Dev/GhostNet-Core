package gsql

import (
	"bytes"
	"crypto/sha256"
	"testing"

	types "github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"github.com/stretchr/testify/assert"
	mems "github.com/traherom/memstream"
)

func TestSqlCreateTable(t *testing.T) {
	gSql := NewGSql("sqlite3")
	gSql.OpenSQL("./")
	gSql.CreateTable("../../db.sqlite3.sql")
}

func TestSqlInsertAndSelectCheck(t *testing.T) {
	gSql := NewGSql("sqlite3")
	gSql.OpenSQL("./")
	gSql.CreateTable("../../db.sqlite3.sql")
	tx := MakeTx()
	gSql.InsertTx(0, tx, 0, 0)
	newTx := gSql.SelectTx(tx.TxId, 0)
	size := tx.Size()
	newSize := newTx.Size()
	assert.Equal(t, size, newSize, "db tx 크기가 다릅니다.")

	stream := mems.NewCapacity(int(size))
	tx.Serialize(stream)
	newStream := mems.NewCapacity(int(size))
	newTx.Serialize(newStream)
	result := bytes.Compare(stream.Bytes(), newStream.Bytes())
	assert.Equal(t, 0, result, "tx가 다릅니다.")

}

func MakeTxOutput() types.TxOutput {
	dummy := make([]byte, 4)
	hash := sha256.New()
	hash.Write(dummy)
	key := hash.Sum((nil))

	output := types.TxOutput{
		Addr:         key,
		BrokerAddr:   key,
		Type:         0,
		Value:        1212,
		ScriptSize:   4,
		ScriptPubKey: dummy,
	}
	return output
}

func MakeTxInput() types.TxInput {
	dummy := make([]byte, 4)
	hash := sha256.New()
	hash.Write(dummy)
	key := hash.Sum((nil))

	input := types.TxInput{
		PrevOut: types.TxOutPoint{
			TxId:       key,
			TxOutIndex: 0,
		},
		Sequence:   3232,
		ScriptSize: 4,
		ScriptSig:  dummy,
	}
	return input
}

func MakeTxBody() types.TxBody {
	return types.TxBody{
		InputCounter: 2,
		Vin: []types.TxInput{
			MakeTxInput(),
			MakeTxInput(),
		},
		OutputCounter: 1,
		Vout: []types.TxOutput{
			MakeTxOutput(),
		},
		Nonce:    2233,
		LockTime: 1234,
	}
}

func MakeTx() types.GhostTransaction {
	txBody := MakeTxBody()
	stream := mems.NewCapacity(int(txBody.Size()))
	txBody.Serialize(stream)
	hash := sha256.New()
	hash.Write(stream.Bytes())
	txId := hash.Sum((nil))
	return types.GhostTransaction{txId, txBody}
}
