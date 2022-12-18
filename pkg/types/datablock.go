package types

import (
	"crypto/sha256"
	"unsafe"

	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
	mems "github.com/traherom/memstream"
)

type GhostNetDataBlock struct {
	Header      GhostNetDataBlockHeader
	Transaction []GhostDataTransaction
}

type GhostNetDataBlockHeader struct {
	Id                      uint32
	Version                 uint32
	PreviousBlockHeaderHash gbytes.HashBytes
	MerkleRoot              gbytes.HashBytes
	Nonce                   uint32
	TransactionCount        uint32
}

func (header *GhostNetDataBlockHeader) Size() uint32 {
	return uint32(unsafe.Sizeof(header.Id)) +
		uint32(unsafe.Sizeof(header.Version)) + gbytes.HashSize*2 +
		uint32(unsafe.Sizeof(header.Nonce)) +
		uint32(unsafe.Sizeof(header.TransactionCount))
}

func (dataBlock *GhostNetDataBlock) Size() uint32 {
	var txSize uint32 = 0
	for _, tx := range dataBlock.Transaction {
		txSize += tx.Size()
	}
	return dataBlock.Header.Size() + txSize
}

func (block *GhostNetDataBlock) GetHashKey() []byte {
	size := block.Header.Size()
	stream := mems.NewCapacity(int(size))
	hash := sha256.New()
	block.Header.Serialize(stream)
	hash.Write(stream.Bytes())
	return hash.Sum(nil)
}
