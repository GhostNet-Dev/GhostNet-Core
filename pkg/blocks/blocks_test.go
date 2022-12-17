package blocks

import (
	"bytes"
	"crypto/sha256"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/txs"

	"github.com/stretchr/testify/assert"
)

var (
	Sender = gcrypto.GenerateKeyPair()
	Broker = gcrypto.GenerateKeyPair()
	Recver = gcrypto.GenerateKeyPair()

	gScript        = gvm.NewGScript()
	gVm            = gvm.NewGVM()
	BlockContainer = store.NewBlockContainer()
	Txs            = txs.NewTXs(gScript, BlockContainer, gVm)
	blocks         = NewBlocks(BlockContainer, Txs, 1)
)

func TestNewBlocks(t *testing.T) {
	assert.Equal(t, nil, blocks, "")
}

func TestMerkleTree(t *testing.T) {
	testBuf := [][]byte{{0x11, 0x22, 0x33},
		{0x22, 0x33, 0x44}, {0x33, 0x44, 0x55}}
	hashs := make([][]byte, 3)
	hash := sha256.New()
	for i := range hashs {
		hash.Write(testBuf[i])
		hashs[i] = hash.Sum(nil)
		hash.Reset()
	}
	result := CreateMerkleRoot(hashs)
	result2 := CreateMerkleRoot(hashs)
	compareResult := bytes.Compare(result, result2)
	assert.Equal(t, 0, compareResult, "The merkle root is not the same.")

	testBuf[0][0] = 0x66
	hash.Write(testBuf[0])
	hashs[0] = hash.Sum(nil)
	hash.Reset()
	result2 = CreateMerkleRoot(hashs)
	compareResult = bytes.Compare(result, result2)
	assert.Equal(t, true, compareResult != 0, "The merkle root has not changed.")
}
