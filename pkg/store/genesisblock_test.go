package store

import (
	"bytes"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"github.com/stretchr/testify/assert"
)

var (
	gScript = gvm.NewGCompiler()
)

func TestGenesisBlockLoad(t *testing.T) {
	pair := GenesisBlock()
	block := pair.Block
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

func TestAdamLoad(t *testing.T) {
	pubKey := AdamsAddress()
	assert.Equal(t, true, pubKey != nil, "could not find adams address")
}

func TestBetweenDbAndFile(t *testing.T) {
	blockContainer := NewBlockContainer("sqlite3")
	blockContainer.BlockContainerOpen("./")
	gene := GenesisBlock()
	blockContainer.InsertBlock(gene)
	pairedBlock := blockContainer.gSql.SelectBlock(1)

	fromDb := pairedBlock.Block.GetHashKey()
	fromFile := gene.Block.GetHashKey()
	result := bytes.Equal(fromDb, fromFile)
	assert.Equal(t, true, result, "genesis block is not equal between db and file")

}
