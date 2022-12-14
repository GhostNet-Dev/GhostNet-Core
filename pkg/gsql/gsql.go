package gsql

import (
	"log"

	types "github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

// GSql sql Instance
type GSql interface {
	OpenSQL(path string) error
	CloseSQL()
	CreateTable(schemaFile string) error
	DropTable()
	InsertTx(blockId uint32, tx *types.GhostTransaction, txType uint32, txIndexInBlock uint32)
	InsertDataTx(blockId uint32, dataTx *types.GhostDataTransaction, txIndexInBlock uint32)
	SelectTx(TxId []byte) *types.GhostTransaction
	SelectData(TxId []byte) *types.GhostDataTransaction
	SelectUnusedOutputs(TxType types.TxOutputType, ToAddr []byte) []types.PrevOutputParam
	InsertBlock(pair *types.PairedBlock)
	SelectBlock(blockId uint32) *types.PairedBlock
	CheckExistTxId(txId []byte) bool
	CheckExistRefOutout(refTxId []byte, outIndex uint32, notTxId []byte) bool
	GetBlockHeight() uint32
	SelectTxsPool(poolId uint32) []types.GhostTransaction
	SelectDataTxsPool(poolId uint32) []types.GhostDataTransaction
	SelectCandidateTxCount() uint32
	GetMinPoolId() uint32
	GetMaxPoolId() uint32
	UpdatePoolId(oldPoolId uint32, newPoolId uint32)
}

// NewGSql sql instance를 생성한다.
func NewGSql(sqlType string) GSql {
	var gSql GSql
	switch sqlType {
	case "postgres":
	case "sqlite3":
		gSql = new(GSqlite3)
	}

	if gSql == nil {
		log.Fatal("Failed to create db.")
	}

	return gSql
}
