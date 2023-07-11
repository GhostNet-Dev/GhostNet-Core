package types

import (
	"crypto/sha256"
	"unsafe"

	"github.com/GhostNet-Dev/GhostNet-Core/libs/gbytes"
)

type GhostDataTransaction struct {
	TxId           gbytes.HashBytes
	LogicalAddress gbytes.HashBytes
	DataSize       uint32
	Data           []byte
}

func (dataTx *GhostDataTransaction) Size() uint32 {
	return uint32(gbytes.HashSize)*2 + //address
		uint32(unsafe.Sizeof(dataTx.DataSize)) +
		dataTx.DataSize
}

func (dataTx *GhostDataTransaction) GetHashKey() []byte {
	dummy := make([]byte, gbytes.HashSize)
	backup := dataTx.TxId
	dataTx.TxId = dummy
	hash := sha256.New()
	hash.Write(dataTx.SerializeToByte())
	hashKey := hash.Sum(nil)
	dataTx.TxId = backup
	return hashKey
}
