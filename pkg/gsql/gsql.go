package gsql

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gsql/sqlite"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
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
	SelectBlockHeader(blockId uint32) (*types.GhostNetBlockHeader, *types.GhostNetDataBlockHeader)
	CheckExistTxId(txId []byte) bool
	CheckExistRefOutout(refTxId []byte, outIndex uint32, notTxId []byte) bool
	GetBlockHeight() uint32
	SelectTxsPool(poolId uint32) []types.GhostTransaction
	SelectDataTxsPool(poolId uint32) []types.GhostDataTransaction
	DeleteAfterTargetId(blockId uint32)
}

type GCandidateSql interface {
	InsertCandidateTx(tx *types.GhostTransaction, poolId uint32)
	InsertCandidateDataTx(dataTx *types.GhostDataTransaction, poolId uint32)
	SelectCandidateTxCount() uint32
	GetMinPoolId() uint32
	GetMaxPoolId() uint32
	UpdatePoolId(oldPoolId uint32, newPoolId uint32)
}

type MasterNodeStore interface {
	GetMasterNodeList() []*ptypes.GhostUser
	GetMasterNodeSearch(pubKey string) []*ptypes.GhostUser
	GetMasterNodeSearchPick(pubKey string) *ptypes.GhostUser
}

// NewGSql sql instance를 생성한다.
func NewGSql(sqlType string) GSql {
	var gSql GSql
	switch sqlType {
	case "postgres":
	case "sqlite3":
		gSql = sqlite.GSqlite
	}

	if gSql == nil {
		log.Fatal("Failed to create db.")
	}

	return gSql
}

func NewGCandidateSql(sqlType string) GCandidateSql {
	var gCandidate GCandidateSql
	switch sqlType {
	case "postgres":
	case "sqlite3":
		gCandidate = sqlite.GSqlite
	}

	if gCandidate == nil {
		log.Fatal("Failed to create db.")
	}

	return gCandidate
}

func NewAccountSql(sqlType string) MasterNodeStore {
	return &sqlite.GhostAccount{}
}
