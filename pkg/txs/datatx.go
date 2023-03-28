package txs

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

func (txs *TXs) CreateDataTx(info TransferTxInfo, logicalAddr uint64,
	data []byte) (*types.GhostTransaction, *types.GhostDataTransaction) {
	dataTx := txs.MakeDataTx(logicalAddr, data)
	dataTxId := dataTx.GetHashKey()
	nextOutputParam := map[types.TxOutputType][]types.NextOutputParam{
		types.TxTypeDataStore: {
			{
				TxType:         types.TxTypeDataStore,
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
				OutputScript: gvm.MakeLockScriptOut(info.FeeAddr),
				TransferCoin: 0, // 현재는 free
			},
		},
	}
	return txs.MakeTransaction(info, info.Prevs, nextOutputParam), dataTx
}
