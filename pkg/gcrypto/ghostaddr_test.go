package gcrypto

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrivateKeySerialize(t *testing.T) {
	ghostAddress := GenerateKeyPair()
	pubKey := ghostAddress.GetPubAddress()
	der := ghostAddress.PrivateKeySerialize()
	newGhostAddress := &GhostAddress{}
	newGhostAddress.PrivateKeyDeserialize(der)
	newPubKey := newGhostAddress.GetPubAddress()
	assert.Equal(t, pubKey, newPubKey, "different between new and old")
}
