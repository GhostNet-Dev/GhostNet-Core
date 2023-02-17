package bootloader

import (
	"fmt"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/stretchr/testify/assert"
)

var (
	gVm = gvm.NewGVM()
)

func TestLoadGenesis(t *testing.T) {
	loader := NewLoadGenesis(gVm, "../blocks/samples/")
	assert.Equal(t, true, loader.pairedBlock != nil, "genesis block is nil")
	creator := loader.CreatorList()
	assert.Equal(t, "Adam", creator["Adam"].Nickname, "Adam is not found")

	w, err := loader.LoadCreatorKeyFile("Adam", creator["Adam"].PubKey,
		gcrypto.PasswordToSha256("test"))
	assert.Equal(t, true, err == nil, fmt.Sprint("fail to load creatorkeyfile: ", err))
	assert.Equal(t, creator["Adam"].PubKey, w.GetPubAddress(), "difference between pubkeys")
}
