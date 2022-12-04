package txs

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/crypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

const (
	FeeRatio = 0.0129
)

type TransferCoinInfo struct {
	MyWallet     crypto.Wallet
	ToAddr       []byte
	Broker       []byte
	FeeAddr      []byte
	FeeBroker    []byte
	Prevs        map[uint32][]types.PrevOutputParam
	TransferCoin uint64
}

func (txs *TXs) TransferCoin(info TransferCoinInfo) *types.GhostTransaction {
	nextOutputParam := map[uint32][]types.NextOutputParam{
		types.TxTypeCoinTransfer: {
			{
				TxType:       types.TxTypeCoinTransfer,
				RecvAddr:     info.ToAddr,
				Broker:       info.Broker,
				TransferCoin: info.TransferCoin,
			},
			{
				TxType:       types.TxTypeCoinTransfer,
				RecvAddr:     info.FeeAddr,
				Broker:       info.FeeBroker,
				TransferCoin: uint64(float64(info.TransferCoin) * FeeRatio),
			},
		},
	}
	return txs.MakeTransaction(info, info.Prevs, nextOutputParam)
}
