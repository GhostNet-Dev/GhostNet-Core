package blocks

import (
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"

	"github.com/stretchr/testify/assert"
)

var (
	Sender = gcrypto.GenerateKeyPair()
	Broker = gcrypto.GenerateKeyPair()
	Recver = gcrypto.GenerateKeyPair()

	gScript        = gvm.NewGScript()
	gVm            = gvm.NewGVM()
	blockContainer = store.NewBlockContainer()
)

func TestNewBlocks(t *testing.T) {
	blocks := NewBlocks(blockContainer)

	assert.Equal(t, nil, blocks, "")
}
