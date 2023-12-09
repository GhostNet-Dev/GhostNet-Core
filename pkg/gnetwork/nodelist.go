package gnetwork

import (
	"log"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"google.golang.org/protobuf/proto"
)

type GhostAccount struct {
	masterNodeList   *sync.Map
	masterNodeLen    *atomic.Uint32
	nodeList         *sync.Map
	nicknameToPubKey *sync.Map
	liteStore        *store.LiteStore
}

const MaxGetherNodeList = 10

// TODO: Warning!! give me Refactoring!!!
func NewGhostAccount(liteStore *store.LiteStore) *GhostAccount {
	account := &GhostAccount{
		masterNodeList:   new(sync.Map),
		masterNodeLen:    &atomic.Uint32{},
		nodeList:         new(sync.Map),
		nicknameToPubKey: new(sync.Map),
		liteStore:        liteStore,
	}

	account.loadNodeList(store.DefaultNickTable)
	account.loadNodeList(store.DefaultMastersTable)

	return account
}

func (account *GhostAccount) StoreMasterNode(node *ptypes.GhostUser) {
	account.masterNodeList.Store(node.PubKey, node)
	account.masterNodeLen.Add(1)
}

func (account *GhostAccount) AddUserNode(node *ptypes.GhostUser) {
	account.nodeList.Store(node.PubKey, node)
	account.nicknameToPubKey.Store(node.Nickname, node)
	account.saveToDb(store.DefaultNickTable, node)
}

func (account *GhostAccount) AddMasterNode(node *ptypes.GhostUser) {
	log.Print("Add Master = ", node.Nickname)
	account.StoreMasterNode(node)

	account.nicknameToPubKey.Store(node.Nickname, node)

	account.saveToDb(store.DefaultNickTable, node)
	account.saveToDb(store.DefaultMastersTable, node)
}

func (account *GhostAccount) AddMasterUserList(userList []*ptypes.GhostUser) {
	for _, user := range userList {
		account.StoreMasterNode(user)
		log.Print("Add Master = ", user.Nickname)
		account.saveToDb(store.DefaultMastersTable, user)
	}
}

func (account *GhostAccount) GetUserNode(pubKey string) *ptypes.GhostUser {
	find, exist := account.nodeList.Load(pubKey)
	if !exist {
		log.Fatal("pubkey not found")
	}
	return find.(*ptypes.GhostUser)
}

func (account *GhostAccount) GetNodeInfo(pubKey string) *ptypes.GhostUser {
	find, exist := account.masterNodeList.Load(pubKey)
	if !exist {
		if find, exist = account.nodeList.Load(pubKey); !exist {
			return nil
		}
	}
	return find.(*ptypes.GhostUser)
}

func (account *GhostAccount) GetNodeByNickname(nickname string) *ptypes.GhostUser {
	find, exist := account.nicknameToPubKey.Load(nickname)
	if !exist {
		if node, err := account.loadFromDb(nickname); err == nil && node != nil {
			ghostNode := node
			account.AddUserNode(ghostNode)
			return ghostNode
		}
		log.Print("nickname not found = ", nickname)
	}
	return find.(*ptypes.GhostUser)
}

func (account *GhostAccount) GetMasterNodeUserList(startIndex uint32) ([]*ptypes.GhostUser, uint32) {
	startItem := startIndex * MaxGetherNodeList
	endItem := startIndex + MaxGetherNodeList
	totalItem := account.masterNodeLen.Load()
	if startItem > totalItem {
		return nil, 0
	}

	if endItem > totalItem {
		endItem = totalItem
	}

	//TODO: need to load from database
	i := uint32(0)
	userList := []*ptypes.GhostUser{}
	account.masterNodeList.Range(func(key, value any) bool {
		if i >= startItem || i < endItem {
			userList = append(userList, value.(*ptypes.GhostUser))
		}
		i++
		return true
	})
	return userList, totalItem
}

func (account *GhostAccount) GetMasterNodeList() []*ptypes.GhostUser {
	userList := []*ptypes.GhostUser{}
	account.masterNodeList.Range(func(key, value any) bool {
		userList = append(userList, value.(*ptypes.GhostUser))
		return true
	})
	return userList
}

func (account *GhostAccount) GetMasterNodeSearch(pubKey string) []*ptypes.GhostUser {
	userList := []*ptypes.GhostUser{}
	account.masterNodeList.Range(func(key, value any) bool {
		if strings.HasPrefix(value.(*ptypes.GhostUser).PubKey, pubKey) {
			userList = append(userList, value.(*ptypes.GhostUser))
		}
		return true
	})
	return userList
}

func (account *GhostAccount) GetMasterNodeSearchPick(pubKey string) (ret *ptypes.GhostUser) {
	ret = nil
	account.masterNodeList.Range(func(key, value any) bool {
		if strings.HasPrefix(value.(*ptypes.GhostUser).PubKey, pubKey) {
			ret = value.(*ptypes.GhostUser)
			return false
		}
		return true
	})
	return ret
}

func (account *GhostAccount) loadNodeList(table string) {
	if k, v, err := account.liteStore.LoadEntry(table); err == nil {
		for idx, nodeByte := range v {
			masterNode := &ptypes.GhostUser{}
			if err := proto.Unmarshal(nodeByte, masterNode); err != nil {
				log.Fatal(err)
				account.liteStore.DelEntry(table, k[idx])
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
