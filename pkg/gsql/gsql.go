package gsql

import (
	types "github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

// GSql sql Instance
type GSql interface {
	OpenSQL(path string) error
	CloseSQL()
	CreateTable(schemaFile string) error
	DropTable()
	InsertTx(blockId uint32, tx types.GhostTransaction, txType uint32, txIndexInBlock uint32)
	InsertDataTx(blockId uint32, dataTx types.GhostDataTransaction, txIndexInBlock uint32)
	SelectTx(TxId []byte, txType uint32) *types.GhostTransaction
	SelectData(TxId []byte) *types.GhostDataTransaction
	InsertBlock(pair types.PairedBlock)
	SelectBlock(blockId uint32) *types.PairedBlock
}

// NewGSql sql instance를 생성한다.
func NewGSql(sqlType string) GSql {
	var gSql GSql
	switch sqlType {
	case "postgres":
	case "sqlite3":
		gSql = new(GSqlite3)
	}

	return gSql
}
