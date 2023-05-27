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
	makeInOutFunc  map[types.TxOutputType]func(TransferTxInfo, []types.PrevOutputParam, []types.NextOutputParam) ([]types.TxInput, []types.TxOutput)
}

func NewTXs(g *gvm.GScript, b *store.BlockContainer, e *gvm.GVM) *TXs {
	txs := &TXs{
		gScript:        g,
		blockContainer: b,
		gVmExe:         e,
		makeInOutFunc:  make(map[types.TxOutputType]func(TransferTxInfo, []types.PrevOutputParam, []types.NextOutputParam) ([]types.TxInput, []types.TxOutput)),
	}

	txs.makeInOutFunc[types.TxTypeDataStore] = func(info TransferTxInfo, prev []types.PrevOutputParam,
		next []types.NextOutputParam) (inputs []types.TxInput, outputs []types.TxOutput) {
		for _, outputParam := range prev {
			input := types.MakeTxInputFromOutputParam(&outputParam)
			inputs = append(inputs, *input)
		}

		for _, newOutput := range next {
			output := types.MakeTxOutputFromOutputParam(&newOutput)
			outputs = append(outputs, *output)
		}
		return inputs, outputs
	}

	txs.makeInOutFunc[types.TxTypeFSRoot] = func(info TransferTxInfo, prev []types.PrevOutputParam,
		next []types.NextOutputParam) (inputs []types.TxInput, outputs []types.TxOutput) {
		inputs = append(inputs, *types.MakeEmptyInput())
		for _, newOutput := range next {
			output := types.MakeTxOutputFromOutputParam(&newOutput)
			outputs = append(outputs, *output)
		}
		return inputs, outputs
	}
	txs.makeInOutFunc[types.TxTypeShare] = txs.makeInOutFunc[types.TxTypeFSRoot]

	txs.makeInOutFunc[types.TxTypeCoinTransfer] = func(info TransferTxInfo, prev []types.PrevOutputParam,
		next []types.NextOutputParam) (inputs []types.TxInput, outputs []types.TxOutput) {
		var totalCoin uint64 = 0
		var transferCoin uint64 = 0
		for _, outputParam := range prev {
			input := types.MakeTxInputFromOutputParam(&outputParam)
			inputs = append(inputs, *input)
			totalCoin += outputParam.Vout.Value
		}

		for _, newOutput := range next {
			output := types.MakeTxOutputFromOutputParam(&newOutput)
			outputs = append(outputs, *output)
			transferCoin += newOutput.TransferCoin
		}

		if totalCoin > transferCoin {
			script := gvm.MakeLockScriptOut(info.MyWallet.MyPubKey())
			output := types.TxOutput{
				Addr:         info.MyWallet.MyPubKey(),
				BrokerAddr:   info.Broker,
				Type:         types.TxTypeCoinTransfer,
				Value:        totalCoin - transferCoin,
				ScriptPubKey: script,
				ScriptSize:   uint32(len(script)),
			}
			outputs = append(outputs, output)
		}
		return inputs, outputs
	}

	return txs
}

type TransferTxInfo struct {
	MyWallet     *gcrypto.Wallet
	ToAddr       []byte
	Broker       []byte
	FeeAddr      []byte
	FeeBroker    []byte
	Prevs        map[types.TxOutputType][]types.PrevOutputParam
	TransferCoin uint64
}

func (txs *TXs) InkTheContract(tx *types.GhostTransaction,
	ghostAddr *gcrypto.GhostAddress) *types.GhostTransaction {
	txs.gScript.MakeScriptSigExecuteUnlock(tx, ghostAddr)
	tx.TxId = tx.GetHashKey()
	return tx
}

func (txs *TXs) MakeDataTx(logicalAddr []byte, data []byte) *types.GhostDataTransaction {
	dummy := make([]byte, gbytes.HashSize)
	logicalAddrBuf := make([]byte, gbytes.HashSize)
	copy(logicalAddrBuf, logicalAddr)

	dataTx := types.GhostDataTransaction{
		TxId:           dummy,
		LogicalAddress: logicalAddrBuf,
		Data:           data,
		DataSize:       uint32(len(data)),
	}
	dataTx.TxId = dataTx.GetHashKey()
	return &dataTx
}

func (txs *TXs) MakeTransaction(info TransferTxInfo, prev map[types.TxOutputType][]types.PrevOutputParam,
	next map[types.TxOutputType][]types.NextOutputParam) *types.GhostTransaction {
	dummy := make([]byte, gbytes.HashSize)
	inputs := []types.TxInput{}
	outputs := []types.TxOutput{}

	// prev가 없을 수도 있다.
	for txType, nextOutputParams := range next {
		input, output := txs.makeInOutFunc[txType](info, prev[txType], nextOutputParams)
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
