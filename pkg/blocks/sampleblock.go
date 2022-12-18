package blocks

import (
	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

func (blocks *Blocks) MakeGenesisBlock(creator []string) (*types.PairedBlock, map[string]*gcrypto.GhostAddress) {
	accountFile := map[string]*gcrypto.GhostAddress{}
	txs := blocks.txs
	tx, root := txs.MakeSampleRootAccount("Adam", nil)
	accountFile["Adam@"+root.GetPubAddress()+".ghost"] = root

	newTxs := []types.GhostTransaction{*tx}
	for _, name := range creator {
		tx, address := txs.MakeSampleRootAccount(name, root.Get160PubKey())
		newTxs = append(newTxs, *tx)
		accountFile[name+"@"+address.GetPubAddress()+".ghost"] = address
	}

	msg := make([]byte, gbytes.HashSize)
	msg2 := make([]byte, gbytes.HashSize)
	copy(msg, []byte("Was it a cat I saw?"))
	copy(msg2, []byte("I show you how deep the rabbit hole goes."))

	pair := &types.PairedBlock{
		Block:     *blocks.CreateGhostNetBlock(1, msg, msg2, root, root.Get160PubKey(), newTxs),
		DataBlock: *blocks.CreateGhostNetDataBlock(1, msg, nil),
	}

	return pair, accountFile
}
