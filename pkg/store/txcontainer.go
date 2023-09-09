package store

import (
	"github.com/GhostNet-Dev/GhostNet-Core/libs/container"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gsql"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type TxContainer struct {
	gSql             gsql.GSql
	gCandidateSql    gsql.GCandidateSql
	CandidateTxPools *container.Queue
	CurrentPoolId    uint32
}

func NewTxContainer(g gsql.GSql, gCandidateSql gsql.GCandidateSql) *TxContainer {
	return &TxContainer{
		gSql:             g,
		gCandidateSql:    gCandidateSql,
		CandidateTxPools: container.NewQueue(),
	}
}

func (txContainer *TxContainer) Initialize() {
	txContainer.CurrentPoolId = txContainer.gCandidateSql.GetMaxPoolId()
}

func (txContainer *TxContainer) SaveTransaction(blockId uint32, tx *types.GhostTransaction, txIndexInBlock uint32) {
	txContainer.gSql.InsertTx(blockId, tx, types.NormalTx, txIndexInBlock)
}

func (txContainer *TxContainer) GetUnusedOutputList(txType types.TxOutputType, toAddr []byte) []types.PrevOutputParam {
	return txContainer.gSql.SelectUnusedOutputs(txType, toAddr)
}

func (txContainer *TxContainer) SearchStringOutputList(txType types.TxOutputType,
	toAddr, searchText []byte) []types.PrevOutputParam {
	return txContainer.gSql.SearchStringOutputs(txType, toAddr, searchText)
}

func (txContainer *TxContainer) SearchOutputList(txType types.TxOutputType, toAddr []byte) []types.PrevOutputParam {
	return txContainer.gSql.SearchOutputs(txType, toAddr)
}

func (txContainer *TxContainer) SearchOutput(txType types.TxOutputType, toAddr, uniqKey []byte) []types.PrevOutputParam {
	return txContainer.gSql.SearchOutput(txType, toAddr, uniqKey)
}

// to find the latest account
func (txContainer *TxContainer) SelectOutputList(txType types.TxOutputType, start, count int) []types.PrevOutputParam {
	return txContainer.gSql.SelectOutputs(txType, start, count)
}
func (txContainer *TxContainer) CheckExistTxId(txId []byte) bool {
	return txContainer.gSql.CheckExistTxId(txId)
}

func (txContainer *TxContainer) GetTx(txId []byte) (*types.GhostTransaction, uint32) {
	return txContainer.gSql.SelectTx(txId)
}

func (txContainer *TxContainer) GetDataTx(txId []byte) *types.GhostDataTransaction {
	return txContainer.gSql.SelectData(txId)
}

func (txContainer *TxContainer) CheckRefExist(refTxId []byte, outIndex uint32, notTxId []byte) bool {
	return txContainer.gSql.CheckExistRefOutput(refTxId, outIndex, notTxId)
}

func (txContainer *TxContainer) CheckExistFsRoot(nickname []byte) bool {
	return txContainer.gSql.CheckExistFsRoot(nickname)
}

func (txContainer *TxContainer) GetMaxLogicalAddress(owner []byte) (uint64, error) {
	return txContainer.gSql.GetMaxLogicalAddress(owner)
}

func (txContainer *TxContainer) GetNicknameToAddress(nickname string) []byte {
	return txContainer.gSql.GetNicknameToAddress([]byte(nickname))
}
