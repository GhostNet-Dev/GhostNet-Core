package txs

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
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

func ExtractFSRootGParam(tx *types.GhostTransaction) ([]gvm.GFuncParam, []byte) {
	if tx.Body.Vout[0].Type != types.TxTypeFSRoot {
		log.Fatal("genesis has wrong transactions")
		return nil, nil
	}

	var gFuncParam []gvm.GFuncParam
	input := tx.Body.Vin[0]
	scriptSig := input.ScriptSig

	gFuncParam = append(gFuncParam, gvm.GFuncParam{
		InputSig:      scriptSig,
		ScriptPubbKey: tx.Body.Vout[0].ScriptPubKey,
		TxType:        types.TxTypeFSRoot,
	})

	dummyBuf4 := make([]byte, 4)
	dummyHash := make([]byte, gbytes.HashSize)
	tx.TxId = dummyHash

	tx.Body.Vin[0] = types.TxInput{
		PrevOut: types.TxOutPoint{
			TxId: dummyHash,
		},
		Sequence:   0xFFFFFFFF,
		ScriptSize: uint32(len(dummyBuf4)),
		ScriptSig:  dummyBuf4,
	}
	return gFuncParam, tx.Body.Vout[0].Addr
}
