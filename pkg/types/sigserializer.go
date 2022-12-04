package types

import (
	"bytes"
	"encoding/binary"

	mems "github.com/traherom/memstream"
)

func (sig *SigHash) SerializeToByte() []byte {
	size := sig.Size()
	stream := mems.NewCapacity(int(size))
	sig.Serialize(stream)
	return stream.Bytes()
}

func (sig *SigHash) Serialize(stream *mems.MemoryStream) {
	bs := []byte{sig.DER, sig.SignatureSize, sig.RType, sig.RSize}
	bs4 := make([]byte, 4)
	stream.Write(bs)
	stream.Write(sig.RBuf)
	bs = []byte{sig.SType, sig.SSize}
	stream.Write(bs)
	stream.Write(sig.SBuf)
	binary.LittleEndian.PutUint32(bs4, uint32(sig.SignatureType))
	stream.Write(bs4)
	bs = []byte{sig.PubKeySize}
	stream.Write(bs)
	stream.Write(sig.PubKey)
}

func (sig *SigHash) DeserializeSigHashFromByte(buf []byte) {
	byteBuf := bytes.NewBuffer(buf)
	sig.DeserializeSigHash(byteBuf)
}

func (sig *SigHash) DeserializeSigHash(byteBuf *bytes.Buffer) {
	binary.Read(byteBuf, binary.LittleEndian, &sig.DER)
	binary.Read(byteBuf, binary.LittleEndian, &sig.SignatureSize)
	binary.Read(byteBuf, binary.LittleEndian, &sig.RType)
	binary.Read(byteBuf, binary.LittleEndian, &sig.RSize)
	sig.RBuf = make([]byte, sig.RSize)
	binary.Read(byteBuf, binary.LittleEndian, sig.RBuf)
	binary.Read(byteBuf, binary.LittleEndian, &sig.SType)
	binary.Read(byteBuf, binary.LittleEndian, &sig.SSize)
	sig.SBuf = make([]byte, sig.SSize)
	binary.Read(byteBuf, binary.LittleEndian, sig.SBuf)
	binary.Read(byteBuf, binary.LittleEndian, &sig.SignatureType)
	binary.Read(byteBuf, binary.LittleEndian, &sig.PubKeySize)
	sig.PubKey = make([]byte, sig.PubKeySize)
	binary.Read(byteBuf, binary.LittleEndian, sig.PubKey)
}
