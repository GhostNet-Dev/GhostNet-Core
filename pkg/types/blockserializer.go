package types

import (
	"bytes"
	"encoding/binary"

	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	mems "github.com/traherom/memstream"
)

type BinaryPairedBlockHeader struct {
	BinaryVersion uint32
	BlockSize     uint32
}

func (binaryHeader *BinaryPairedBlockHeader) Serialize(stream *mems.MemoryStream) {
	bs4 := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs4, uint32(binaryHeader.BinaryVersion))
	stream.Write(bs4)
	binary.LittleEndian.PutUint32(bs4, uint32(binaryHeader.BlockSize))
	stream.Write(bs4)
}

func (binaryHeader *BinaryPairedBlockHeader) Deserialize(byteBuf *bytes.Buffer) {
	binary.Read(byteBuf, binary.LittleEndian, &binaryHeader.BinaryVersion)
	binary.Read(byteBuf, binary.LittleEndian, &binaryHeader.BlockSize)
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
	binary.LittleEndian.PutUint32(bs4, uint32(header.AliceCount))
	stream.Write(bs4)
	binary.LittleEndian.PutUint32(bs4, uint32(header.TransactionCount))
	stream.Write(bs4)
	binary.LittleEndian.PutUint32(bs4, uint32(header.SignatureSize))
	stream.Write(bs4)
	header.BlockSignature.Serialize(stream)
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
	binary.Read(byteBuf, binary.LittleEndian, &header.AliceCount)
	binary.Read(byteBuf, binary.LittleEndian, &header.TransactionCount)
	binary.Read(byteBuf, binary.LittleEndian, &header.SignatureSize)
	header.BlockSignature.DeserializeSigHash(byteBuf)
}

func (block *GhostNetBlock) Serialize(stream *mems.MemoryStream) {
	block.Header.Serialize(stream)
	for _, tx := range block.Alice {
		tx.Serialize(stream)
	}
	for _, tx := range block.Transaction {
		tx.Serialize(stream)
	}
}

func (block *GhostNetBlock) DeserializeBlock(byteBuf *bytes.Buffer) {
	block.Header.DeserializeBlockHeader(byteBuf)
	block.Alice = make([]GhostTransaction, block.Header.AliceCount)
	for i := 0; i < int(block.Header.AliceCount); i++ {
		block.Alice[i].Deserialize(byteBuf)
	}
	block.Transaction = make([]GhostTransaction, block.Header.TransactionCount)
	for i := 0; i < int(block.Header.TransactionCount); i++ {
		block.Transaction[i].Deserialize(byteBuf)
	}
}

func (pair *PairedBlock) Serialize(stream *mems.MemoryStream) {
	binaryPairHeader := BinaryPairedBlockHeader{
		BlockSize: pair.Block.Size(),
	}
	binaryPairHeader.Serialize(stream)
	pair.Block.Serialize(stream)
}

func (pair *PairedBlock) Deserialize(byteBuf *bytes.Buffer) {
	binaryPairHeader := BinaryPairedBlockHeader{}
	binaryPairHeader.Deserialize(byteBuf)
	pair.Block.DeserializeBlock(byteBuf)
}
