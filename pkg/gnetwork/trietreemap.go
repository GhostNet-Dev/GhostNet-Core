package gnetwork

import (
	"container/list"
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gsql"
)

type TrieTreeMap struct {
	keyMap       map[byte]struct{}
	nodes        map[byte]*TNode
	checkList    map[string]byte
	ownerPubKey  string
	store        gsql.MasterNodeStore
	currentLevel int
}

const (
	ClientTreeLevel  = 0
	DefaultTreeLevel = 1
	MaxNodeDepth     = 5
)

func NewTrieTreeMap(owner string, store gsql.MasterNodeStore) *TrieTreeMap {
	return &TrieTreeMap{
		keyMap:       make(map[byte]struct{}),
		nodes:        make(map[byte]*TNode),
		checkList:    make(map[string]byte),
		ownerPubKey:  owner,
		store:        store,
		currentLevel: DefaultTreeLevel,
	}
}

func (tTree *TrieTreeMap) LoadTrieTree() {
	for _, user := range tTree.store.GetMasterNodeList() {
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
		tTree.nodes[key] = NewTNode(tTree.ownerPubKey, tTree.store,
			tTree.currentLevel+1, MaxNodeDepth, key)
		tTree.keyMap[key] = struct{}{}
	}

	tTree.nodes[key].AddChildNode(pubKey)
}

func (tTree *TrieTreeMap) DelNode(pubKey string) {
	if _, exist := tTree.checkList[pubKey]; exist == true {
		return
	}
	delete(tTree.checkList, pubKey)
	key := pubKey[tTree.currentLevel]
	if node, exist := tTree.nodes[key]; exist == false {
		return
	} else {
		node.RemoveChildNode(pubKey)
		if node.childNum == 0 {
			delete(tTree.nodes, key)
			delete(tTree.keyMap, key)
		}
	}
}

func (tTree *TrieTreeMap) GetLevelMasterList(targetLevel int) *list.List {
	if targetLevel < DefaultTreeLevel {
		log.Fatal("Level 0 is not exist!!")
	}

	if targetLevel > DefaultTreeLevel {
		c := tTree.ownerPubKey[tTree.currentLevel]
		return tTree.nodes[c].GetLevelMasterList(targetLevel, "1")
	}

	userList := list.New()
	for c := range tTree.keyMap {
		node := tTree.store.GetMasterNodeSearchPick("1" + string(c))
		if node != nil {
			userList.PushFront(node)
		}
	}
	return userList
}
