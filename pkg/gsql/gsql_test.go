package gsql

import (
	"bytes"
	"crypto/sha256"
	"log"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"github.com/stretchr/testify/assert"
	mems "github.com/traherom/memstream"
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func TestSqlCreateTable(t *testing.T) {
	gSql := NewGSql("sqlite3")
	err := gSql.OpenSQL("./", "block.db")
	if err == nil {
		defer gSql.CloseSQL()
	}
	gSql.CreateTable("../../db.sqlite3.sql")
}

func TestSqlInsertAndSelectCheck(t *testing.T) {
	gSql := NewGSql("sqlite3")
	err := gSql.OpenSQL("./", "block.db")
	if err == nil {
		defer gSql.CloseSQL()
	}
	if err = gSql.CreateTable("../../db.sqlite3.sql"); err == nil {
		defer gSql.DropTable()
	}
	tx := MakeTx()
	gSql.InsertTx(0, &tx, 0, 0)
	newTx := gSql.SelectTx(tx.TxId)
	size := tx.Size()
	sizeNew := newTx.Size()

	assert.Equal(t, size, sizeNew, "db tx 크기가 다릅니다.")

	result := bytes.Compare(tx.SerializeToByte(), newTx.SerializeToByte())
	assert.Equal(t, 0, result, "tx가 다릅니다.")

}

func TestSqlInAndOutBlock(t *testing.T) {
	gSql := NewGSql("sqlite3")
	err := gSql.OpenSQL("./", "block.db")
	if err == nil {
		defer gSql.CloseSQL()
	}
	if err = gSql.CreateTable("../../db.sqlite3.sql"); err == nil {
		defer gSql.DropTable()
	}

	pair := MakePairBlock()
	gSql.InsertBlock(&pair)
	pairNew := gSql.SelectBlock(2)

	size := pair.Size()
	sizeNew := pairNew.Size()
	assert.Equal(t, size, sizeNew, "db block 크기가 다릅니다.")

	result := bytes.Compare(pair.SerializeToByte(), pairNew.SerializeToByte())
	assert.Equal(t, 0, result, "block이 다릅니다.")
}

func MakePairBlock() types.PairedBlock {
	dummy := make([]byte, 4)
	hash := sha256.New()
	hash.Write(dummy)
	key := hash.Sum((nil))

	pair := types.PairedBlock{
		Block: types.GhostNetBlock{
			Header: types.GhostNetBlockHeader{
				Id:                      2,
				Version:                 1,
				PreviousBlockHeaderHash: key,
				MerkleRoot:              key,
				DataBlockHeaderHash:     key,
				TimeStamp:               123,
				Bits:                    456,
				Nonce:                   789,
				AliceCount:              1,
				TransactionCount:        1,
				SignatureSize:           4,
				BlockSignature:          types.SigHash{},
			},
			Alice:       []types.GhostTransaction{MakeTx()},
			Transaction: []types.GhostTransaction{MakeTx()},
		},
		DataBlock: types.GhostNetDataBlock{
			Header: types.GhostNetDataBlockHeader{
				Id:                      2,
				Version:                 1,
				PreviousBlockHeaderHash: key,
				MerkleRoot:              key,
				Nonce:                   123,
				TransactionCount:        0,
			},
		},
	}
	return pair
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
	return types.GhostTransaction{TxId: txId, Body: txBody}
}
