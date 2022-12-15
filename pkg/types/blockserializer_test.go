package types

import (
	"bytes"
	"crypto/sha256"
	"testing"
	"time"

	gbytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	"github.com/stretchr/testify/assert"
	mems "github.com/traherom/memstream"
)

func TestBlockSerialize(t *testing.T) {
	txs := []GhostTransaction{MakeTx(), MakeTx()}
	block := &GhostNetBlock{
		Header:      *MakeHeader(uint32(len(txs)), uint32(len(txs))),
		Alice:       txs,
		Transaction: txs,
	}
	blockByte := block.SerializeToByte()

	newBlock := GhostNetBlock{}
	byteBuf := bytes.NewBuffer(blockByte)
	newBlock.Deserialize(byteBuf)
	newBlockByte := newBlock.SerializeToByte()

	result := bytes.Compare(blockByte, newBlockByte)

	assert.Equal(t, 0, result, "binary is different.")
}

func TestHeaderSerialize(t *testing.T) {
	header := MakeHeader(2, 2)
	size := header.Size()
	stream := mems.NewCapacity(int(size))
	header.Serialize(stream)

	oriBuf := stream.Bytes()
	byteBuf := bytes.NewBuffer(oriBuf)
	newHeader := GhostNetBlockHeader{}
	newHeader.Deserialize(byteBuf)

	size = newHeader.Size()
	stream = mems.NewCapacity(int(size))
	newHeader.Serialize(stream)
	newByteBuf := stream.Bytes()
	result := bytes.Compare(oriBuf, newByteBuf)
	assert.Equal(t, 0, result, "binary is different.")
}

func MakeHeader(aliceCount uint32, txsCount uint32) *GhostNetBlockHeader {
	hash := sha256.New()
	hash.Write(make([]byte, gbytes.HashSize))
	hashByte := hash.Sum((nil))
	sigHash := SigHash{}
	return &GhostNetBlockHeader{
		Id:                      1,
		Version:                 2,
		PreviousBlockHeaderHash: hashByte,
		MerkleRoot:              hashByte,
		DataBlockHeaderHash:     hashByte,
		AliceCount:              aliceCount,
		TransactionCount:        txsCount,
		TimeStamp:               uint64(time.Now().Unix()),
		BlockSignature:          sigHash,
		SignatureSize:           uint32(sigHash.Size()),
	}
}
