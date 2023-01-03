package gnetwork

import (
	"log"
	"net"
	"strconv"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
)

type MasterNode struct {
	User    *ptypes.GhostUser
	NetAddr *net.UDPAddr
}

const MaxGetherNodeList = 10

func (master *MasterNetwork) AddMasterNode(node *MasterNode) {
	master.nodeList[node.User.PubKey] = node
	master.nicknameToPubKey[node.User.Nickname] = node
}

func (master *MasterNetwork) AddMasterUserList(userList []*ptypes.GhostUser) {
	for _, user := range userList {
		portInt, _ := strconv.Atoi(user.Ip.Port)
		master.nodeList[user.PubKey] = &MasterNode{
			User: user,
			NetAddr: &net.UDPAddr{
				Port: portInt,
				IP:   net.ParseIP(user.Ip.Ip),
			},
		}
	}
}

func (master *MasterNetwork) GetMasterNode(pubKey string) *MasterNode {
	find, err := master.nodeList[pubKey]
	if err == false {
		log.Fatal("pubkey not found")
	}
	return find
}

func (master *MasterNetwork) GetMasterNodeByNickname(nickname string) *MasterNode {
	find, err := master.nicknameToPubKey[nickname]
	if err == false {
		log.Fatal("nickname not found")
	}
	return find
}

func (master *MasterNetwork) GetMasterNodeUserList(startIndex uint32) ([]*ptypes.GhostUser, uint32) {
	startItem := startIndex * MaxGetherNodeList
	endItem := startIndex + MaxGetherNodeList
	totalItem := uint32(len(master.nodeList))
	if startItem > totalItem {
		return nil, 0
	}

	if endItem > totalItem {
		endItem = totalItem
	}

	//TODO: need to load from database
	i := uint32(0)
	userList := []*ptypes.GhostUser{}
	for _, item := range master.nodeList {
		if i >= startItem || i < endItem {
			userList = append(userList, item.User)
		}
		i++
	}
	return userList, totalItem
}
