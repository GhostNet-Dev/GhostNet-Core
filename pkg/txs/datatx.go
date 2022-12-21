package txs

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

func (txs *TXs) CreateDataTx(info TransferCoinInfo, nickname string) *types.GhostTransaction {
	nextOutputParam := map[uint32][]types.NextOutputParam{
		types.TxTypeDataTransfer: {
			{
				TxType:       types.TxTypeDataTransfer,
				RecvAddr:     info.ToAddr,
				Broker:       info.Broker,
				OutputScript: txs.gScript.MakeRootAccount(info.ToAddr, nickname),
				TransferCoin: 0,
			},
		},
	}
	return txs.MakeTransaction(info, info.Prevs, nextOutputParam)
}
