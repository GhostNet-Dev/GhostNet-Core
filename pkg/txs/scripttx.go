package txs

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

func (txs *TXs) CreateScriptTx(info TransferTxInfo, uniqKey []byte,
	data []byte) (*types.GhostTransaction, *types.GhostDataTransaction) {
	dataTx := txs.MakeDataTx(uniqKey, data)
	dataTxId := dataTx.TxId
	nextOutputParam := map[types.TxOutputType][]types.NextOutputParam{
		types.TxTypeScript: {
			{
				TxType:         types.TxTypeScript,
				RecvAddr:       info.ToAddr,
				Broker:         info.Broker,
				OutputScript:   gvm.MakeDataMapping(info.ToAddr),
				OutputScriptEx: dataTxId,
				TransferCoin:   0,
			},
			{
				TxType:       types.TxTypeCoinTransfer,
				RecvAddr:     info.FeeAddr,
				Broker:       info.FeeBroker,
				OutputScript: gvm.MakeDataMapping(info.ToAddr),
				TransferCoin: 0, // 현재는 free
			},
		},
	}
	return txs.MakeTransaction(info, info.Prevs, nextOutputParam), dataTx
}

func (txs *TXs) CreateScriptDataTx(info TransferTxInfo, uniqKey []byte,
	data []byte) (*types.GhostTransaction, *types.GhostDataTransaction) {
	dataTx := txs.MakeDataTx(uniqKey, data)
	dataTxId := dataTx.TxId
	nextOutputParam := map[types.TxOutputType][]types.NextOutputParam{
		types.TxTypeScriptStore: {
			{
				TxType:         types.TxTypeScriptStore,
				RecvAddr:       info.ToAddr,
				Broker:         info.Broker,
				OutputScript:   uniqKey,
				OutputScriptEx: dataTxId,
				TransferCoin:   0,
			},
			{
				TxType:       types.TxTypeCoinTransfer,
				RecvAddr:     info.FeeAddr,
				Broker:       info.FeeBroker,
				OutputScript: nil,
				TransferCoin: 0, // 현재는 free
			},
		},
	}
	return txs.MakeTransaction(info, info.Prevs, nextOutputParam), dataTx
}
