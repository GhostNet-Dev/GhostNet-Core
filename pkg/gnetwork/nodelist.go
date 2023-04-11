package gnetwork

import (
	"log"
	"net"
	"strings"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"google.golang.org/protobuf/proto"
)

type GhostAccount struct {
	nodeList         map[string]*GhostNode
	masterNodeList   map[string]*GhostNode
	nicknameToPubKey map[string]*GhostNode
	liteStore        *store.LiteStore
}

type GhostNode struct {
	User    *ptypes.GhostUser
	NetAddr *net.UDPAddr
}

const MaxGetherNodeList = 10

func NewGhostAccount(liteStore *store.LiteStore) *GhostAccount {
	account := &GhostAccount{
		nodeList:         make(map[string]*GhostNode),
		masterNodeList:   make(map[string]*GhostNode),
		nicknameToPubKey: make(map[string]*GhostNode),
		liteStore:        liteStore,
	}

	account.LoadNodeList(store.DefaultNickTable)
	account.LoadNodeList(store.DefaultMastersTable)

	return account
}

func (account *GhostAccount) LoadNodeList(table string) {
	if _, v, err := account.liteStore.LoadEntry(table); err == nil {
		for _, nodeByte := range v {
			masterNode := &ptypes.GhostUser{}
			if err := proto.Unmarshal(nodeByte, masterNode); err != nil {
				log.Fatal(err)
			}
			if table == store.DefaultMastersTable {
				account.AddMasterNode(&GhostNode{User: masterNode, NetAddr: masterNode.Ip.GetUdpAddr()})
			}
		}
	}
}

func (account *GhostAccount) AddUserNode(node *GhostNode) {
	account.nodeList[node.User.PubKey] = node
	account.nicknameToPubKey[node.User.Nickname] = node
	account.SaveToDb(store.DefaultNickTable, node.User)
}

func (account *GhostAccount) AddMasterNode(node *GhostNode) {
	log.Print("Add Master = ", node.User.Nickname)
	account.masterNodeList[node.User.PubKey] = node
	account.nicknameToPubKey[node.User.Nickname] = node
	account.SaveToDb(store.DefaultNickTable, node.User)
	account.SaveToDb(store.DefaultMastersTable, node.User)
}

func (account *GhostAccount) AddMasterUserList(userList []*ptypes.GhostUser) {
	for _, user := range userList {
		account.masterNodeList[user.PubKey] = &GhostNode{
			User:    user,
			NetAddr: user.Ip.GetUdpAddr(),
		}
		log.Print("Add Master = ", user.Nickname)
		account.SaveToDb(store.DefaultMastersTable, user)
	}
}

func (account *GhostAccount) SaveToDb(table string, masterNode *ptypes.GhostUser) {
	// Save To Db
	nodeByte, err := proto.Marshal(masterNode)
	if err != nil {
		log.Fatal(err)
	}
	if err := account.liteStore.SaveEntry(table, []byte(masterNode.Nickname), nodeByte); err != nil {
		log.Fatal(err)
	}
}

func (account *GhostAccount) LoadFromDb(nickname string) (masterNode *ptypes.GhostUser, err error) {
	nodeByte, err := account.liteStore.SelectEntry(store.DefaultNodeTable, []byte(nickname))
	if err != nil || nodeByte == nil {
		return nil, err
	}
	if err := proto.Unmarshal(nodeByte, masterNode); err != nil {
		log.Fatal(err)
	}
	return masterNode, nil
}

func (account *GhostAccount) GetUserNode(pubKey string) *GhostNode {
	find, exist := account.nodeList[pubKey]
	if !exist {
		log.Fatal("pubkey not found")
	}
	return find
}

func (account *GhostAccount) GetNodeInfo(pubKey string) *GhostNode {
	find, exist := account.masterNodeList[pubKey]
	if !exist {
		if find, exist = account.nodeList[pubKey]; !exist {
			return nil
		}
	}
	return find
}

func (account *GhostAccount) GetNodeByNickname(nickname string) *GhostNode {
	find, exist := account.nicknameToPubKey[nickname]
	if !exist {
		if node, err := account.LoadFromDb(nickname); err == nil && node != nil {
			ghostNode := &GhostNode{
				User:    node,
				NetAddr: node.Ip.GetUdpAddr(),
			}
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
			userList = append(userList, item.User)
		}
		i++
	}
	return userList, totalItem
}

func (account *GhostAccount) GetMasterNodeList() []*ptypes.GhostUser {
	userList := []*ptypes.GhostUser{}
	for _, item := range account.masterNodeList {
		userList = append(userList, item.User)
	}
	return userList
}

func (account *GhostAccount) GetMasterNodeSearch(pubKey string) []*ptypes.GhostUser {
	userList := []*ptypes.GhostUser{}
	for _, item := range account.masterNodeList {
		if strings.HasPrefix(item.User.PubKey, pubKey) {
			userList = append(userList, item.User)
		}
	}
	return userList
}

func (account *GhostAccount) GetMasterNodeSearchPick(pubKey string) *ptypes.GhostUser {
	for _, item := range account.masterNodeList {
		if strings.HasPrefix(item.User.PubKey, pubKey) {
			return item.User
		}
	}
	return nil
}
