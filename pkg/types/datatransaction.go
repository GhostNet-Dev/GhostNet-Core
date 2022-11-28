package types

import (
	"unsafe"

	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
)

type GhostDataTransaction struct {
	TxId           ghostBytes.HashBytes
	LogicalAddress uint64
	DataSize       uint32
	Data           []byte
}

func (dataTx *GhostDataTransaction) Size() uint32 {
	return uint32(ghostBytes.HashSize) + //address
		uint32(unsafe.Sizeof(dataTx.LogicalAddress)) +
		uint32(unsafe.Sizeof(dataTx.DataSize)) +
		dataTx.DataSize
}
