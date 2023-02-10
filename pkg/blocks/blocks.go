package blocks

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/txs"
)

type Blocks struct {
	txs            *txs.TXs
	blockContainer *store.BlockContainer
	gScript        *gvm.GScript
	Version        uint32
	miningFlag     bool
}

func NewBlocks(b *store.BlockContainer, t *txs.TXs, version uint32) *Blocks {
	return &Blocks{
		txs:            t,
		blockContainer: b,
		Version:        version,
		miningFlag:     false,
	}
}

func (blocks *Blocks) ResetBlocks() {
	blocks.blockContainer.Close()
	blocks.blockContainer.Reset()
	blocks.blockContainer.Reopen()
}
