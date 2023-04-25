package gnetwork

import (
	"fmt"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/stretchr/testify/assert"
)

var (
	TestNode = 10
)

func TestMakeTrieTree(t *testing.T) {
	account := NewGhostAccount(liteStore)
	tTreeMap := NewTrieTreeMap(owner.GetPubAddress(), account)
	MakePatternNode(tTreeMap, account)
	tTreeMap.LoadTrieTree()
	for i := uint32(1); i <= MaxNodeDepth; i++ {
		ret := tTreeMap.GetLevelMasterList(i)
		if MaxNodeDepth != i {
			assert.Equal(t, 1, len(ret), fmt.Sprint("not expected node items = ", i))
		} else {
			assert.Equal(t, TestNode, len(ret), fmt.Sprint("not expected node items = ", i))
		}
	}
	assert.Equal(t, TestNode, tTreeMap.GetTotalNodeNum(), "total num is wrong")
}

func TestAddRemove(t *testing.T) {
	account := NewGhostAccount(liteStore)
	tTreeMap := NewTrieTreeMap(owner.GetPubAddress(), account)
	MakePatternNode(tTreeMap, account)
	tTreeMap.LoadTrieTree()
	tTreeMap.AddNode("1aaaab")
	tTreeMap.AddNode("1aaaac")
	tTreeMap.AddNode("1aaaad")
	assert.Equal(t, TestNode+3, tTreeMap.GetTotalNodeNum(), "total num is wrong")
	tTreeMap.DelNode("1aaaac")
	tTreeMap.DelNode("1aaaab")
	assert.Equal(t, TestNode+1, tTreeMap.GetTotalNodeNum(), "total num is wrong")
	tTreeMap.AddNode("1aaaab")
	tTreeMap.AddNode("1aaaac")
	tTreeMap.DelNode("1aaaad")
	assert.Equal(t, TestNode+2, tTreeMap.GetTotalNodeNum(), "total num is wrong")
}

func MakePatternNode(tTreeMap *TrieTreeMap, account *GhostAccount) {
	tTreeMap.ownerPubKey = "1aaaaa"
	for i := 0; i < TestNode; i++ {
		owner = gcrypto.GenerateKeyPair()
		ghostUser := &ptypes.GhostUser{
			PubKey: fmt.Sprint("1aaaa", i),
			Ip:     ghostIp,
		}
		account.AddMasterNode(ghostUser)
	}
}

func MakeDummyNode(account *GhostAccount) {
	for i := 0; i < TestNode; i++ {
		owner = gcrypto.GenerateKeyPair()
		ghostUser := &ptypes.GhostUser{
			PubKey: owner.GetPubAddress(),
			Ip:     ghostIp,
		}
		account.AddMasterNode(ghostUser)
	}
}
