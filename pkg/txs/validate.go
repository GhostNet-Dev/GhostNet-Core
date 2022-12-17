package txs

import (
	gbytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

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
	if txResult == nil {
		return "success"
	}
	resultString := "error"
	switch txResult.result {
	case TxChkResult_Success:
		resultString = "success"
	}
	return resultString
}

func (txs *TXs) TransactionChecker(tx *types.GhostTransaction, dataTx *types.GhostDataTransaction,
	txContainer *store.TxContainer) *TxChkResult {
	var transferCoin, getherCoin uint64 = 0, 0
	var gFuncParam []gvm.GFuncParam

	if tx.Body.InputCounter != uint32(len(tx.Body.Vin)) {
		return &TxChkResult{TxChkResult_CounterMismatch}
	}
	if tx.Body.OutputCounter != uint32(len(tx.Body.Vout)) {
		return &TxChkResult{TxChkResult_CounterMismatch}
	}

	for _, input := range tx.Body.Vin {
		prevOutpointer := input.PrevOut
		// Check Validate TxId
		prevTx := txContainer.GetTx(prevOutpointer.TxId)
		if prevTx == nil {
			return &TxChkResult{TxChkResult_MissingRefTx}
		}
		// Check Script Format
		prevOutput := prevTx.Body.Vout[prevOutpointer.TxOutIndex]
		if prevOutput.ScriptSize != uint32(len(prevOutput.ScriptPubKey)) {
			return &TxChkResult{TxChkResult_FormatMismatch}
		}
		// Check Coin
		if prevOutput.Type == types.TxTypeCoinTransfer {
			getherCoin += prevOutput.Value
		}

		if prevOutput.Type != types.TxTypeFSRoot &&
			txContainer.CheckRefExist(prevOutpointer.TxId, prevOutpointer.TxOutIndex, tx.TxId) == true {
			return &TxChkResult{TxChkResult_MissingRefTx}
		}

		if input.ScriptSig == nil {
			return &TxChkResult{TxChkResult_FormatMismatch}
		}

		scriptSig := input.ScriptSig
		input.ScriptSig = prevOutput.ScriptPubKey
		input.ScriptSize = prevOutput.ScriptSize

		gFuncParam = append(gFuncParam, gvm.GFuncParam{
			InputSig:      scriptSig,
			ScriptPubbKey: prevOutput.ScriptPubKey,
			TxType:        prevOutput.Type,
		})
	}

	for _, output := range tx.Body.Vout {
		if output.Type == types.TxTypeCoinTransfer {
			transferCoin += output.Value
		}
	}

	if transferCoin != getherCoin {
		return &TxChkResult{TxChkResult_Error}
	}

	dummy := make([]byte, gbytes.HashSize)
	tx.TxId = dummy

	if txs.gVmExe.ExecuteGFunction(tx.SerializeToByte(), gFuncParam) == false {
		return &TxChkResult{TxChkResult_ScriptError}
	}

	return nil
}
