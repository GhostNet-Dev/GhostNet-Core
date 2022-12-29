package txs

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

func (txs *TXs) CreateDataTx(info TransferCoinInfo, nickname string) *types.GhostTransaction {
	nextOutputParam := map[types.TxOutputType][]types.NextOutputParam{
		types.TxTypeDataTransfer: {
			{
				TxType:       types.TxTypeDataTransfer,
				RecvAddr:     info.ToAddr,
				Broker:       info.Broker,
				OutputScript: gvm.MakeRootAccount(info.ToAddr, nickname),
				TransferCoin: 0,
			},
		},
	}
	return txs.MakeTransaction(info, info.Prevs, nextOutputParam)
}
