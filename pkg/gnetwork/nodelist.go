package gnetwork

import (
	"log"
	"strings"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"google.golang.org/protobuf/proto"
)

type GhostAccount struct {
	nodeList         map[string]*ptypes.GhostUser
	masterNodeList   map[string]*ptypes.GhostUser
	nicknameToPubKey map[string]*ptypes.GhostUser
	liteStore        *store.LiteStore
}

const MaxGetherNodeList = 10

// TODO: Warning!! give me Refactoring!!!
func NewGhostAccount(liteStore *store.LiteStore) *GhostAccount {
	account := &GhostAccount{
		nodeList:         make(map[string]*ptypes.GhostUser),
		masterNodeList:   make(map[string]*ptypes.GhostUser),
		nicknameToPubKey: make(map[string]*ptypes.GhostUser),
		liteStore:        liteStore,
	}

	account.loadNodeList(store.DefaultNickTable)
	account.loadNodeList(store.DefaultMastersTable)

	return account
}

func (account *GhostAccount) AddUserNode(node *ptypes.GhostUser) {
	account.nodeList[node.PubKey] = node
	account.nicknameToPubKey[node.Nickname] = node
	account.saveToDb(store.DefaultNickTable, node)
}

func (account *GhostAccount) AddMasterNode(node *ptypes.GhostUser) {
	log.Print("Add Master = ", node.Nickname)
	account.masterNodeList[node.PubKey] = node
	account.nicknameToPubKey[node.Nickname] = node
	account.saveToDb(store.DefaultNickTable, node)
	account.saveToDb(store.DefaultMastersTable, node)
}

func (account *GhostAccount) AddMasterUserList(userList []*ptypes.GhostUser) {
	for _, user := range userList {
		account.masterNodeList[user.PubKey] = user
		log.Print("Add Master = ", user.Nickname)
		account.saveToDb(store.DefaultMastersTable, user)
	}
}

func (account *GhostAccount) GetUserNode(pubKey string) *ptypes.GhostUser {
	find, exist := account.nodeList[pubKey]
	if !exist {
		log.Fatal("pubkey not found")
	}
	return find
}

func (account *GhostAccount) GetNodeInfo(pubKey string) *ptypes.GhostUser {
	find, exist := account.masterNodeList[pubKey]
	if !exist {
		if find, exist = account.nodeList[pubKey]; !exist {
			return nil
		}
	}
	return find
}

func (account *GhostAccount) GetNodeByNickname(nickname string) *ptypes.GhostUser {
	find, exist := account.nicknameToPubKey[nickname]
	if !exist {
		if node, err := account.loadFromDb(nickname); err == nil && node != nil {
			ghostNode := node
			account.AddUserNode(ghostNode)
			return ghostNode
		}
		log.Print("nickname not found = ", nickname)
	}
	return find
}

func (account *GhostAccount) GetMasterNodeUserList(startIndex uint32) ([]*ptypes.GhostUser, uint32) {
	startItem := startIndex * MaxGetherNodeList
	endItem := startIndex + MaxGetherNodeList
	totalItem := uint32(len(account.masterNodeList))
	if startItem > totalItem {
		return nil, 0
	}

	if endItem > totalItem {
		endItem = totalItem
	}

	//TODO: need to load from database
	i := uint32(0)
	userList := []*ptypes.GhostUser{}
	for _, item := range account.masterNodeList {
		if i >= startItem || i < endItem {
			userList = append(userList, item)
		}
		i++
	}
	return userList, totalItem
}

func (account *GhostAccount) GetMasterNodeList() []*ptypes.GhostUser {
	userList := []*ptypes.GhostUser{}
	for _, item := range account.masterNodeList {
		userList = append(userList, item)
	}
	return userList
}

func (account *GhostAccount) GetMasterNodeSearch(pubKey string) []*ptypes.GhostUser {
	userList := []*ptypes.GhostUser{}
	for _, item := range account.masterNodeList {
		if strings.HasPrefix(item.PubKey, pubKey) {
			userList = append(userList, item)
		}
	}
	return userList
}

func (account *GhostAccount) GetMasterNodeSearchPick(pubKey string) *ptypes.GhostUser {
	for _, item := range account.masterNodeList {
		if strings.HasPrefix(item.PubKey, pubKey) {
			return item
		}
	}
	return nil
}

func (account *GhostAccount) loadNodeList(table string) {
	if _, v, err := account.liteStore.LoadEntry(table); err == nil {
		for _, nodeByte := range v {
			masterNode := &ptypes.GhostUser{}
			if err := proto.Unmarshal(nodeByte, masterNode); err != nil {
				log.Fatal(err)
			}
			if table == store.DefaultMastersTable {
				account.AddMasterNode(masterNode)
			}
		}
	}
}
func (account *GhostAccount) saveToDb(table string, masterNode *ptypes.GhostUser) {
	// Save To Db
	nodeByte, err := proto.Marshal(masterNode)
	if err != nil {
		log.Fatal(err)
	}
	if err := account.liteStore.SaveEntry(table, []byte(masterNode.Nickname), nodeByte); err != nil {
		log.Fatal(err)
	}
}

func (account *GhostAccount) loadFromDb(nickname string) (masterNode *ptypes.GhostUser, err error) {
	nodeByte, err := account.liteStore.SelectEntry(store.DefaultNodeTable, []byte(nickname))
	if err != nil || nodeByte == nil {
		return nil, err
	}
	if err := proto.Unmarshal(nodeByte, masterNode); err != nil {
		log.Fatal(err)
	}
	return masterNode, nil
}
