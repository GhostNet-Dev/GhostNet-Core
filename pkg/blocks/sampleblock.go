package blocks

import (
	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

var (
	creator = []string{
		"Alice", "Bob", "Carol", "Dave", "Eve", "Frank", "Grace",
		"Heidi", "Ivan", "Judy", "Michael", "Niaj", "Oscar", "Peggy",
		"Root", "Sybil", "Theo", "Terry", "Victor", "Walter", "Wendy",
	}
)

func (blocks *Blocks) MakeGenesisBlock(saveCall func(string, *gcrypto.GhostAddress)) *types.PairedBlock {
	txs := blocks.txs
	tx, root := txs.MakeSampleRootAccount("Adam", nil)
	saveCall("Adam", root)

	newTxs := []types.GhostTransaction{*tx}
	for _, name := range creator {
		tx, address := txs.MakeSampleRootAccount(name, root.Get160PubKey())
		newTxs = append(newTxs, *tx)
		saveCall(name, address)
	}

	msg := make([]byte, gbytes.HashSize)
	msg2 := make([]byte, gbytes.HashSize)
	copy(msg, []byte("Was it a cat I saw?"))
	copy(msg2, []byte("I show you how deep the rabbit hole goes."))

	pair := &types.PairedBlock{
		Block:     *blocks.CreateGhostNetBlock(1, msg, msg2, root, root.Get160PubKey(), newTxs),
		DataBlock: *blocks.CreateGhostNetDataBlock(1, msg, nil),
	}

	return pair
}
