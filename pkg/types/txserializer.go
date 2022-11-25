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

func (output *TxOutput) Serialize2(byteBuf *bytes.Buffer) {
	bs4 := make([]byte, 4)
	bs8 := make([]byte, 8)

	byteBuf.Write(output.Addr[:])
	byteBuf.Write(output.BrokerAddr[:])
	binary.LittleEndian.PutUint64(bs8, uint64(output.Value))
	byteBuf.Write(bs8)
	binary.LittleEndian.PutUint32(bs4, uint32(output.ScriptSize))
	byteBuf.Write(bs4)
	byteBuf.Write(output.ScriptPubKey)
}

func (output *TxOutput) DeserializeTxOutput(byteBuf *bytes.Buffer) {
	output.Addr = make([]byte, ghostBytes.HashSize)
	output.BrokerAddr = make([]byte, ghostBytes.HashSize)
	binary.Read(byteBuf, binary.LittleEndian, output.Addr)
	binary.Read(byteBuf, binary.LittleEndian, output.BrokerAddr)
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

func (input *TxInput) DeserializeTxInput(byteBuf *bytes.Buffer) {
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
	binary.LittleEndian.PutUint32(bs, body.VinCounter)
	stream.Write(bs)
	for i := 0; i < int(body.VinCounter); i++ {
		body.Vin[i].Serialize(stream)
	}
	binary.LittleEndian.PutUint32(bs, body.VoutCounter)
	stream.Write(bs)
	for i := 0; i < int(body.VoutCounter); i++ {
		body.Vout[i].Serialize(stream)
	}
	binary.LittleEndian.PutUint32(bs, body.Nonce)
	stream.Write(bs)
	binary.LittleEndian.PutUint32(bs, body.LockTime)
	stream.Write(bs)
}

func (body *TxBody) DeserializeTxBody(byteBuf *bytes.Buffer) {
	binary.Read(byteBuf, binary.LittleEndian, &body.VinCounter)
	body.Vin = make([]TxInput, body.VinCounter)
	for i := 0; i < int(body.VinCounter); i++ {
		body.Vin[i].DeserializeTxInput(byteBuf)
	}
	binary.Read(byteBuf, binary.LittleEndian, &body.VoutCounter)
	body.Vout = make([]TxOutput, body.VoutCounter)
	for i := 0; i < int(body.VoutCounter); i++ {
		body.Vout[i].DeserializeTxOutput(byteBuf)
	}
	binary.Read(byteBuf, binary.LittleEndian, &body.Nonce)
	binary.Read(byteBuf, binary.LittleEndian, &body.LockTime)
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
	tx.Body.DeserializeTxBody(byteBuf)
}

func (dataTx *GhostDataTrasaction) Serialize(stream *mems.MemoryStream) {
	bs := make([]byte, 4)
	stream.Write(dataTx.TxId[:])
	stream.Write(dataTx.UniqHashKey[:])
	binary.LittleEndian.PutUint32(bs, dataTx.DataSize)
	if dataTx.Data != nil {
		stream.Write(dataTx.Data)
	}
}

func (dataTx *GhostDataTrasaction) Deserialize(byteBuf *bytes.Buffer) {
	dataTx.TxId = make([]byte, ghostBytes.HashSize)
	dataTx.UniqHashKey = make([]byte, ghostBytes.HashSize)
	binary.Read(byteBuf, binary.LittleEndian, dataTx.TxId)
	binary.Read(byteBuf, binary.LittleEndian, dataTx.UniqHashKey)
	binary.Read(byteBuf, binary.LittleEndian, &dataTx.DataSize)
	if byteBuf.Len() > 0 {
		dataTx.Data = make([]byte, dataTx.DataSize)
		if err := binary.Read(byteBuf, binary.LittleEndian, dataTx.Data); err != nil {
			log.Fatal("GhostDataTrasaction.Deserialize error: ", err)
		}
	}
}
