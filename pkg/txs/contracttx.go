package txs

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

func (txs *TXs) CreateContractTx(info TransferTxInfo, dataTxId []byte) *types.GhostTransaction {
	nextOutputParam := map[types.TxOutputType][]types.NextOutputParam{
		types.TxTypeCoinTransfer: {
			{
				TxType:       types.TxTypeContract,
				RecvAddr:     info.ToAddr,
				Broker:       info.Broker,
				OutputScript: gvm.MakeContractScript(info.ToAddr, dataTxId),
				TransferCoin: info.TransferCoin,
			},
		},
	}
	return txs.MakeTransaction(info, info.Prevs, nextOutputParam)
}

func (txs *TXs) CompleteContractTx(info TransferTxInfo, dataTxId, token, dataOwner []byte) *types.GhostTransaction {
	nextOutputParam := map[types.TxOutputType][]types.NextOutputParam{
		types.TxTypeCoinTransfer: {
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
				RecvAddr:     dataOwner,
				Broker:       info.Broker,
				OutputScript: gvm.MakeLockScriptOut(dataOwner),
				TransferCoin: info.TransferCoin,
			},
		},
	}
	return txs.MakeTransaction(info, info.Prevs, nextOutputParam)
}
