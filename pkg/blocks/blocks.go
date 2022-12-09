package blocks

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
)

type Blocks struct {
	blockContainer *store.BlockContainer
	Version        uint32
}

func NewBlocks(b *store.BlockContainer) *Blocks {
	return &Blocks{
		blockContainer: b,
		Version:        1,
	}
}
