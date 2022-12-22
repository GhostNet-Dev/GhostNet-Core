package txs

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

func (txs *TXs) CreateContractTx(info TransferCoinInfo, nickname string) *types.GhostTransaction {
	nextOutputParam := map[types.TxOutputType][]types.NextOutputParam{
		types.TxTypeContract: {
			{
				TxType:       types.TxTypeContract,
				RecvAddr:     info.ToAddr,
				Broker:       info.Broker,
				OutputScript: txs.gScript.MakeRootAccount(info.ToAddr, nickname),
				TransferCoin: 0,
			},
		},
	}
	return txs.MakeTransaction(info, info.Prevs, nextOutputParam)
}
