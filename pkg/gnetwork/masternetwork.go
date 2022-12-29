package gnetwork

import (
	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
)

type MasterNetwork struct {
	// My Nickname
	nickname string
	// My Master Node, not me
	masterInfo *MasterNode
	// connected Master Nodes
	nodeList         map[string]*MasterNode
	nicknameToPubKey map[string]*MasterNode
	udp              *p2p.UdpServer
	owner            *gcrypto.GhostAddress
	ipAddr           *ptypes.GhostIp
	config           *gconfig.GConfig
	blockContainer   *store.BlockContainer
}

func NewMasterNode(nickname string, myAddr *gcrypto.GhostAddress, myIpAddr *ptypes.GhostIp,
	config *gconfig.GConfig, packetFactory *p2p.PacketFactory, udp *p2p.UdpServer,
	blockContainer *store.BlockContainer) *MasterNetwork {
	masterNode := &MasterNetwork{
		nickname:         nickname,
		nodeList:         make(map[string]*MasterNode),
		nicknameToPubKey: make(map[string]*MasterNode),
		udp:              udp,
		owner:            myAddr,
		ipAddr:           myIpAddr,
		blockContainer:   blockContainer,
	}
	packetFactory.RegisterPacketHandler(masterNode)
	return masterNode
}

func (node *MasterNetwork) BroadcastMasterNodeNotification() {

}
