package store

import (
	"bytes"
	_ "embed"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

//go:embed GenesisBlocks
var genesisBlock []byte

func GenesisBlock() *types.GhostNetBlock {
	var pair *types.PairedBlock = new(types.PairedBlock)
	seriBuf := bytes.NewBuffer(genesisBlock)
	pair.Deserialize(seriBuf)
	return &pair.Block
}

func AdamsAddress() []byte {
	address := gcrypto.GenerateKeyPair()
	return address.PubKey
}
