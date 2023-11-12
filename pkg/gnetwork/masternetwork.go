package gnetwork

import (
	"log"
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"google.golang.org/protobuf/proto"
)

type MasterNetwork struct {
	// My Nickname
	nickname string
	// My Master Node, not me
	masterInfo *ptypes.GhostUser
	// GhostNetVersion
	ghostNetVersion uint32
	// connected Master Nodes
	udp            *p2p.UdpServer
	owner          *gcrypto.GhostAddress
	localGhostIp   *ptypes.GhostIp
	blockContainer *store.BlockContainer
	tTreeMap       *TrieTreeMap
	account        *GhostAccount
	blockHandlerSq func(*packets.Header, *net.UDPAddr) []p2p.ResponseHeaderInfo
	blockHandlerCq func(*packets.Header, *net.UDPAddr)

	packetSqHandler map[packets.PacketSecondType]p2p.FuncPacketHandler
	packetCqHandler map[packets.PacketSecondType]p2p.FuncPacketHandler
}

func NewMasterNode(ghostNetVersion uint32, w *gcrypto.Wallet, myIpAddr *ptypes.GhostIp,
	packetFactory *p2p.PacketFactory, udp *p2p.UdpServer,
	blockContainer *store.BlockContainer, account *GhostAccount, tTreeMap *TrieTreeMap) *MasterNetwork {
	masterNode := &MasterNetwork{
		nickname:        w.GetNickname(),
		udp:             udp,
		owner:           w.GetGhostAddress(),
		localGhostIp:    myIpAddr,
		ghostNetVersion: ghostNetVersion,
		blockContainer:  blockContainer,
		account:         account,
		tTreeMap:        tTreeMap,
		masterInfo:      w.GetGhostUser(),
	}
	masterNode.RegisterHandler(packetFactory)
	tTreeMap.LoadTrieTree()
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
	node.packetSqHandler[packets.PacketSecondType_SearchUserInfoByPubKey] = node.SearchMasterPubKeySq
	node.packetSqHandler[packets.PacketSecondType_BlockChain] = node.BlockChainSq
	node.packetSqHandler[packets.PacketSecondType_Forwarding] = node.ForwardingSq

	node.packetCqHandler[packets.PacketSecondType_GetGhostNetVersion] = node.GetGhostNetVersionCq
	node.packetCqHandler[packets.PacketSecondType_NotificationMasterNode] = node.NotificationMasterNodeCq
	node.packetCqHandler[packets.PacketSecondType_ConnectToMasterNode] = node.ConnectToMasterNodeCq
	node.packetCqHandler[packets.PacketSecondType_SearchGhostPubKey] = node.SearchGhostPubKeyCq
	node.packetCqHandler[packets.PacketSecondType_RequestMasterNodeList] = node.RequestMasterNodeListCq
	node.packetCqHandler[packets.PacketSecondType_ResponseMasterNodeList] = node.ResponseMasterNodeListCq
	node.packetCqHandler[packets.PacketSecondType_SearchUserInfoByPubKey] = node.SearchMasterPubKeyCq
	node.packetCqHandler[packets.PacketSecondType_BlockChain] = node.BlockChainCq
	node.packetCqHandler[packets.PacketSecondType_Forwarding] = node.ForwardingCq

	packetFactory.RegisterPacketHandler(packets.PacketType_MasterNetwork, node.packetSqHandler, node.packetCqHandler)
}

func (master *MasterNetwork) RegisterMyMasterNode(user *ptypes.GhostUser) {
	master.masterInfo = user
	master.account.AddMasterNode(master.masterInfo)
	master.tTreeMap.AddNode(user.GetPubKey())
}

func (master *MasterNetwork) RegisterMasterNode(user *ptypes.GhostUser) {
	master.account.AddMasterNode(user)
	master.tTreeMap.AddNode(user.GetPubKey())
}

func (master *MasterNetwork) getGhostUser() *ptypes.GhostUser {
	return &ptypes.GhostUser{
		Nickname:     master.nickname,
		PubKey:       master.owner.GetPubAddress(),
		MasterPubKey: master.masterInfo.PubKey,
		Ip:           master.localGhostIp,
	}
}

func (master *MasterNetwork) RequestGhostNetVersion() {
	if master.masterInfo == nil {
		return
	}
	sq := packets.VersionInfoSq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), nil, 0, master.localGhostIp),
	}
	sendData, err := proto.Marshal(&sq)
	if err != nil {
		log.Fatal(err)
	}
	headerInfo := &p2p.ResponseHeaderInfo{
		ToAddr:     master.masterInfo.Ip.GetUdpAddr(),
		PacketType: packets.PacketType_MasterNetwork,
		SecondType: packets.PacketSecondType_NotificationMasterNode,
		RequestId:  sq.Master.GetRequestId(),
		PacketData: sendData,
		SqFlag:     true,
	}
	master.udp.SendUdpPacket(headerInfo, master.masterInfo.Ip.GetUdpAddr())
}

