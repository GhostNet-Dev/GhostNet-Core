package gsql

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gsql/sqlite"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

// GSql sql Instance
type GSql interface {
	OpenSQL(path string, filename string) error
	CloseSQL()
	CreateTable() error
	DropTable()
	InsertTx(blockId uint32, tx *types.GhostTransaction, txType uint32, txIndexInBlock uint32)
	InsertDataTx(blockId uint32, dataTx *types.GhostDataTransaction, txIndexInBlock uint32)
	SelectTx(TxId []byte) (tx *types.GhostTransaction, blockId uint32)
	SelectData(TxId []byte) *types.GhostDataTransaction
	SelectUnusedOutputs(TxType types.TxOutputType, ToAddr []byte) []types.PrevOutputParam
	SelectOutputs(TxType types.TxOutputType, start, count int) []types.PrevOutputParam
	SelectOutputLatests(txType types.TxOutputType, toAddr, keyword []byte, start, count int) []types.PrevOutputParam
	SearchStringOutputs(txType types.TxOutputType, toAddr, keyword []byte) []types.PrevOutputParam
	SearchOutputs(TxType types.TxOutputType, ToAddr []byte) []types.PrevOutputParam
	SearchOutput(TxType types.TxOutputType, ToAddr, UniqKey []byte) []types.PrevOutputParam
	InsertBlock(pair *types.PairedBlock)
	SelectBlock(blockId uint32) *types.PairedBlock
	SelectBlockHeader(blockId uint32) (*types.GhostNetBlockHeader, *types.GhostNetDataBlockHeader)
	CheckExistBlockId(blockId uint32) bool
	CheckExistTxId(txId []byte) bool
	CheckExistTxBefore(txId []byte, curBlockId uint32) bool
	CheckExistRefOutput(refTxId []byte, outIndex uint32, notTxId []byte) bool
	CheckExistFsRoot(nickname []byte) bool
	CheckExistFsRootWithoutCurrentTx(nickname, txId []byte) bool
	GetBlockHeight() uint32
	SelectTxsPool(poolId uint32) []types.GhostTransaction
	SelectDataTxsPool(poolId uint32) []types.GhostDataTransaction
	DeleteBlock(blockId uint32)
	DeleteAfterTargetId(blockId uint32)
	GetMaxLogicalAddress(toAddr []byte) (uint64, error)
	GetNicknameToAddress(nickname []byte) []byte
	GetIssuedCoin(blockId uint32) uint64
}

type GCandidateSql interface {
	InsertCandidateTx(tx *types.GhostTransaction, poolId uint32)
	InsertCandidateDataTx(dataTx *types.GhostDataTransaction, poolId uint32)
	SelectCandidateTxCount() uint32
	DeleteCandidatePool(poolId uint32)
	DeleteCandidateTx(txId []byte)
	GetMinPoolId() uint32
	GetMaxPoolId() uint32
	UpdatePoolId(oldPoolId uint32, newPoolId uint32)
	CheckExistCandidateTxId(txId []byte) bool
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
		gSql = sqlite.NewGSqlite3()
	}

	if gSql == nil {
		log.Fatal("Failed to create db.")
	}

	return gSql
}

func NewMergeGSql(sqlType string) GSql {
	var gCandidate GSql
	switch sqlType {
	case "postgres":
	case "sqlite3":
		gCandidate = sqlite.NewMergeGSqlite()
	}

	if gCandidate == nil {
		log.Fatal("Failed to create db.")
	}

	return gCandidate
}

func NewAccountSql(sqlType string) MasterNodeStore {
	return &sqlite.GhostAccount{}
}
