package blocks

import (
	"bytes"
	"crypto/sha256"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/txs"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"

	"github.com/stretchr/testify/assert"
)

var (
	Miner  = gcrypto.GenerateKeyPair()
	Sender = gcrypto.GenerateKeyPair()
	Broker = gcrypto.GenerateKeyPair()
	Recver = gcrypto.GenerateKeyPair()

	gScript        = gvm.NewGScript()
	gVm            = gvm.NewGVM()
	BlockContainer = store.NewBlockContainer()
	Txs            = txs.NewTXs(gScript, BlockContainer, gVm)
	blocks         = NewBlocks(BlockContainer, Txs, 1)
)

func init() {
	BlockContainer.BlockContainerOpen("../../db.sqlite3.sql", "./")
}

func TestNewBlocks(t *testing.T) {
	pair := MakeNewPair()
	result := blocks.BlockValidation(pair, nil)
	assert.Equal(t, true, result, "block is not valid")
}

func TestBlockSignature(t *testing.T) {
	block := MakeNewPair().Block
	// verify
	sig := block.Header.BlockSignature

	block.Header.BlockSignature = types.SigHash{}
	block.Header.SignatureSize = uint32(block.Header.BlockSignature.Size())
	sigPack := gcrypto.SignaturePackage{
		PubKey:    sig.PubKey,
		Signature: append(sig.RBuf, sig.SBuf...),
		Text:      block.Header.SerializeToByte(),
	}
	result := gcrypto.SignVerify(&sigPack)
	assert.Equal(t, true, result, "block signature is not valid")
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

func MakeNewPair() *types.PairedBlock {
	tx, _ := Txs.MakeSampleRootAccount("test", Broker.Get160PubKey())
	txs := []types.GhostTransaction{*tx}
	msg := make([]byte, gbytes.HashSize)
	msg2 := make([]byte, gbytes.HashSize)
	copy(msg, []byte("test"))
	copy(msg2, []byte("test is important"))
	return &types.PairedBlock{
		*blocks.CreateGhostNetBlock(1, msg, msg2, Miner, Broker.Get160PubKey(), txs),
		*blocks.CreateGhostNetDataBlock(1, msg, nil),
	}
}
