package txs

import (
	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type TXs struct {
	gScript        *gvm.GScript
	gVmExe         *gvm.GVM
	blockContainer *store.BlockContainer
}

func NewTXs(g *gvm.GScript, b *store.BlockContainer, e *gvm.GVM) *TXs {
	return &TXs{
		gScript:        g,
		blockContainer: b,
		gVmExe:         e,
	}
}

func (txs *TXs) InkTheContract(tx *types.GhostTransaction,
	ghostAddr *gcrypto.GhostAddress) *types.GhostTransaction {
	txs.gScript.MakeScriptSigExecuteUnlock(tx, ghostAddr)
	tx.TxId = tx.GetHashKey()
	return tx
}

func (txs *TXs) MakeDataTx(logicalAddr uint64, data []byte) *types.GhostDataTransaction {
	dummy := make([]byte, gbytes.HashSize)
	dataTx := types.GhostDataTransaction{
		TxId:           dummy,
		LogicalAddress: logicalAddr,
		Data:           data,
		DataSize:       uint32(len(data)),
	}
	dataTx.TxId = dataTx.GetHashKey()
	return &dataTx
}

func (txs *TXs) MakeContractTx(prev types.PrevOutputParam,
	next []types.NextOutputParam) *types.GhostTransaction {
	dummy := make([]byte, gbytes.HashSize)
	inputs := []types.TxInput{
		{
			PrevOut:    prev.VOutPoint,
			Sequence:   0xffffffff,
			ScriptSize: prev.Vout.ScriptSize,
			ScriptSig:  prev.Vout.ScriptPubKey,
		},
	}

	outputs := []types.TxOutput{}
	for _, newScript := range next {
		output := types.TxOutput{
			Addr:         newScript.RecvAddr,
			BrokerAddr:   newScript.Broker,
			Value:        newScript.TransferCoin,
			ScriptSize:   uint32(len(newScript.OutputScript)),
			ScriptPubKey: newScript.OutputScript,
		}
		outputs = append(outputs, output)
	}

	return &types.GhostTransaction{
		TxId: dummy,
		Body: types.TxBody{
			Vin:           inputs,
			InputCounter:  uint32(len(inputs)),
			Vout:          outputs,
			OutputCounter: uint32(len(outputs)),
		},
	}
}

func (txs *TXs) MakeTransaction(info TransferCoinInfo, prev map[types.TxOutputType][]types.PrevOutputParam,
	next map[types.TxOutputType][]types.NextOutputParam) *types.GhostTransaction {
	dummy := make([]byte, gbytes.HashSize)
	inputs := []types.TxInput{}
	outputs := []types.TxOutput{}

	// prev가 없을 수도 있다.
	for txType, nextOutputParams := range next {
		input, output := txs.MakeInputOutput(txType, info, prev[txType], nextOutputParams)
		inputs = append(inputs, input...)
		outputs = append(outputs, output...)
	}

	tx := &types.GhostTransaction{
		TxId: dummy,
		Body: types.TxBody{
			InputCounter:  uint32(len(inputs)),
			Vin:           inputs,
			OutputCounter: uint32(len(outputs)),
			Vout:          outputs,
			Nonce:         0,
			LockTime:      0,
		},
	}
	return tx
}

func (txs *TXs) MakeInputOutput(txType types.TxOutputType, info TransferCoinInfo, prev []types.PrevOutputParam,
	next []types.NextOutputParam) (inputs []types.TxInput, outputs []types.TxOutput) {
	var totalCoin uint64 = 0
	var transferCoin uint64 = 0

	// for rootfs tx, it's not necessary previous tx
	if len(prev) == 0 && txType == types.TxTypeFSRoot {
		dummyBuf4 := make([]byte, 4)
		dummyHash := make([]byte, gbytes.HashSize)
		inputs = append(inputs, types.TxInput{
			PrevOut: types.TxOutPoint{
				TxId: dummyHash,
			},
			Sequence:   0xFFFFFFFF,
			ScriptSize: uint32(len(dummyBuf4)),
			ScriptSig:  dummyBuf4,
		})
	} else {
		// for normal tx
		for _, outputParam := range prev {
			input := types.TxInput{
				PrevOut:    outputParam.VOutPoint,
				Sequence:   0xFFFFFFFF,
				ScriptSize: outputParam.Vout.ScriptSize,
				ScriptSig:  outputParam.Vout.ScriptPubKey, // 서명후 새로 생성된 서명으로 교체된다.
			}
			inputs = append(inputs, input)
			totalCoin += outputParam.Vout.Value
		}
	}

	for _, newOutput := range next {
		output := types.TxOutput{
			Addr:         newOutput.RecvAddr,
			BrokerAddr:   newOutput.Broker,
			Value:        newOutput.TransferCoin,
			Type:         txType,
			ScriptSize:   uint32(len(newOutput.OutputScript)),
			ScriptPubKey: newOutput.OutputScript,
		}
		outputs = append(outputs, output)
		transferCoin += newOutput.TransferCoin
	}

	if txType == types.TxTypeCoinTransfer && totalCoin > transferCoin {
		script := gvm.MakeLockScriptOut(info.MyWallet.MyPubKey())
		output := types.TxOutput{
			Addr:         info.MyWallet.MyPubKey(),
			BrokerAddr:   info.Broker,
			Type:         txType,
			Value:        totalCoin - transferCoin,
			ScriptPubKey: script,
			ScriptSize:   uint32(len(script)),
		}
		outputs = append(outputs, output)
	}

	return inputs, outputs
}
