package gnetwork

import (
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
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
	tTreeMap         *TrieTreeMap
	blockHandlerSq   func(*packets.Header, *net.UDPAddr) []p2p.PacketHeaderInfo
	blockHandlerCq   func(*packets.Header, *net.UDPAddr)

	packetSqHandler map[packets.PacketSecondType]p2p.FuncPacketHandler
	packetCqHandler map[packets.PacketSecondType]p2p.FuncPacketHandler
}

func NewMasterNode(nickname string, myAddr *gcrypto.GhostAddress, myIpAddr *ptypes.GhostIp,
	config *gconfig.GConfig, packetFactory *p2p.PacketFactory, udp *p2p.UdpServer,
	blockContainer *store.BlockContainer, tTreeMap *TrieTreeMap) *MasterNetwork {
	masterNode := &MasterNetwork{
		nickname:         nickname,
		nodeList:         make(map[string]*MasterNode),
		nicknameToPubKey: make(map[string]*MasterNode),
		udp:              udp,
		owner:            myAddr,
		ipAddr:           myIpAddr,
		config:           config,
		blockContainer:   blockContainer,
		tTreeMap:         tTreeMap,
	}
	masterNode.RegisterHandler(packetFactory)
	return masterNode
}

func (node *MasterNetwork) RegisterHandler(packetFactory *p2p.PacketFactory) {
	node.packetSqHandler = make(map[packets.PacketSecondType]p2p.FuncPacketHandler)
	node.packetCqHandler = make(map[packets.PacketSecondType]p2p.FuncPacketHandler)

	node.packetSqHandler[packets.PacketSecondType_GetGhostNetVersion] = node.GetGhostNetVersionSq
	node.packetSqHandler[packets.PacketSecondType_NotificationMasterNode] = node.NotificationMasterNodeSq
	node.packetSqHandler[packets.PacketSecondType_ConnectToMasterNode] = node.ConnectToMasterNodeSq
	node.packetSqHandler[packets.PacketSecondType_SearchGhostPubKey] = node.SearchGhostPubKeySq
	node.packetSqHandler[packets.PacketSecondType_RequestMasterNodeList] = node.RequestMasterNodeListSq
	node.packetSqHandler[packets.PacketSecondType_ResponseMasterNodeList] = node.ResponseMasterNodeListSq
	node.packetSqHandler[packets.PacketSecondType_SearchMasterPubKey] = node.SearchMasterPubKeySq
	node.packetSqHandler[packets.PacketSecondType_BlockChain] = node.SearchMasterPubKeySq
	node.packetSqHandler[packets.PacketSecondType_Forwarding] = node.ForwardingSq

	node.packetCqHandler[packets.PacketSecondType_GetGhostNetVersion] = node.GetGhostNetVersionCq
	node.packetCqHandler[packets.PacketSecondType_NotificationMasterNode] = node.NotificationMasterNodeCq
	node.packetCqHandler[packets.PacketSecondType_ConnectToMasterNode] = node.ConnectToMasterNodeCq
	node.packetCqHandler[packets.PacketSecondType_SearchGhostPubKey] = node.SearchGhostPubKeyCq
	node.packetCqHandler[packets.PacketSecondType_RequestMasterNodeList] = node.RequestMasterNodeListCq
	node.packetCqHandler[packets.PacketSecondType_ResponseMasterNodeList] = node.ResponseMasterNodeListCq
	node.packetCqHandler[packets.PacketSecondType_SearchMasterPubKey] = node.SearchMasterPubKeyCq
	node.packetCqHandler[packets.PacketSecondType_BlockChain] = node.SearchMasterPubKeyCq
	node.packetCqHandler[packets.PacketSecondType_Forwarding] = node.ForwardingCq

	packetFactory.RegisterPacketHandler(packets.PacketType_MasterNetwork, node.packetSqHandler, node.packetCqHandler)
}

func (node *MasterNetwork) BroadcastMasterNodeNotification() {

}
