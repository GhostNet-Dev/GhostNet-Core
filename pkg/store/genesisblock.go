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
)

func GenesisBlock() *types.PairedBlock {
	if genesisPairedBlock == nil {
		genesisPairedBlock = new(types.PairedBlock)
		seriBuf := bytes.NewBuffer(genesisBlockByte)
		genesisPairedBlock.Deserialize(seriBuf)
	}
	return genesisPairedBlock
}

func AdamsAddress() *gcrypto.GhostAddress {
	return gcrypto.GenerateKeyPair()
}
