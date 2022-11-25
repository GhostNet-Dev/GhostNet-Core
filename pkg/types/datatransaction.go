package types

import (
	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
)

type GhostDataTrasaction struct {
	TxId        ghostBytes.HashBytes
	UniqHashKey ghostBytes.HashBytes
	DataSize    uint32
	Data        []byte
}
