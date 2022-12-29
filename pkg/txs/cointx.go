package txs

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

const (
	FeeRatio = 0.0129
)

type TransferCoinInfo struct {
	MyWallet     gcrypto.Wallet
	ToAddr       []byte
	Broker       []byte
	FeeAddr      []byte
	FeeBroker    []byte
	Prevs        map[types.TxOutputType][]types.PrevOutputParam
	TransferCoin uint64
}

func (txs *TXs) TransferCoin(info TransferCoinInfo) *types.GhostTransaction {
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
