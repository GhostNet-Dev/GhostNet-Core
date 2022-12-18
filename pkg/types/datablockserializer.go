package types

import (
	"bytes"
	"encoding/binary"

	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
	mems "github.com/traherom/memstream"
)

func (header *GhostNetDataBlockHeader) Serialize(stream *mems.MemoryStream) {
	bs4 := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs4, uint32(header.Id))
	stream.Write(bs4)
	binary.LittleEndian.PutUint32(bs4, uint32(header.Version))
	stream.Write(bs4)
	stream.Write(header.PreviousBlockHeaderHash)
	stream.Write(header.MerkleRoot)
	binary.LittleEndian.PutUint32(bs4, uint32(header.Nonce))
	stream.Write(bs4)
	binary.LittleEndian.PutUint32(bs4, uint32(header.TransactionCount))
	stream.Write(bs4)
}

func (header *GhostNetDataBlockHeader) Deserialize(byteBuf *bytes.Buffer) {
	header.PreviousBlockHeaderHash = make([]byte, gbytes.HashSize)
	header.MerkleRoot = make([]byte, gbytes.HashSize)
	binary.Read(byteBuf, binary.LittleEndian, &header.Id)
	binary.Read(byteBuf, binary.LittleEndian, &header.Version)
	binary.Read(byteBuf, binary.LittleEndian, header.PreviousBlockHeaderHash)
	binary.Read(byteBuf, binary.LittleEndian, header.MerkleRoot)
	binary.Read(byteBuf, binary.LittleEndian, &header.Nonce)
	binary.Read(byteBuf, binary.LittleEndian, &header.TransactionCount)
}

func (dataBlock *GhostNetDataBlock) Serialize(stream *mems.MemoryStream) {
	dataBlock.Header.Serialize(stream)
	for _, tx := range dataBlock.Transaction {
		tx.Serialize(stream)
	}
}

func (dataBlock *GhostNetDataBlock) Deserialize(byteBuf *bytes.Buffer) {
	dataBlock.Header.Deserialize(byteBuf)
	dataBlock.Transaction = make([]GhostDataTransaction, dataBlock.Header.TransactionCount)
	for i := 0; i < int(dataBlock.Header.TransactionCount); i++ {
		dataBlock.Transaction[i].Deserialize(byteBuf)
	}
}
