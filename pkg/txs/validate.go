package txs

import (
	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
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

// TODO: reduce duplicated code
// 1. format check
// 2. db check
// 3. gvm excute
func TransactionDefaultValidation(tx *types.GhostTransaction, dataTx *types.GhostDataTransaction,
	gVmExe *gvm.GVM) *TxChkResult {
	var gFuncParam []gvm.GFuncParam
	if tx.Body.InputCounter != uint32(len(tx.Body.Vin)) {
		return &TxChkResult{TxChkResult_CounterMismatch}
	}
	if tx.Body.OutputCounter != uint32(len(tx.Body.Vout)) {
		return &TxChkResult{TxChkResult_CounterMismatch}
	}

	if tx.Body.Vout[0].Type != types.TxTypeFSRoot {
		return &TxChkResult{TxChkResult_Success}
	}

	dummyBuf4 := make([]byte, 4)
	input := tx.Body.Vin[0]
	scriptSig := input.ScriptSig
	input.ScriptSig = dummyBuf4
	input.ScriptSize = uint32(len(dummyBuf4))

	gFuncParam = append(gFuncParam, gvm.GFuncParam{
		InputSig:      scriptSig,
		ScriptPubbKey: tx.Body.Vout[0].ScriptPubKey,
		TxType:        types.TxTypeFSRoot,
	})

	dummy := make([]byte, gbytes.HashSize)
	txId := tx.TxId
	tx.TxId = dummy

	if !gVmExe.ExecuteGFunction(tx.SerializeToByte(), gFuncParam) {
		return &TxChkResult{TxChkResult_ScriptError}
	}

	tx.TxId = txId

	return &TxChkResult{TxChkResult_Success}

}

func (txs *TXs) TransactionValidation(tx *types.GhostTransaction, dataTx *types.GhostDataTransaction,
	txContainer *store.TxContainer) *TxChkResult {
	var transferCoin, getherCoin uint64 = 0, 0
	var gFuncParam []gvm.GFuncParam

	if tx.Body.InputCounter != uint32(len(tx.Body.Vin)) {
		return &TxChkResult{TxChkResult_CounterMismatch}
	}
	if tx.Body.OutputCounter != uint32(len(tx.Body.Vout)) {
		return &TxChkResult{TxChkResult_CounterMismatch}
	}

	if tx.Body.Vout[0].Type == types.TxTypeFSRoot {
		dummyBuf4 := make([]byte, 4)
		input := tx.Body.Vin[0]
		scriptSig := input.ScriptSig
		input.ScriptSig = dummyBuf4
		input.ScriptSize = uint32(len(dummyBuf4))
		nickname := string(tx.Body.Vout[0].ScriptEx)
		if txContainer.CheckExistFsRoot([]byte(nickname)) {
			return &TxChkResult{TxChkResult_AlreadyExist}
		}

		gFuncParam = append(gFuncParam, gvm.GFuncParam{
			InputSig:      scriptSig,
			ScriptPubbKey: tx.Body.Vout[0].ScriptPubKey,
			TxType:        types.TxTypeFSRoot,
		})
	} else {
		for _, output := range tx.Body.Vout {
			if output.Type == types.TxTypeCoinTransfer {
				transferCoin += output.Value
			}
		}

		for _, input := range tx.Body.Vin {
			prevOutpointer := input.PrevOut
			// Check Validate TxId
			prevTx, _ := txContainer.GetTx(prevOutpointer.TxId)
			if prevTx == nil {
				return &TxChkResult{TxChkResult_MissingRefTx}
			}
			// Check Script Format
			prevOutput := prevTx.Body.Vout[prevOutpointer.TxOutIndex]
			if prevOutput.ScriptSize != uint32(len(prevOutput.ScriptPubKey)) {
				return &TxChkResult{TxChkResult_FormatMismatch}
			}

			if !txContainer.CheckRefExist(prevOutpointer.TxId, prevOutpointer.TxOutIndex, tx.TxId) {
				return &TxChkResult{TxChkResult_MissingRefTx}
			}

			if input.ScriptSig == nil {
				return &TxChkResult{TxChkResult_FormatMismatch}
			}

			// check only transfer coin
			if prevOutput.Type == types.TxTypeCoinTransfer {
				getherCoin += prevOutput.Value
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
	}

	if transferCoin != getherCoin {
		return &TxChkResult{TxChkResult_Error}
	}

	dummy := make([]byte, gbytes.HashSize)
	txId := tx.TxId
	tx.TxId = dummy

	if !txs.gVmExe.ExecuteGFunction(tx.SerializeToByte(), gFuncParam) {
		return &TxChkResult{TxChkResult_ScriptError}
	}

	tx.TxId = txId

	return &TxChkResult{TxChkResult_Success}
}

func (txs *TXs) TransactionMergeValidation(tx *types.GhostTransaction, dataTx *types.GhostDataTransaction,
	txContainer *store.TxContainer, mergeTxContainer *store.TxContainer) *TxChkResult {
	var transferCoin, getherCoin uint64 = 0, 0
	var gFuncParam []gvm.GFuncParam

	if tx.Body.InputCounter != uint32(len(tx.Body.Vin)) {
		return &TxChkResult{TxChkResult_CounterMismatch}
	}
	if tx.Body.OutputCounter != uint32(len(tx.Body.Vout)) {
		return &TxChkResult{TxChkResult_CounterMismatch}
	}

	if tx.Body.Vout[0].Type == types.TxTypeFSRoot {
		dummyBuf4 := make([]byte, 4)
		input := tx.Body.Vin[0]
		scriptSig := input.ScriptSig
		input.ScriptSig = dummyBuf4
		input.ScriptSize = uint32(len(dummyBuf4))

		gFuncParam = append(gFuncParam, gvm.GFuncParam{
			InputSig:      scriptSig,
			ScriptPubbKey: tx.Body.Vout[0].ScriptPubKey,
			TxType:        types.TxTypeFSRoot,
		})
	} else {
		for _, output := range tx.Body.Vout {
			if output.Type == types.TxTypeCoinTransfer {
				transferCoin += output.Value
			}
		}

		for _, input := range tx.Body.Vin {
			prevOutpointer := input.PrevOut
			// Check Validate TxId
			prevTx, _ := txContainer.GetTx(prevOutpointer.TxId)
			if prevTx == nil {
				//TODO: need to separate between normal validate and merge validate
				if mergeTxContainer != nil {
					if prevTx, _ = mergeTxContainer.GetTx(prevOutpointer.TxId); prevTx == nil {
						return &TxChkResult{TxChkResult_MissingRefTx}
					}
				} else {
					return &TxChkResult{TxChkResult_MissingRefTx}
				}
			}
			// Check Script Format
			prevOutput := prevTx.Body.Vout[prevOutpointer.TxOutIndex]
			if prevOutput.ScriptSize != uint32(len(prevOutput.ScriptPubKey)) {
				return &TxChkResult{TxChkResult_FormatMismatch}
			}

			if !txContainer.CheckRefExist(prevOutpointer.TxId, prevOutpointer.TxOutIndex, tx.TxId) {
				//TODO: need to separate between normal validate and merge validate
				if mergeTxContainer != nil && mergeTxContainer.CheckRefExist(prevOutpointer.TxId, prevOutpointer.TxOutIndex, tx.TxId) {
				} else {
					return &TxChkResult{TxChkResult_MissingRefTx}
				}
			}

			if input.ScriptSig == nil {
				return &TxChkResult{TxChkResult_FormatMismatch}
			}

			// check only transfer coin
			if prevOutput.Type == types.TxTypeCoinTransfer {
				getherCoin += prevOutput.Value
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
	}

	if transferCoin != getherCoin {
		return &TxChkResult{TxChkResult_Error}
	}

	dummy := make([]byte, gbytes.HashSize)
	txId := tx.TxId
	tx.TxId = dummy

	if !txs.gVmExe.ExecuteGFunction(tx.SerializeToByte(), gFuncParam) {
		return &TxChkResult{TxChkResult_ScriptError}
	}

	tx.TxId = txId

	return &TxChkResult{TxChkResult_Success}
}
