package types

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBetweenBlockandProto(t *testing.T) {
	txs := []GhostTransaction{MakeTx(), MakeTx()}
	block := &GhostNetBlock{
		Header:      *MakeHeader(uint32(len(txs)), uint32(len(txs))),
		Alice:       txs,
		Transaction: txs,
	}
	pair := &PairedBlock{
		Block:     *block,
		DataBlock: GhostNetDataBlock{},
	}
	protoBlock := GhostBlockToProtoType(pair)
	newPair := ProtoTypeToGhostBlock(protoBlock)
	compOri := block.SerializeToByte()
	compCpy := newPair.Block.SerializeToByte()
	ret := bytes.Compare(compOri, compCpy)
	assert.Equal(t, 0, ret, "different between protoblock and block")
}
