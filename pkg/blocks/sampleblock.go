package blocks

import (
	"github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

func (blocks *Blocks) MakeGenesisBlock(creator []string) (*types.GhostNetBlock, map[string][]byte) {
	accountFile := map[string][]byte{}
	txs := blocks.txs
	tx, broker := txs.MakeSampleRootAccount("Adam", nil)
	accountFile["Adam@"+broker.GetPubAddress()+".ghost"] = broker.PrivateKeySerialize()

	newTxs := []types.GhostTransaction{*tx}
	for _, name := range creator {
		tx, address := txs.MakeSampleRootAccount(name, broker.Get160PubKey())
		newTxs = append(newTxs, *tx)
		accountFile[name+"@"+address.GetPubAddress()+".ghost"] = address.PrivateKeySerialize()
	}

	msg := make([]byte, bytes.HashSize)
	msg2 := make([]byte, bytes.HashSize)
	copy(msg, []byte("Was it a cat I saw?"))
	copy(msg2, []byte("I show you how deep the rabbit hole goes."))

	return blocks.CreateGhostNetBlock(1, msg, msg2, broker, newTxs), accountFile
}
