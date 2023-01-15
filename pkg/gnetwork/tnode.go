package gnetwork

import (
	"bytes"
	"container/list"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gsql"
)

type TNode struct {
	level        int
	maxNodeDepth int
	childNum     int
	nodeKey      byte
	nodes        map[byte]*TNode
	keyMap       map[byte]struct{}
	pubKeyAddr   string
	masterStore  gsql.MasterNodeStore
}

const ByteMaxValue = '{' - '0'

func NewTNode(owner string, masterStore gsql.MasterNodeStore, level int, maxNodeDepth int, key byte) *TNode {
	tNode := &TNode{
		level:        level,
		maxNodeDepth: maxNodeDepth,
		childNum:     0,
		nodeKey:      key,
		pubKeyAddr:   owner,
		masterStore:  masterStore,
	}

	if level < maxNodeDepth {
		tNode.nodes = make(map[byte]*TNode)
		tNode.keyMap = make(map[byte]struct{})
	}

	return tNode
}

func (tNode *TNode) AddChildNode(pubKey string) {
	if tNode.level == tNode.maxNodeDepth {
		return
	}

	nextKey := pubKey[tNode.level]
	if _, exist := tNode.nodes[nextKey]; exist == false {
		tNode.nodes[nextKey] = NewTNode(tNode.pubKeyAddr, tNode.masterStore, tNode.level+1, tNode.maxNodeDepth, nextKey)
		tNode.keyMap[nextKey] = struct{}{}
	}
	tNode.childNum++
	tNode.nodes[nextKey].AddChildNode(pubKey)
}

func (tNode *TNode) RemoveChildNode(pubKey string) {
	if tNode.level == tNode.maxNodeDepth {
		return
	}

	key := pubKey[tNode.level]
	node, exist := tNode.nodes[key]
	if exist == false {
		return
	}
	node.RemoveChildNode(pubKey)
	if node.childNum == 0 {
		delete(tNode.nodes, key)
		delete(tNode.keyMap, key)
	}

}

func (tNode *TNode) MakePubKey(pubKey string) {
	builder := bytes.NewBufferString(pubKey)
	builder.WriteByte(tNode.nodeKey)
	if tNode.level == tNode.maxNodeDepth {
		return
	}
	for char, _ := range tNode.keyMap {
		tNode.nodes[char].MakePubKey(builder.String())
	}
}

func (tNode *TNode) MakePubKeyFull(pubKey string, result []string) {
	builder := bytes.NewBufferString(pubKey)
	builder.WriteByte(tNode.nodeKey)
	if tNode.level == tNode.maxNodeDepth {
		result = append(result, builder.String())
		return
	}
	for char, _ := range tNode.keyMap {
		tNode.nodes[char].MakePubKeyFull(builder.String(), result)
	}
}

func (tNode *TNode) GetLevelMasterList(targetLevel int, searchPubKey string) *list.List {
	userList := list.New()
	builder := bytes.NewBufferString(searchPubKey)
	builder.WriteByte(tNode.nodeKey)

	if tNode.level == tNode.maxNodeDepth {
		tNode.masterStore.GetMasterNodeSearch(builder.String())
	} else if targetLevel == tNode.level {
		for c, _ := range tNode.keyMap {
			builder.WriteByte(c)
			node := tNode.masterStore.GetMasterNodeSearchPick(builder.String())
			if node != nil {
				userList.PushFront(node)
			}
		}
	} else {
		c := tNode.pubKeyAddr[tNode.level]
		userList = tNode.nodes[c].GetLevelMasterList(tNode.level, builder.String())
	}
	return userList
}

func (tNode *TNode) GetEnoughCluster(num int, pubKey string, searchPubKey string) string {
	key := pubKey[tNode.level]
	builder := bytes.NewBufferString(searchPubKey)
	builder.WriteByte(tNode.nodeKey)

	if tNode.level == tNode.maxNodeDepth {
		return builder.String()
	}

	if tNode.nodes[key].childNum > num {
		return tNode.nodes[key].GetEnoughCluster(num, pubKey, builder.String())
	}

	return builder.String()
}
