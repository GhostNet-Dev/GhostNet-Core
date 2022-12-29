package txs

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

func (txs *TXs) CreateRootFsTx(info TransferCoinInfo, nickname string) *types.GhostTransaction {
	nextOutputParam := map[types.TxOutputType][]types.NextOutputParam{
		types.TxTypeFSRoot: {
			{
				TxType:       types.TxTypeFSRoot,
				RecvAddr:     info.ToAddr,
				Broker:       info.Broker,
				OutputScript: gvm.MakeRootAccount(info.ToAddr, nickname),
				TransferCoin: 0,
			},
		},
	}
	return txs.MakeTransaction(info, info.Prevs, nextOutputParam)
}
