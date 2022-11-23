package store

import (
	"bytes"
	_ "embed"

	types "github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

//go:embed GenesisBlocks
var genesisBlock []byte

func GenesisBlock() *types.GhostNetBlock {
	var block *types.GhostNetBlock = new(types.GhostNetBlock)
	seriBuf := bytes.NewBuffer(genesisBlock)
	block.DeserializeBlock(seriBuf)
	return block
}
