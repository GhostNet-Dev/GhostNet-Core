package store

import (
	"bytes"
	_ "embed"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

var (
	//go:embed genesisblock
	genesisBlockByte   []byte
	genesisPairedBlock *types.PairedBlock
	adamsPubKey        []byte
)

func GenesisBlock() *types.PairedBlock {
	if genesisPairedBlock == nil {
		genesisPairedBlock = new(types.PairedBlock)
		seriBuf := bytes.NewBuffer(genesisBlockByte)
		genesisPairedBlock.Deserialize(seriBuf)
	}
	return genesisPairedBlock
}

func AdamsAddress() []byte {
	if adamsPubKey != nil {
		return adamsPubKey
	}
	block := GenesisBlock().Block
	sigPubKey := block.Header.BlockSignature.PubKey
	pubKey := gcrypto.TranslateSigPubTo160PubKey(sigPubKey)
	txs := block.Transaction
	result := 0

	for _, tx := range txs {
		output := tx.Body.Vout[0]
		if output.Type == types.TxTypeFSRoot {
			result = bytes.Compare(pubKey, output.Addr)
			if result == 0 {
				adamsPubKey = pubKey
				break
			}
		}
	}
	if result != 0 {
		return nil
	}

	return adamsPubKey
}
