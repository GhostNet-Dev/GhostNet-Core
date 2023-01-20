package gnetwork

import (
	"log"
	"net"
	"strconv"
	"strings"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
)

type GhostAccount struct {
	nodeList         map[string]*GhostNode
	masterNodeList   map[string]*GhostNode
	nicknameToPubKey map[string]*GhostNode
}

type GhostNode struct {
	User    *ptypes.GhostUser
	NetAddr *net.UDPAddr
}

const MaxGetherNodeList = 10

func NewGhostAccount() *GhostAccount {
	return &GhostAccount{
		nodeList:         make(map[string]*GhostNode),
		masterNodeList:   make(map[string]*GhostNode),
		nicknameToPubKey: make(map[string]*GhostNode),
	}
}

func (account *GhostAccount) AddUserNode(node *GhostNode) {
	account.nodeList[node.User.PubKey] = node
	account.nicknameToPubKey[node.User.Nickname] = node
}

func (account *GhostAccount) AddMasterNode(node *GhostNode) {
	account.masterNodeList[node.User.PubKey] = node
	account.nicknameToPubKey[node.User.Nickname] = node
}

func (account *GhostAccount) AddMasterUserList(userList []*ptypes.GhostUser) {
	for _, user := range userList {
		portInt, _ := strconv.Atoi(user.Ip.Port)
		account.masterNodeList[user.PubKey] = &GhostNode{
			User: user,
			NetAddr: &net.UDPAddr{
				Port: portInt,
				IP:   net.ParseIP(user.Ip.Ip),
			},
		}
	}
}

func (account *GhostAccount) GetUserNode(pubKey string) *GhostNode {
	find, err := account.nodeList[pubKey]
	if err == false {
		log.Fatal("pubkey not found")
	}
	return find
}

func (account *GhostAccount) GetNodeInfo(pubKey string) *GhostNode {
	find, err := account.masterNodeList[pubKey]
	if err == false {
		if find, err = account.nodeList[pubKey]; err == false {
			return nil
		}
	}
	return find
}

func (account *GhostAccount) GetNodeByNickname(nickname string) *GhostNode {
	find, err := account.nicknameToPubKey[nickname]
	if err == false {
		log.Fatal("nickname not found")
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
		if strings.HasPrefix(item.User.PubKey, pubKey) == true {
			userList = append(userList, item.User)
		}
	}
	return userList
}

func (account *GhostAccount) GetMasterNodeSearchPick(pubKey string) *ptypes.GhostUser {
	for _, item := range account.masterNodeList {
		if strings.HasPrefix(item.User.PubKey, pubKey) == true {
			return item.User
		}
	}
	return nil
}
