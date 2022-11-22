package types

import (
	"bytes"
	"encoding/binary"

	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	mems "github.com/traherom/memstream"
)

func (output *TxOutput) Serialize() []byte {
	size := output.Size()
	stream := mems.NewCapacity(int(size))
	bs := make([]byte, 8)

	stream.Write(output.Addr[:])
	stream.Write(output.BrokerAddr[:])
	binary.LittleEndian.PutUint64(bs, uint64(output.Value))
	stream.Write(bs)
	binary.LittleEndian.PutUint32(bs, uint32(output.ScriptSize))
	stream.Write(bs)
	stream.Write(output.ScriptPubKey)
	return stream.Bytes()
}

func (output *TxOutput) DeserializeTxOutput(buf []byte) {
	//txOutput.Addr = buf[:offset+len(txOutput.Addr)]
	byteBuf := bytes.NewBuffer(buf)
	output.Addr = make([]byte, ghostBytes.HashSize)
	output.BrokerAddr = make([]byte, ghostBytes.HashSize)
	binary.Read(byteBuf, binary.LittleEndian, output.Addr)
	binary.Read(byteBuf, binary.LittleEndian, output.BrokerAddr)
	binary.Read(byteBuf, binary.LittleEndian, &output.Value)
	binary.Read(byteBuf, binary.LittleEndian, &output.ScriptSize)
	output.ScriptPubKey = make([]byte, output.ScriptSize)
	binary.Read(byteBuf, binary.LittleEndian, output.ScriptPubKey)
	//binary.LittleEndian.PutUint32(buf[offset:], uint32(txOutput.Addr[]))
}

func (input *TxInput) Serialize() []byte {
	size := input.Size()
	stream := mems.NewCapacity(int(size))
	bs := make([]byte, 4)

	stream.Write(input.PrevOut.TxId[:])
	binary.LittleEndian.PutUint32(bs, input.PrevOut.TxOutIndex)
	stream.Write(bs)
	binary.LittleEndian.PutUint32(bs, input.Sequence)
	stream.Write(bs)
	binary.LittleEndian.PutUint32(bs, input.ScriptSize)
	stream.Write(bs)
	stream.Write(input.ScriptSig)
	return stream.Bytes()
}

func (input *TxInput) DeserializeTxInput(buf []byte) {
	byteBuf := bytes.NewBuffer(buf)
	input.PrevOut.TxId = make([]byte, ghostBytes.HashSize)
	binary.Read(byteBuf, binary.LittleEndian, input.PrevOut.TxId)
	binary.Read(byteBuf, binary.LittleEndian, &input.PrevOut.TxOutIndex)
	binary.Read(byteBuf, binary.LittleEndian, &input.Sequence)
	binary.Read(byteBuf, binary.LittleEndian, &input.ScriptSize)
	input.ScriptSig = make([]byte, input.ScriptSize)
	binary.Read(byteBuf, binary.LittleEndian, input.ScriptSig)
}

func (body *TxBody) Serialize() []byte {
	size := body.Size()
	stream := mems.NewCapacity(int(size))
	bs := make([]byte, 4)

	binary.LittleEndian.PutUint32(bs, body.VinCounter)
	stream.Write(bs)
	for i := 0; i < int(body.VinCounter); i++ {
		stream.Write(body.Vin[i].Serialize())
	}
	binary.LittleEndian.PutUint32(bs, body.VoutCounter)
	stream.Write(bs)
	for i := 0; i < int(body.VoutCounter); i++ {
		stream.Write(body.Vout[i].Serialize())
	}
	binary.LittleEndian.PutUint32(bs, body.Nonce)
	stream.Write(bs)
	binary.LittleEndian.PutUint32(bs, body.LockTime)
	stream.Write(bs)
	return stream.Bytes()
}

func (tx *GhostTransaction) Serialize() []byte {
	size := tx.Size()
	stream := mems.NewCapacity(int(size))
	stream.Write(tx.TxId[:])
	stream.Write((tx.Body.Serialize()))
	return stream.Bytes()
}
