package gnetwork

import (
	"log"
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
)

type TrieTreeMap struct {
	keyMap       map[byte]struct{}
	nodes        map[byte]*TNode
	checkList    map[string]byte
	ownerPubKey  string
	account      *GhostAccount
	currentLevel uint32
}

const (
	ClientTreeLevel  = 0
	DefaultTreeLevel = 1
	MaxNodeDepth     = 5
)

func NewTrieTreeMap(owner string, account *GhostAccount) *TrieTreeMap {
	return &TrieTreeMap{
		keyMap:       make(map[byte]struct{}),
		nodes:        make(map[byte]*TNode),
		checkList:    make(map[string]byte),
		ownerPubKey:  owner,
		account:      account,
		currentLevel: DefaultTreeLevel,
	}
}

func (tTree *TrieTreeMap) LoadTrieTree() {
	for _, user := range tTree.account.GetMasterNodeList() {
		tTree.AddNode(user.PubKey)
	}
}

func (tTree *TrieTreeMap) AddNode(pubKey string) {
	if _, exist := tTree.checkList[pubKey]; exist == true {
		return
	}

	key := pubKey[tTree.currentLevel]
	tTree.checkList[pubKey] = key
	if _, exist := tTree.nodes[key]; exist == false {
		tTree.nodes[key] = NewTNode(tTree.ownerPubKey, tTree.account,
			tTree.currentLevel+1, key)
		tTree.keyMap[key] = struct{}{}
	}

	tTree.nodes[key].AddChildNode(pubKey)
}

func (tTree *TrieTreeMap) DelNode(pubKey string) {
	if _, exist := tTree.checkList[pubKey]; exist == false {
		return
	}
	delete(tTree.checkList, pubKey)
	key := pubKey[tTree.currentLevel]
	if node, exist := tTree.nodes[key]; exist == false {
		return
	} else {
		if node.RemoveChildNode(pubKey) == true {
			node.childNum--
		}
		if node.childNum == 0 {
			delete(tTree.nodes, key)
			delete(tTree.keyMap, key)
		}
	}
}

func (tTree *TrieTreeMap) GetTotalNodeNum() int {
	totalNum := 0
	for c := range tTree.keyMap {
		totalNum += tTree.nodes[c].GetTotalNodeNum()
	}
	return totalNum
}

func (tTree *TrieTreeMap) GetLevelMasterList(targetLevel uint32) []*ptypes.GhostUser {
	if targetLevel < DefaultTreeLevel {
		log.Fatal("Level 0 is not exist!!")
	}

	c := tTree.ownerPubKey[tTree.currentLevel]
	if node, exist := tTree.nodes[c]; targetLevel > DefaultTreeLevel && exist == true {
		return node.GetLevelMasterList(targetLevel, "1")
	}

	userList := []*ptypes.GhostUser{}
	for c := range tTree.keyMap {
		node := tTree.account.GetMasterNodeSearchPick("1" + string(c))
		if node != nil {
			userList = append(userList, node)
		}
	}
	return userList
}

func (tTree *TrieTreeMap) GetCharLevelMasterList(searchPubKey string, targetLevel uint32) []*ptypes.GhostUser {
	c := searchPubKey[1]
	if _, exist := tTree.keyMap[c]; exist == false || targetLevel < DefaultTreeLevel {
		log.Fatal("Level 0 or char is not exist!!")
	}

	if targetLevel > DefaultTreeLevel {
		return tTree.nodes[c].GetLevelMasterList(targetLevel, searchPubKey)
	}
	return nil
}

func (tTree *TrieTreeMap) GetTreeClusterPick(searchString string) *net.UDPAddr {
	user := tTree.account.GetMasterNodeSearchPick("1" + searchString[0:1])
	if user == nil {
		for pubKey := range tTree.checkList {
			node := tTree.account.GetNodeInfo(pubKey)
			return node.NetAddr
		}
	}
	return user.Ip.GetUdpAddr()
}
