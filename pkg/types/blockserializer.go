package types

import (
	"bytes"
	"encoding/binary"

	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
	mems "github.com/traherom/memstream"
)

func (header *GhostNetBlockHeader) SerializeToByte() []byte {
	size := header.Size()
	stream := mems.NewCapacity(int(size))
	header.Serialize(stream)
	return stream.Bytes()
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
	stream.Write(header.BlockSignature.SerializeToByte())
}

func (header *GhostNetBlockHeader) Deserialize(byteBuf *bytes.Buffer) {
	header.PreviousBlockHeaderHash = make([]byte, gbytes.HashSize)
	header.MerkleRoot = make([]byte, gbytes.HashSize)
	header.DataBlockHeaderHash = make([]byte, gbytes.HashSize)
	binary.Read(byteBuf, binary.LittleEndian, &header.Id)
	binary.Read(byteBuf, binary.LittleEndian, &header.Version)
	binary.Read(byteBuf, binary.LittleEndian, header.PreviousBlockHeaderHash)
	binary.Read(byteBuf, binary.LittleEndian, header.MerkleRoot)
	binary.Read(byteBuf, binary.LittleEndian, header.DataBlockHeaderHash)
	binary.Read(byteBuf, binary.LittleEndian, &header.TimeStamp)
	binary.Read(byteBuf, binary.LittleEndian, &header.Bits)
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

func (block *GhostNetBlock) SerializeToByte() []byte {
	size := block.Size()
	stream := mems.NewCapacity(int(size))
	block.Serialize(stream)
	return stream.Bytes()
}

func (block *GhostNetBlock) Deserialize(byteBuf *bytes.Buffer) {
	block.Header.Deserialize(byteBuf)
	block.Alice = make([]GhostTransaction, block.Header.AliceCount)
	for i := 0; i < int(block.Header.AliceCount); i++ {
		block.Alice[i].Deserialize(byteBuf)
	}
	block.Transaction = make([]GhostTransaction, block.Header.TransactionCount)
	for i := 0; i < int(block.Header.TransactionCount); i++ {
		block.Transaction[i].Deserialize(byteBuf)
	}
}

func (pair *PairedBlock) SerializeToByte() []byte {
	size := pair.Size()
	stream := mems.NewCapacity(int(size))
	pair.Serialize(stream)
	return stream.Bytes()
}

func (pair *PairedBlock) Serialize(stream *mems.MemoryStream) {
	pair.Block.Serialize(stream)
	pair.DataBlock.Serialize(stream)
}

func (pair *PairedBlock) Deserialize(byteBuf *bytes.Buffer) {
	pair.Block.Deserialize(byteBuf)
	pair.DataBlock.Deserialize(byteBuf)
}
