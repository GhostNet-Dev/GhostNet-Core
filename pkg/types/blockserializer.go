package types

import (
	"bytes"
	"encoding/binary"

	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	mems "github.com/traherom/memstream"
)

var (
	binaryVersion = 0
)

type BinaryPairedBlockHeader struct {
	Version     uint32
	BlockOffset uint32
}

type BinaryGhostBlockHeader struct {
	Version               uint32
	AliceTransactionCount uint32
	TransactionCount      uint32
}

func (binaryHeader *BinaryPairedBlockHeader) Serialize(stream *mems.MemoryStream) {
	bs4 := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs4, uint32(binaryHeader.Version))
	stream.Write(bs4)
	binary.LittleEndian.PutUint32(bs4, uint32(binaryHeader.BlockOffset))
	stream.Write(bs4)
}

func (binaryHeader *BinaryPairedBlockHeader) Deserialize(byteBuf *bytes.Buffer) {
	binary.Read(byteBuf, binary.LittleEndian, &binaryHeader.Version)
	binary.Read(byteBuf, binary.LittleEndian, &binaryHeader.BlockOffset)
}

func (binaryHeader *BinaryGhostBlockHeader) Serialize(stream *mems.MemoryStream) {
	bs4 := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs4, uint32(binaryHeader.Version))
	stream.Write(bs4)
	binary.LittleEndian.PutUint32(bs4, uint32(binaryHeader.AliceTransactionCount))
	stream.Write(bs4)
	binary.LittleEndian.PutUint32(bs4, uint32(binaryHeader.TransactionCount))
	stream.Write(bs4)
}

func (binaryHeader *BinaryGhostBlockHeader) Deserialize(byteBuf *bytes.Buffer) {
	binary.Read(byteBuf, binary.LittleEndian, &binaryHeader.Version)
	binary.Read(byteBuf, binary.LittleEndian, &binaryHeader.AliceTransactionCount)
	binary.Read(byteBuf, binary.LittleEndian, &binaryHeader.TransactionCount)
}

func (header *GhostNetBlockHeader) Serialize(stream *mems.MemoryStream) {
	bs4 := make([]byte, 4)
	bs8 := make([]byte, 8)
	binary.LittleEndian.PutUint32(bs4, uint32(header.Id))
	stream.Write(bs4)
	binary.LittleEndian.PutUint32(bs4, uint32(header.Version))
	stream.Write(bs4)
	stream.Write(header.PreviousBlockHeaderHash)
	stream.Write(header.MerkleRoot)
	stream.Write(header.DataBlockHeaderHash)
	binary.LittleEndian.PutUint64(bs8, uint64(header.TimeStamp))
	stream.Write(bs8)
	binary.LittleEndian.PutUint32(bs4, uint32(header.Bits))
	stream.Write(bs4)
	binary.LittleEndian.PutUint32(bs4, uint32(header.Nonce))
	stream.Write(bs4)
	binary.LittleEndian.PutUint32(bs4, uint32(header.TransactionCount))
	stream.Write(bs4)
}

func (header *GhostNetBlockHeader) DeserializeBlockHeader(byteBuf *bytes.Buffer) {
	header.PreviousBlockHeaderHash = make([]byte, ghostBytes.HashSize)
	header.MerkleRoot = make([]byte, ghostBytes.HashSize)
	header.DataBlockHeaderHash = make([]byte, ghostBytes.HashSize)
	binary.Read(byteBuf, binary.LittleEndian, &header.Id)
	binary.Read(byteBuf, binary.LittleEndian, &header.Version)
	binary.Read(byteBuf, binary.LittleEndian, header.PreviousBlockHeaderHash)
	binary.Read(byteBuf, binary.LittleEndian, header.MerkleRoot)
	binary.Read(byteBuf, binary.LittleEndian, header.DataBlockHeaderHash)
	binary.Read(byteBuf, binary.LittleEndian, &header.TimeStamp)
	binary.Read(byteBuf, binary.LittleEndian, &header.Nonce)
	binary.Read(byteBuf, binary.LittleEndian, &header.TransactionCount)
}

func (headerEx *GhostNetBlockHeaderEx) Serialize(stream *mems.MemoryStream) {
	bs4 := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs4, uint32(headerEx.HeaderSize))
	stream.Write(bs4)
	binary.LittleEndian.PutUint32(bs4, uint32(headerEx.SignatureSize))
	stream.Write(bs4)
	headerEx.BlockSignature.Serialize(stream)
}

func (headerEx *GhostNetBlockHeaderEx) DeserializeBlockHeader(byteBuf *bytes.Buffer) {
	binary.Read(byteBuf, binary.LittleEndian, &headerEx.HeaderSize)
	binary.Read(byteBuf, binary.LittleEndian, &headerEx.SignatureSize)
	headerEx.BlockSignature.DeserializeSigHash(byteBuf)
}

func (block *GhostNetBlock) Serialize(stream *mems.MemoryStream) {
	block.Header.Serialize(stream)
	block.HeaderEx.Serialize(stream)
	block.Alice.Serialize(stream)
	for i := 0; i < int(block.Header.TransactionCount); i++ {
		block.Transaction[i].Serialize(stream)
	}
}

func (block *GhostNetBlock) DeserializeBlock(byteBuf *bytes.Buffer) {
	block.Header.DeserializeBlockHeader(byteBuf)
	block.HeaderEx.DeserializeBlockHeader(byteBuf)
	block.Alice.Deserialize(byteBuf)
	for i := 0; i < int(block.Header.TransactionCount); i++ {
		block.Transaction[i].Deserialize(byteBuf)
	}
}

func (pair *PairedBlock) Serialize(stream *mems.MemoryStream) {
	binaryPairHeader := BinaryPairedBlockHeader{}
	binaryHeader := BinaryGhostBlockHeader{}
	binaryPairHeader.Serialize(stream)
	binaryHeader.Serialize(stream)
	pair.Block.Serialize(stream)
}
