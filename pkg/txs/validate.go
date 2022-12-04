package txs

import "github.com/GhostNet-Dev/GhostNet-Core/pkg/types"

const (
	TxChkResult_Success = iota
	TxChkResult_Pending
	TxChkResult_Error
	TxChkResult_CounterMismatch
	TxChkResult_MissingRefTx
	TxChkResult_FormatMismatch
	TxChkResult_MappingMismatch
	TxChkResult_DoubleSpending
	TxChkResult_DoubleSpendingInCandi
	TxChkResult_ScriptError
	TxChkResult_MasterReject
	TxChkResult_AlreadyExist
)

type TxChkResult struct {
	result uint32
}

func (txResult *TxChkResult) Result() bool {
	return txResult.result == TxChkResult_Success
}

func (txResult *TxChkResult) Error() string {
	resultString := "error"
	switch txResult.result {
	case TxChkResult_Success:
		resultString = "success"
	}
	return resultString
}

func (txs *TXs) TransactionChecker(tx *types.GhostTransaction, dataTx *types.GhostDataTransaction) *TxChkResult {
	if tx.Body.InputCounter != uint32(len(tx.Body.Vin)) {
		return &TxChkResult{TxChkResult_CounterMismatch}
	}
	if tx.Body.OutputCounter != uint32(len(tx.Body.Vout)) {
		return &TxChkResult{TxChkResult_CounterMismatch}
	}

	for _, input := range tx.Body.Vin {
		prevOutput := input.PrevOut
		prevTx := txs.blockContainer.GetTx(prevOutput.TxId)
		if prevTx == nil {
			return &TxChkResult{TxChkResult_MissingRefTx}
		}
		loadedPrevOutput := prevTx.Body.Vout[prevOutput.TxOutIndex]
		if loadedPrevOutput.ScriptSize != uint32(len(loadedPrevOutput.ScriptPubKey)) {
			return &TxChkResult{TxChkResult_FormatMismatch}
		}
	}

	return nil
}
