package txs

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

func (txs *TXs) MakeSampleRootAccount(nickname string, brokerAddr []byte) (*types.GhostTransaction, *gcrypto.GhostAddress) {
	address := gcrypto.GenerateKeyPair()
	toAddr := address.Get160PubKey()
	if nickname == "Adam" {
		brokerAddr = toAddr
	}

	tx := txs.CreateRootFsTx(TransferTxInfo{
		ToAddr:    toAddr,
		Broker:    brokerAddr,
		FeeAddr:   brokerAddr,
		FeeBroker: brokerAddr,
	}, nickname)

	/*
		outputScript := txs.gScript.MakeRootAccount(toAddr, nickname)
		dummyTxId := make([]byte, gbytes.HashSize)
		dummyBuf4 := make([]byte, 4)

		tx := &types.GhostTransaction{
			Body: types.TxBody{
				Vin: []types.TxInput{{
					ScriptSig:  dummyBuf4,
					ScriptSize: uint32(len(dummyBuf4)),
				}},
				InputCounter: 1,
				Vout: []types.TxOutput{{
					Addr:         toAddr,
					BrokerAddr:   brokerAddr,
					Type:         types.TxTypeFSRoot,
					Value:        0,
					ScriptPubKey: outputScript,
					ScriptSize:   uint32(len(outputScript)),
				}},
				OutputCounter: 1,
			},
			TxId: dummyTxId,
		}
	*/

	tx = txs.InkTheContract(tx, address)

	return tx, address
}
