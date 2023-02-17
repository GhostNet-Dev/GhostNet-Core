package types

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
	mems "github.com/traherom/memstream"
)

const (
	SizeError = iota
	Success
	Exception
)

type SerialiazeError struct {
	result int
}

func (err *SerialiazeError) Result() bool {
	if err == nil {
		return true
	}
	return err.result == Success
}

func (err *SerialiazeError) Error() string {
	if err == nil {
		return "success"
	}
	resultString := "error"
	switch err.result {
	case Success:
		resultString = "success"
	case Exception:
		resultString = "Exception"
	}
	return resultString
}

func (output *TxOutput) Serialize(stream *mems.MemoryStream) (sErr *SerialiazeError) {
	defer func() {
		// error catch
		if err := recover(); err != nil {
			fmt.Println(err)
			sErr = &SerialiazeError{result: Exception}
		} else {
			sErr = nil
		}
	}()
	bs4 := make([]byte, 4)
	bs8 := make([]byte, 8)
	if len(output.Addr) != int(gbytes.PubKeySize) {
		return &SerialiazeError{result: SizeError}
	}

	stream.Write(output.Addr[:])
	stream.Write(output.BrokerAddr[:])
	binary.LittleEndian.PutUint32(bs4, uint32(output.Type))
	stream.Write(bs4)
	binary.LittleEndian.PutUint64(bs8, uint64(output.Value))
	stream.Write(bs8)
	binary.LittleEndian.PutUint32(bs4, uint32(output.ScriptSize))
	stream.Write(bs4)
	stream.Write(output.ScriptPubKey)

	return sErr
}

func (output *TxOutput) Deserialize(byteBuf *bytes.Buffer) (sErr *SerialiazeError) {
	defer func() {
		// error catch
		if err := recover(); err != nil {
			fmt.Println(err)
			sErr = &SerialiazeError{result: Exception}
		} else {
			sErr = nil
		}
	}()
	output.Addr = make([]byte, gbytes.PubKeySize)
	output.BrokerAddr = make([]byte, gbytes.PubKeySize)
	binary.Read(byteBuf, binary.LittleEndian, output.Addr)
	binary.Read(byteBuf, binary.LittleEndian, output.BrokerAddr)
	binary.Read(byteBuf, binary.LittleEndian, &output.Type)
	binary.Read(byteBuf, binary.LittleEndian, &output.Value)
	binary.Read(byteBuf, binary.LittleEndian, &output.ScriptSize)
	output.ScriptPubKey = make([]byte, output.ScriptSize)
	binary.Read(byteBuf, binary.LittleEndian, output.ScriptPubKey)
	return sErr
}

func (input *TxInput) Serialize(stream *mems.MemoryStream) (sErr *SerialiazeError) {
	defer func() {
		// error catch
		if err := recover(); err != nil {
			fmt.Println(err)
			sErr = &SerialiazeError{result: Exception}
		} else {
			sErr = nil
		}
	}()
	bs := make([]byte, 4)

	stream.Write(input.PrevOut.TxId[:])
	binary.LittleEndian.PutUint32(bs, input.PrevOut.TxOutIndex)
	stream.Write(bs)
	binary.LittleEndian.PutUint32(bs, input.Sequence)
	stream.Write(bs)
	binary.LittleEndian.PutUint32(bs, input.ScriptSize)
	stream.Write(bs)
	stream.Write(input.ScriptSig)
	return sErr
}

func (input *TxInput) Deserialize(byteBuf *bytes.Buffer) (sErr *SerialiazeError) {
	defer func() {
		// error catch
		if err := recover(); err != nil {
			fmt.Println(err)
			sErr = &SerialiazeError{result: Exception}
		} else {
			sErr = nil
		}
	}()
	input.PrevOut.TxId = make([]byte, gbytes.HashSize)
	binary.Read(byteBuf, binary.LittleEndian, input.PrevOut.TxId)
	binary.Read(byteBuf, binary.LittleEndian, &input.PrevOut.TxOutIndex)
	binary.Read(byteBuf, binary.LittleEndian, &input.Sequence)
	binary.Read(byteBuf, binary.LittleEndian, &input.ScriptSize)
	input.ScriptSig = make([]byte, input.ScriptSize)
	binary.Read(byteBuf, binary.LittleEndian, input.ScriptSig)

	return sErr
}

