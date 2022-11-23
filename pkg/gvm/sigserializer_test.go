package gvm

import (
	"bytes"
	"testing"
	"crypto/sha256"

	"github.com/stretchr/testify/assert"
	mems "github.com/traherom/memstream"
)


func TestSigHashSerializeDeserialize(t *testing.T) {
	dummy := make([]byte, 4)
	hash := sha256.New()
	hash.Write(dummy)
	key := hash.Sum((nil))
	sig := SigHash{
		SSize: 32,
		SBuf: key,
		RSize: 32,
		RBuf: key,
	}
	size := sig.Size()
	stream := mems.NewCapacity(int(size))
	sig.Serialize(stream)
	byteBuf := bytes.NewBuffer(stream.Bytes())

	newSig := SigHash{}
	newSig.DeserializeSigHash(byteBuf)
	assert.Equal(t, sig.SSize, newSig.SSize, "Value가 다릅니다.")

	result := bytes.Compare(sig.RBuf, newSig.RBuf)
	assert.Equal(t, 0, result, "buffer 내용이 다릅니다.")
}
