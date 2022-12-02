package txs

import (
	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

func MakeTransaction(info TransferCoinInfo, prev map[uint32][]types.PrevOutputParam,
	next map[uint32][]types.NextOutputParam) *types.GhostTransaction {
	dummy := make([]byte, ghostBytes.HashSize)
	inputs := []types.TxInput{}
	outputs := []types.TxOutput{}

	for txType, prevOutputParams := range prev {
		input, output := MakeInputOutput(txType, info, prevOutputParams, next[txType])
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

func MakeInputOutput(txType uint32, info TransferCoinInfo, prev []types.PrevOutputParam,
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
		output := types.TxOutput {
			Addr: info.MyWallet.MyPubKey(),
			BrokerAddr: info.Broker,
			Value: totalCoin - transferCoin,
			ScriptPubKey:,
			ScriptSize:,
		}		
	}
}
