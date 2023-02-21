package txs

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

const (
	FeeRatio = 0.0129
)

func (txs *TXs) TransferCoin(info TransferTxInfo) *types.GhostTransaction {
	nextOutputParam := map[types.TxOutputType][]types.NextOutputParam{
		types.TxTypeCoinTransfer: {
			{
				TxType:       types.TxTypeCoinTransfer,
				RecvAddr:     info.ToAddr,
				Broker:       info.Broker,
				OutputScript: gvm.MakeLockScriptOut(info.ToAddr),
				TransferCoin: info.TransferCoin,
			},
			{
				TxType:       types.TxTypeCoinTransfer,
				RecvAddr:     info.FeeAddr,
				Broker:       info.FeeBroker,
				OutputScript: gvm.MakeLockScriptOut(info.FeeAddr),
				TransferCoin: uint64(float64(info.TransferCoin) * FeeRatio),
			},
		},
	}
	return txs.MakeTransaction(info, info.Prevs, nextOutputParam)
}
