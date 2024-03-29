package types

import (
	"crypto/sha256"
	"fmt"
	"unsafe"

	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
	"github.com/btcsuite/btcutil/base58"
	mems "github.com/traherom/memstream"
)

type PairedBlock struct {
	Block     GhostNetBlock
	DataBlock GhostNetDataBlock
}

type GhostNetBlock struct {
	Header      GhostNetBlockHeader `json:"Header"`
	Alice       []GhostTransaction  `json:"Alice"`
	Transaction []GhostTransaction  `json:"Transaction"`
}

type GhostNetBlockHeader struct {
	Id                      uint32           `json:"Id"`
	Version                 uint32           `json:"Version"`
	PreviousBlockHeaderHash gbytes.HashBytes `json:"PreviousBlockHeaderHash"`
	MerkleRoot              gbytes.HashBytes `json:"MerkleRoot"`
	DataBlockHeaderHash     gbytes.HashBytes `json:"DataBlockHeaderHash"`
	TimeStamp               uint64           `json:"TimeStamp"`
	Bits                    uint32           `json:"Bits"`
	Nonce                   uint32           `json:"Nonce"`
	AliceCount              uint32           `json:"AliceCount"`
	TransactionCount        uint32           `json:"TransactionCount"`
	SignatureSize           uint32           `json:"SignatureSize"`
	BlockSignature          SigHash          `json:"BlockSignature"`
}

func (pairedBlock *PairedBlock) BlockId() uint32 {
	return pairedBlock.Block.Header.Id
}

func (pairedBlock *PairedBlock) TxCount() uint32 {
	return pairedBlock.Block.Header.TransactionCount
}

func (pairedBlock *PairedBlock) GetBlockFilename() string {
	return fmt.Sprint(pairedBlock.BlockId(), "@", base58.Encode(pairedBlock.Block.GetHashKey()), ".ghost")
}

func (header *GhostNetBlockHeader) Size() uint32 {
	return uint32(unsafe.Sizeof(header.Id)) +
		uint32(unsafe.Sizeof(header.Version)) + gbytes.HashSize*3 +
		uint32(unsafe.Sizeof(header.TimeStamp)) +
		uint32(unsafe.Sizeof(header.Bits)) +
		uint32(unsafe.Sizeof(header.Nonce)) +
		uint32(unsafe.Sizeof(header.AliceCount)) +
		uint32(unsafe.Sizeof(header.TransactionCount)) +
		uint32(unsafe.Sizeof(header.SignatureSize)) +
		uint32(header.BlockSignature.Size())
}

func (block *GhostNetBlock) Size() uint32 {
	var txSize uint32 = 0
	for _, tx := range block.Alice {
		txSize += tx.Size()
	}
	for _, tx := range block.Transaction {
		txSize += tx.Size()
	}
	return block.Header.Size() + txSize
}

func (header *GhostNetBlockHeader) GetHashKey() []byte {
	size := header.Size()
	stream := mems.NewCapacity(int(size))
	hash := sha256.New()
	header.Serialize(stream)
	hash.Write(stream.Bytes())
	return hash.Sum(nil)
}

func (block *GhostNetBlock) GetHashKey() []byte {
	return block.Header.GetHashKey()
}

func (pair *PairedBlock) Size() uint32 {
	return pair.Block.Size() + pair.DataBlock.Size()
}
