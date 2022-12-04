package types

import (
	"bytes"
	"encoding/binary"
	"log"

	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	mems "github.com/traherom/memstream"
)

func (output *TxOutput) Serialize(stream *mems.MemoryStream) {
	bs4 := make([]byte, 4)
	bs8 := make([]byte, 8)

	stream.Write(output.Addr[:])
	stream.Write(output.BrokerAddr[:])
	binary.LittleEndian.PutUint64(bs8, uint64(output.Value))
	stream.Write(bs8)
	binary.LittleEndian.PutUint32(bs4, uint32(output.ScriptSize))
	stream.Write(bs4)
	stream.Write(output.ScriptPubKey)
}

func (output *TxOutput) Deserialize(byteBuf *bytes.Buffer) {
	output.Addr = make([]byte, ghostBytes.HashSize)
	output.BrokerAddr = make([]byte, ghostBytes.HashSize)
	binary.Read(byteBuf, binary.LittleEndian, output.Addr)
	binary.Read(byteBuf, binary.LittleEndian, output.BrokerAddr)
	binary.Read(byteBuf, binary.LittleEndian, &output.Type)
	binary.Read(byteBuf, binary.LittleEndian, &output.Value)
	binary.Read(byteBuf, binary.LittleEndian, &output.ScriptSize)
	output.ScriptPubKey = make([]byte, output.ScriptSize)
	binary.Read(byteBuf, binary.LittleEndian, output.ScriptPubKey)
}

func (input *TxInput) Serialize(stream *mems.MemoryStream) {
	bs := make([]byte, 4)

	stream.Write(input.PrevOut.TxId[:])
	binary.LittleEndian.PutUint32(bs, input.PrevOut.TxOutIndex)
	stream.Write(bs)
	binary.LittleEndian.PutUint32(bs, input.Sequence)
	stream.Write(bs)
	binary.LittleEndian.PutUint32(bs, input.ScriptSize)
	stream.Write(bs)
	stream.Write(input.ScriptSig)
}

func (input *TxInput) Deserialize(byteBuf *bytes.Buffer) {
	//byteBuf := bytes.NewBuffer(buf)
	input.PrevOut.TxId = make([]byte, ghostBytes.HashSize)
	binary.Read(byteBuf, binary.LittleEndian, input.PrevOut.TxId)
	binary.Read(byteBuf, binary.LittleEndian, &input.PrevOut.TxOutIndex)
	binary.Read(byteBuf, binary.LittleEndian, &input.Sequence)
	binary.Read(byteBuf, binary.LittleEndian, &input.ScriptSize)
	input.ScriptSig = make([]byte, input.ScriptSize)
	binary.Read(byteBuf, binary.LittleEndian, input.ScriptSig)
}

func (body *TxBody) Serialize(stream *mems.MemoryStream) {
	bs := make([]byte, 4)
	binary.LittleEndian.PutUint32(bs, body.InputCounter)
	stream.Write(bs)
	for i := 0; i < int(body.InputCounter); i++ {
		body.Vin[i].Serialize(stream)
	}
	binary.LittleEndian.PutUint32(bs, body.OutputCounter)
	stream.Write(bs)
	for i := 0; i < int(body.OutputCounter); i++ {
		body.Vout[i].Serialize(stream)
	}
	binary.LittleEndian.PutUint32(bs, body.Nonce)
	stream.Write(bs)
	binary.LittleEndian.PutUint32(bs, body.LockTime)
	stream.Write(bs)
}

func (body *TxBody) Deserialize(byteBuf *bytes.Buffer) {
	binary.Read(byteBuf, binary.LittleEndian, &body.InputCounter)
	body.Vin = make([]TxInput, body.InputCounter)
	for i := 0; i < int(body.InputCounter); i++ {
		body.Vin[i].Deserialize(byteBuf)
	}
	binary.Read(byteBuf, binary.LittleEndian, &body.OutputCounter)
	body.Vout = make([]TxOutput, body.OutputCounter)
	for i := 0; i < int(body.OutputCounter); i++ {
		body.Vout[i].Deserialize(byteBuf)
	}
	binary.Read(byteBuf, binary.LittleEndian, &body.Nonce)
	binary.Read(byteBuf, binary.LittleEndian, &body.LockTime)
}

func (tx *GhostTransaction) SerializeToByte() []byte {
	size := tx.Size()
	stream := mems.NewCapacity(int(size))
	tx.Serialize(stream)
	return stream.Bytes()
}

func (tx *GhostTransaction) Serialize(stream *mems.MemoryStream) {
	/*
		size := tx.Size()
		stream := mems.NewCapacity(int(size))*/
	stream.Write(tx.TxId[:])
	tx.Body.Serialize(stream)
}

func (tx *GhostTransaction) Deserialize(byteBuf *bytes.Buffer) {
	tx.TxId = make([]byte, ghostBytes.HashSize)
	binary.Read(byteBuf, binary.LittleEndian, tx.TxId)
	tx.Body.Deserialize(byteBuf)
}

func (tx *GhostDataTransaction) SerializeToByte() []byte {
	size := tx.Size()
	stream := mems.NewCapacity(int(size))
	tx.Serialize(stream)
	return stream.Bytes()
}

func (tx *GhostDataTransaction) Serialize(stream *mems.MemoryStream) {
	bs4 := make([]byte, 4)
	bs8 := make([]byte, 8)
	stream.Write(tx.TxId[:])
	binary.LittleEndian.PutUint64(bs8, tx.LogicalAddress)
	stream.Write(bs8)
	binary.LittleEndian.PutUint32(bs4, tx.DataSize)
	stream.Write(bs4)
	if tx.Data != nil {
		stream.Write(tx.Data)
	}
}

func (tx *GhostDataTransaction) Deserialize(byteBuf *bytes.Buffer) {
	tx.TxId = make([]byte, ghostBytes.HashSize)
	binary.Read(byteBuf, binary.LittleEndian, tx.TxId)
	binary.Read(byteBuf, binary.LittleEndian, &tx.LogicalAddress)
	binary.Read(byteBuf, binary.LittleEndian, &tx.DataSize)
	if byteBuf.Len() > 0 {
		tx.Data = make([]byte, tx.DataSize)
		if err := binary.Read(byteBuf, binary.LittleEndian, tx.Data); err != nil {
			log.Fatal("GhostDataTrasaction.Deserialize error: ", err)
		}
	}
}