func (body *TxBody) Serialize(stream *mems.MemoryStream) (sErr *SerialiazeError) {
	defer func() {
		// error catch
		if err := recover(); err != nil {
			fmt.Println(err)
			sErr = &SerialiazeError{result: Exception}
		} else {
			sErr = nil
		}
	}()
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

	return sErr
}

func (body *TxBody) Deserialize(byteBuf *bytes.Buffer) (sErr *SerialiazeError) {
	defer func() {
		// error catch
		if err := recover(); err != nil {
			fmt.Println(err)
			sErr = &SerialiazeError{result: Exception}
		} else {
			sErr = nil
		}
	}()
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

	return sErr
}

func (tx *GhostTransaction) SerializeToByte() []byte {
	size := tx.Size()
	stream := mems.NewCapacity(int(size))
	if !tx.Serialize(stream).Result() {
		return nil
	}
	return stream.Bytes()
}

func (tx *GhostTransaction) Serialize(stream *mems.MemoryStream) (sErr *SerialiazeError) {
	defer func() {
		// error catch
		if err := recover(); err != nil {
			fmt.Println(err)
			sErr = &SerialiazeError{result: Exception}
		} else {
			sErr = nil
		}
	}()
	/*
		size := tx.Size()
		stream := mems.NewCapacity(int(size))*/
	stream.Write(tx.TxId[:])
	tx.Body.Serialize(stream)
	return sErr
}

func (tx *GhostTransaction) Deserialize(byteBuf *bytes.Buffer) (sErr *SerialiazeError) {
	defer func() {
		// error catch
		if err := recover(); err != nil {
			fmt.Println(err)
			sErr = &SerialiazeError{result: Exception}
		} else {
			sErr = nil
		}
	}()
	tx.TxId = make([]byte, gbytes.HashSize)
	binary.Read(byteBuf, binary.LittleEndian, tx.TxId)
	tx.Body.Deserialize(byteBuf)
	return sErr
}

func (tx *GhostDataTransaction) SerializeToByte() []byte {
	size := tx.Size()
	stream := mems.NewCapacity(int(size))
	if !tx.Serialize(stream).Result() {
		return nil
	}
	return stream.Bytes()
}

func (tx *GhostDataTransaction) Serialize(stream *mems.MemoryStream) (sErr *SerialiazeError) {
	defer func() {
		// error catch
		if err := recover(); err != nil {
			fmt.Println(err)
			sErr = &SerialiazeError{result: Exception}
		} else {
			sErr = nil
		}

	}()
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
	return sErr
}

func (tx *GhostDataTransaction) Deserialize(byteBuf *bytes.Buffer) (sErr *SerialiazeError) {
	defer func() {
		// error catch
		if err := recover(); err != nil {
			fmt.Println(err)
			sErr = &SerialiazeError{result: Exception}
		} else {
			sErr = nil
		}
	}()
	tx.TxId = make([]byte, gbytes.HashSize)
	binary.Read(byteBuf, binary.LittleEndian, tx.TxId)
	binary.Read(byteBuf, binary.LittleEndian, &tx.LogicalAddress)
	binary.Read(byteBuf, binary.LittleEndian, &tx.DataSize)
	if byteBuf.Len() > 0 {
		tx.Data = make([]byte, tx.DataSize)
		if err := binary.Read(byteBuf, binary.LittleEndian, tx.Data); err != nil {
			log.Fatal("GhostDataTrasaction.Deserialize error: ", err)
		}
	}
	return sErr
}