func (master *MasterNetwork) RequestMasterNodeList(index uint32, toAddr *net.UDPAddr) {
	if master.masterInfo == nil {
		return
	}
	sq := packets.RequestMasterNodeListSq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), nil, 0, master.localGhostIp),
	}
	sendData, err := proto.Marshal(&sq)
	if err != nil {
		log.Fatal(err)
	}
	headerInfo := &p2p.ResponseHeaderInfo{
		ToAddr:     master.masterInfo.Ip.GetUdpAddr(),
		PacketType: packets.PacketType_MasterNetwork,
		SecondType: packets.PacketSecondType_RequestMasterNodeList,
		RequestId:  sq.Master.GetRequestId(),
		PacketData: sendData,
		SqFlag:     true,
	}
	master.udp.SendUdpPacket(headerInfo, master.masterInfo.Ip.GetUdpAddr())
}

func (master *MasterNetwork) ConnectToMasterNode() {
	if master.masterInfo == nil {
		return
	}
	master.sendMasterUserInfo(packets.PacketSecondType_ConnectToMasterNode)
}

func (master *MasterNetwork) BroadcastMasterNodeNotification() {
	if master.masterInfo == nil {
		return
	}
	master.sendMasterUserInfo(packets.PacketSecondType_NotificationMasterNode)
}

func (master *MasterNetwork) sendMasterUserInfo(secondType packets.PacketSecondType) {

	sq := packets.MasterNodeUserInfoSq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), nil, 0, master.localGhostIp),
		User:   master.getGhostUser(),
	}
	sendData, err := proto.Marshal(&sq)
	if err != nil {
		log.Fatal(err)
	}
	headerInfo := &p2p.ResponseHeaderInfo{
		PacketType: packets.PacketType_MasterNetwork,
		SecondType: secondType,
		RequestId:  sq.Master.GetRequestId(),
		PacketData: sendData,
		SqFlag:     true,
	}
	master.SendToMasterNodeGrpSq(packets.RoutingType_BroadCasting, DefaultTreeLevel, headerInfo)
}

func (master *MasterNetwork) SendToMasterNodeSq(third packets.PacketThirdType, pubKey string, packet []byte, reqId []byte) {
	node := master.account.GetNodeInfo(pubKey)
	if node == nil {
		log.Fatal("node key not found")
	}
	headerInfo := &p2p.ResponseHeaderInfo{
		ToAddr:     node.Ip.GetUdpAddr(),
		PacketType: packets.PacketType_MasterNetwork,
		SecondType: packets.PacketSecondType_BlockChain,
		ThirdType:  third,
		RequestId:  reqId,
		PacketData: packet,
		SqFlag:     true,
	}
	master.udp.SendUdpPacket(headerInfo, headerInfo.ToAddr)
}

func (master *MasterNetwork) SendToMasterNodeGrpSq(routingType packets.RoutingType, level uint32, headerInfo *p2p.ResponseHeaderInfo) {
	if headerInfo.PacketType == packets.PacketType_Reserved0 {
		log.Fatal("PacketType not defined.")
	}
	// p2p.PacketHeaderInfo -> packets.Header
	header := master.udp.TranslationToHeader(headerInfo)
	packetList := master.makeForwadingPacket(routingType, level, header)
	for _, packet := range packetList {
		if packet.ToAddr.String() == master.localGhostIp.GetUdpAddr().String() {
			continue
		}
		master.udp.SendUdpPacket(&packet, packet.ToAddr)
	}
}

// routingtype, level, packet
func (master *MasterNetwork) makeForwadingPacket(routingType packets.RoutingType, level uint32, header *packets.Header) []p2p.ResponseHeaderInfo {
	if level > MaxNodeDepth || (routingType == packets.RoutingType_BroadCastingLevelZero && level > 1) {
		return nil
	}
	forwardingSq := packets.ForwardingSq{
		Master:           p2p.MakeMasterPacket(master.owner.GetPubAddress(), nil, 0, master.localGhostIp),
		ForwardingHeader: header,
	}
	forwardingSq.Master.RoutingT = routingType
	forwardingSq.Master.Level = level
	forwardingData, err := proto.Marshal(&forwardingSq)
	if err != nil {
		log.Fatal(err)
	}
	packetList := []p2p.ResponseHeaderInfo{}

	userList := master.tTreeMap.GetLevelMasterList(level)
	for _, user := range userList {
		if user.PubKey == master.owner.GetPubAddress() {
			continue
		}
		ghostUser := user
		from, _ := net.ResolveUDPAddr("udp", ghostUser.Ip.Ip+":"+ghostUser.Ip.Port)
		packetList = append(packetList, p2p.ResponseHeaderInfo{
			ToAddr:     from,
			PacketType: packets.PacketType_MasterNetwork,
			SecondType: packets.PacketSecondType_Forwarding,
			RequestId:  p2p.NewRequestId(),
			PacketData: forwardingData,
			SqFlag:     true,
		})
	}
	return packetList
}
