package crypto

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSignerVerify(t *testing.T) {
	ghostAddr := GenerateKeyPair()
	text := []byte("hello")
	sig := Signer(text, ghostAddr)

	assert.Equal(t, true, SignVerify(sig), "인증에 실패했습니다.")
}

func TestPrivateSerializeDesrialize(t *testing.T) {
	ghostAddr := GenerateKeyPair()
	buf := ghostAddr.PrivateKeySerialize()
	ghostAddr.PrivateKeyDeserialize(buf)
	buf2 := ghostAddr.PrivateKeySerialize()
	result := bytes.Compare(buf2, buf)
	assert.Equal(t, 0, result, "binary가 다릅니다.")
}
