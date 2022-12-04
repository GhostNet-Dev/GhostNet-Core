package types

type SigHash struct {
	DER           byte
	SignatureSize byte
	RType         byte
	RSize         byte
	RBuf          []byte
	SType         byte
	SSize         byte
	SBuf          []byte
	SignatureType uint32
	PubKeySize    byte
	PubKey        []byte
}

func (sig *SigHash) SigSize() int32 {
	return int32(sig.SSize) + int32(sig.RSize) + 6
}

func (sig *SigHash) Size() int32 {
	return sig.SigSize() + int32(sig.PubKeySize) + 1 /*byte*/ + 4 /*uint32*/
}
