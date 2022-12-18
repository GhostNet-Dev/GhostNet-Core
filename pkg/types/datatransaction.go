package types

import (
	"crypto/sha256"
	"unsafe"

	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
)

type GhostDataTransaction struct {
	TxId           gbytes.HashBytes
	LogicalAddress uint64
	DataSize       uint32
	Data           []byte
}

func (dataTx *GhostDataTransaction) Size() uint32 {
	return uint32(gbytes.HashSize) + //address
		uint32(unsafe.Sizeof(dataTx.LogicalAddress)) +
		uint32(unsafe.Sizeof(dataTx.DataSize)) +
		dataTx.DataSize
}

func (dataTx *GhostDataTransaction) GetHashKey() []byte {
	hash := sha256.New()
	hash.Write(dataTx.SerializeToByte())
	return hash.Sum(nil)
}
