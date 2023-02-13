package gnetwork

import (
	"bytes"
	"fmt"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
)

type TNode struct {
	level      uint32
	childNum   int
	nodeKey    byte
	nodes      map[byte]*TNode
	keyMap     map[byte]struct{}
	pubKeyAddr string
	account    *GhostAccount
}

const ByteMaxValue = '{' - '0'

func NewTNode(owner string, account *GhostAccount, level uint32, key byte) *TNode {
	tNode := &TNode{
		level:      level,
		childNum:   0,
		nodeKey:    key,
		pubKeyAddr: owner,
		account:    account,
	}

	if level < MaxNodeDepth {
		tNode.nodes = make(map[byte]*TNode)
		tNode.keyMap = make(map[byte]struct{})
	}

	return tNode
}

func (tNode *TNode) AddChildNode(pubKey string) {
	if tNode.level == MaxNodeDepth {
		tNode.childNum++
		return
	}

	nextKey := pubKey[tNode.level]
	if _, exist := tNode.nodes[nextKey]; !exist {
		tNode.nodes[nextKey] = NewTNode(tNode.pubKeyAddr, tNode.account, tNode.level+1, nextKey)
		tNode.keyMap[nextKey] = struct{}{}
	}
	tNode.childNum++
	tNode.nodes[nextKey].AddChildNode(pubKey)
}

func (tNode *TNode) RemoveChildNode(pubKey string) bool {
	if tNode.level == MaxNodeDepth {
		tNode.childNum--
		return true
	}

	key := pubKey[tNode.level]
	node, exist := tNode.nodes[key]
	if !exist {
		return false
	}
	if node.RemoveChildNode(pubKey) {
		tNode.childNum--
	}
	if node.childNum == 0 {
		delete(tNode.nodes, key)
		delete(tNode.keyMap, key)
		return true
	}
	return false
}

func (tNode *TNode) MakePubKey(pubKey string) {
	builder := bytes.NewBufferString(pubKey)
	builder.WriteByte(tNode.nodeKey)
	if tNode.level == MaxNodeDepth {
		return
	}
	for char := range tNode.keyMap {
		tNode.nodes[char].MakePubKey(builder.String())
	}
}

func (tNode *TNode) MakePubKeyFull(pubKey string, result []string) {
	builder := bytes.NewBufferString(pubKey)
	builder.WriteByte(tNode.nodeKey)
	if tNode.level == MaxNodeDepth {
		result = append(result, builder.String())
		return
	}
	for char := range tNode.keyMap {
		tNode.nodes[char].MakePubKeyFull(builder.String(), result)
	}
}

func (tNode *TNode) GetLevelMasterList(targetLevel uint32, buildString string) []*ptypes.GhostUser {
	builder := bytes.NewBufferString(buildString)
	builder.WriteByte(tNode.nodeKey)

	if tNode.level == MaxNodeDepth {
		return tNode.account.GetMasterNodeSearch(builder.String())
	} else if targetLevel == tNode.level {
		userList := []*ptypes.GhostUser{}
		for c := range tNode.keyMap {
			node := tNode.account.GetMasterNodeSearchPick(fmt.Sprint(builder.String(), string(c)))
			if node != nil {
				userList = append(userList, node)
			}
		}
		return userList
	} else {
		c := tNode.pubKeyAddr[tNode.level]
		if node, exist := tNode.nodes[c]; exist {
			return node.GetLevelMasterList(targetLevel, builder.String())
		} else {
			return nil
		}
	}
}

func (tNode *TNode) GetEnoughCluster(num int, pubKey string, searchPubKey string) string {
	key := pubKey[tNode.level]
	builder := bytes.NewBufferString(searchPubKey)
	builder.WriteByte(tNode.nodeKey)

	if tNode.level == MaxNodeDepth {
		return builder.String()
	}

	if tNode.nodes[key].childNum > num {
		return tNode.nodes[key].GetEnoughCluster(num, pubKey, builder.String())
	}

	return builder.String()
}

func (tNode *TNode) GetTotalNodeNum() int {
	if tNode.level == MaxNodeDepth-1 {
		return tNode.childNum
	}
	totalNum := 0
	for c := range tNode.keyMap {
		totalNum += tNode.nodes[c].GetTotalNodeNum()
	}
	return totalNum
}
