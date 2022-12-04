package txs

import (
	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/crypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type TXs struct {
	gScript        *gvm.GScript
	blockContainer *store.BlockContainer
}

func NewTXs(g *gvm.GScript, b *store.BlockContainer) *TXs {
	return &TXs{
		gScript:        g,
		blockContainer: b,
	}
}

func (txs *TXs) InkTheContract(tx *types.GhostTransaction,
	ghostAddr *crypto.GhostAddress) *types.GhostTransaction {
	txs.gScript.MakeScriptSigExecuteUnlock(tx, ghostAddr)
	tx.TxId = tx.GetHashKey()
	return tx
}

func (txs *TXs) MakeDataTx(logicalAddr uint64, data []byte) *types.GhostDataTransaction {
	dummy := make([]byte, ghostBytes.HashSize)
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
	dummy := make([]byte, ghostBytes.HashSize)
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

func (txs *TXs) MakeTransaction(info TransferCoinInfo, prev map[uint32][]types.PrevOutputParam,
	next map[uint32][]types.NextOutputParam) *types.GhostTransaction {
	dummy := make([]byte, ghostBytes.HashSize)
	inputs := []types.TxInput{}
	outputs := []types.TxOutput{}

	for txType, prevOutputParams := range prev {
		input, output := txs.MakeInputOutput(txType, info, prevOutputParams, next[txType])
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

func (txs *TXs) MakeInputOutput(txType uint32, info TransferCoinInfo, prev []types.PrevOutputParam,
	next []types.NextOutputParam) (inputs []types.TxInput, outputs []types.TxOutput) {
	var totalCoin uint64 = 0
	var transferCoin uint64 = 0

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

	for _, newOutput := range next {
		output := types.TxOutput{
			Addr:         newOutput.RecvAddr,
			BrokerAddr:   newOutput.Broker,
			Value:        newOutput.TransferCoin,
			ScriptSize:   uint32(len(newOutput.OutputScript)),
			ScriptPubKey: newOutput.OutputScript,
		}
		outputs = append(outputs, output)
		transferCoin += newOutput.TransferCoin
	}

	if txType == types.TxTypeCoinTransfer && totalCoin > transferCoin {
		script := txs.gScript.MakeLockScriptOut(info.MyWallet.MyPubKey())
		output := types.TxOutput{
			Addr:         info.MyWallet.MyPubKey(),
			BrokerAddr:   info.Broker,
			Value:        totalCoin - transferCoin,
			ScriptPubKey: script,
			ScriptSize:   uint32(len(script)),
		}
		outputs = append(outputs, output)
	}

	return inputs, outputs
}
