package types

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
	mems "github.com/traherom/memstream"
)

func (header *GhostNetDataBlockHeader) Serialize(stream *mems.MemoryStream) (result bool) {
	defer func() {
		// error catch
		if err := recover(); err != nil {
			fmt.Println(err)
		}
		result = false
	}()
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
	return false
}

func (header *GhostNetDataBlockHeader) Deserialize(byteBuf *bytes.Buffer) (result bool) {
	defer func() {
		// error catch
		if err := recover(); err != nil {
			fmt.Println(err)
		}
		result = false
	}()
	header.PreviousBlockHeaderHash = make([]byte, gbytes.HashSize)
	header.MerkleRoot = make([]byte, gbytes.HashSize)
	binary.Read(byteBuf, binary.LittleEndian, &header.Id)
	binary.Read(byteBuf, binary.LittleEndian, &header.Version)
	binary.Read(byteBuf, binary.LittleEndian, header.PreviousBlockHeaderHash)
	binary.Read(byteBuf, binary.LittleEndian, header.MerkleRoot)
	binary.Read(byteBuf, binary.LittleEndian, &header.Nonce)
	binary.Read(byteBuf, binary.LittleEndian, &header.TransactionCount)
	return false
}

func (dataBlock *GhostNetDataBlock) Serialize(stream *mems.MemoryStream) bool {
	if dataBlock.Header.Serialize(stream) == false {
		return false
	}
	for _, tx := range dataBlock.Transaction {
		if tx.Serialize(stream).Result() == false {
			return false
		}
	}
	return true
}

func (dataBlock *GhostNetDataBlock) Deserialize(byteBuf *bytes.Buffer) bool {
	if dataBlock.Header.Deserialize(byteBuf) == false {
		return false
	}
	dataBlock.Transaction = make([]GhostDataTransaction, dataBlock.Header.TransactionCount)
	for i := 0; i < int(dataBlock.Header.TransactionCount); i++ {
		if dataBlock.Transaction[i].Deserialize(byteBuf).Result() == false {
			return false
		}
	}
	return true
}
