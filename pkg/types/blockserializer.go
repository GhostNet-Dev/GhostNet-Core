package types

import (
	"bytes"
	"encoding/binary"
	"fmt"

	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
	mems "github.com/traherom/memstream"
)

func (header *GhostNetBlockHeader) SerializeToByte() []byte {
	size := header.Size()
	stream := mems.NewCapacity(int(size))
	if !header.Serialize(stream) {
		return nil
	}
	return stream.Bytes()
}

func (header *GhostNetBlockHeader) Serialize(stream *mems.MemoryStream) (result bool) {
	defer func() {
		// error catch
		if err := recover(); err != nil {
			fmt.Println(err)
			result = false
		} else {
			result = true
		}
	}()
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
	return result
}

func (header *GhostNetBlockHeader) Deserialize(byteBuf *bytes.Buffer) (result bool) {
	defer func() {
		// error catch
		if err := recover(); err != nil {
			fmt.Println(err)
			result = false
		} else {
			result = true
		}
	}()

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
	return result
}

func (block *GhostNetBlock) Serialize(stream *mems.MemoryStream) (result bool) {
	if result = block.Header.Serialize(stream); !result {
		return false
	}
	for _, tx := range block.Alice {
		if !tx.Serialize(stream).Result() {
			return false
		}
	}
	for _, tx := range block.Transaction {
		if !tx.Serialize(stream).Result() {
			return false
		}
	}
	return true
}

func (block *GhostNetBlock) SerializeToByte() []byte {
	size := block.Size()
	stream := mems.NewCapacity(int(size))
	if !block.Serialize(stream) {
		return nil
	}
	return stream.Bytes()
}

func (block *GhostNetBlock) Deserialize(byteBuf *bytes.Buffer) (result bool) {
	if result = block.Header.Deserialize(byteBuf); !result {
		return false
	}
	block.Alice = make([]GhostTransaction, block.Header.AliceCount)
	for i := 0; i < int(block.Header.AliceCount); i++ {
		if !block.Alice[i].Deserialize(byteBuf).Result() {
			return false
		}
	}
	block.Transaction = make([]GhostTransaction, block.Header.TransactionCount)
	for i := 0; i < int(block.Header.TransactionCount); i++ {
		if !block.Transaction[i].Deserialize(byteBuf).Result() {
			return false
		}
	}
	return true
}

func (pair *PairedBlock) SerializeToByte() []byte {
	size := pair.Size()
	stream := mems.NewCapacity(int(size))
	if !pair.Serialize(stream) {
		return nil
	}
	return stream.Bytes()
}

func (pair *PairedBlock) Serialize(stream *mems.MemoryStream) bool {
	if !pair.Block.Serialize(stream) || !pair.DataBlock.Serialize(stream) {
		return false
	}
	return true
}

func (pair *PairedBlock) Deserialize(byteBuf *bytes.Buffer) bool {
	if !pair.Block.Deserialize(byteBuf) || !pair.DataBlock.Deserialize(byteBuf) {
		return false
	}
	return true
}
