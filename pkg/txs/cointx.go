package txs

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
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
				OutputScript: txs.gScript.MakeLockScriptOut(info.ToAddr),
				TransferCoin: info.TransferCoin,
			},
			{
				TxType:       types.TxTypeCoinTransfer,
				RecvAddr:     info.FeeAddr,
				Broker:       info.FeeBroker,
				OutputScript: txs.gScript.MakeLockScriptOut(info.FeeAddr),
				TransferCoin: uint64(float64(info.TransferCoin) * FeeRatio),
			},
		},
	}
	return txs.MakeTransaction(info, info.Prevs, nextOutputParam)
}
