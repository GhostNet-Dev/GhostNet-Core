package bootloader

import (
	"log"
	"math/rand"
	"net"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"google.golang.org/protobuf/proto"
)

/*
connect to master node
1. check exist node db
2. net open
3. request node list
4. random select
*/

type ConnectMaster struct {
	db              *store.LiteStore
	udp             *p2p.UdpServer
	wallet          *gcrypto.Wallet
	masterNode      *ptypes.GhostUser
	table           string
	eventChannel    chan bool
	eventWait       bool
	packetSqHandler map[packets.PacketSecondType]p2p.FuncPacketHandler
	packetCqHandler map[packets.PacketSecondType]p2p.FuncPacketHandler
}

const RootUrl = "www.ghostnetroot.com"

func NewConnectMaster(table string, db *store.LiteStore, packetFactory *p2p.PacketFactory,
	udp *p2p.UdpServer, w *gcrypto.Wallet) *ConnectMaster {

	conn := &ConnectMaster{
		db:              db,
		udp:             udp,
		wallet:          w,
		table:           table,
		eventChannel:    make(chan bool),
		eventWait:       false,
		packetSqHandler: make(map[packets.PacketSecondType]p2p.FuncPacketHandler),
		packetCqHandler: make(map[packets.PacketSecondType]p2p.FuncPacketHandler),
	}
	conn.RegisterPacketHandler(packetFactory)
	return conn
}

func GetRootIpAddress() *net.UDPAddr {
	ips, err := net.LookupIP(RootUrl)
	if err != nil {
		log.Fatal(err)
	}
	to, _ := net.ResolveUDPAddr("udp", ips[0].String()+":50029")
	return to
}

func GetGhostRootIpAddress() *ptypes.GhostIp {
	ips, err := net.LookupIP(RootUrl)
	if err != nil {
		log.Fatal(err)
	}
	return &ptypes.GhostIp{
		Ip:   ips[0].String(),
		Port: "50029",
	}
}

func (conn *ConnectMaster) SaveMasterNodeList(nodes []*ptypes.GhostUser) {
	for _, node := range nodes {
		nodeByte, err := proto.Marshal(node)
		if err != nil {
			log.Fatal(err)
		}
		if err := conn.db.SaveEntry(conn.table, []byte(node.Nickname), nodeByte); err != nil {
			log.Fatal(err)
		}
	}
}

func (conn *ConnectMaster) LoadMasterNode() *ptypes.GhostUser {
	_, nodes, err := conn.db.LoadEntry(conn.table)
	if err != nil {
		log.Fatal(err)
	}
	if nodes == nil {
		//need to request from adam
		return nil
	}
	randPick := rand.Uint32() % uint32(len(nodes))
	ghostUser := &ptypes.GhostUser{}
	for i := 0; i < len(nodes); i++ {
		node := nodes[randPick]
		randPick = (randPick + 1) % uint32(len(nodes))
		if err := proto.Unmarshal(node, ghostUser); err != nil {
			log.Fatal(err)
		}
		if ghostUser.PubKey == conn.wallet.GetPubAddress() {
			continue
		}
		break
	}
	if !ghostUser.Validate() {
		return nil
	}
	return ghostUser
}

func (conn *ConnectMaster) WaitEvent() (timeout bool) {
	conn.eventWait = true
	select {
	case <-conn.eventChannel:
		timeout = false
	case <-time.After(time.Second * 8):
		timeout = true
	}
	conn.eventWait = false
	return timeout
}

func (conn *ConnectMaster) GetGhostNetVersion(masterNode *ptypes.GhostUser) {
	sq := &packets.VersionInfoSq{
		Master: p2p.MakeMasterPacket(conn.wallet.GetGhostAddress().GetPubAddress(), nil, 0, conn.udp.GetLocalIp()),
	}
	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}
	headerInfo := &p2p.ResponseHeaderInfo{
		ToAddr:     masterNode.Ip.GetUdpAddr(),
		PacketType: packets.PacketType_MasterNetwork,
		SecondType: packets.PacketSecondType_GetGhostNetVersion,
		RequestId:  sq.Master.GetRequestId(),
		PacketData: sendData,
		SqFlag:     true,
	}
	conn.udp.SendUdpPacket(headerInfo, headerInfo.ToAddr)
}

func (conn *ConnectMaster) ConnectToMaster(masterNode *ptypes.GhostUser) {
	sq := &packets.MasterNodeUserInfoSq{
		Master: p2p.MakeMasterPacket(conn.wallet.GetPubAddress(), nil, 0, conn.udp.GetLocalIp()),
		User:   conn.wallet.GetGhostUser(),
	}
	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}
	headerInfo := &p2p.ResponseHeaderInfo{
		ToAddr:     masterNode.Ip.GetUdpAddr(),
		PacketType: packets.PacketType_MasterNetwork,
		SecondType: packets.PacketSecondType_ConnectToMasterNode,
		RequestId:  sq.Master.GetRequestId(),
		PacketData: sendData,
		SqFlag:     true,
	}
	conn.udp.SendUdpPacket(headerInfo, headerInfo.ToAddr)
}

func (conn *ConnectMaster) RequestMasterNodesToAdam() {
	sq := &packets.RequestMasterNodeListSq{
		Master:     p2p.MakeMasterPacket(conn.wallet.GetPubAddress(), nil, 0, conn.udp.GetLocalIp()),
		StartIndex: 0,
	}
	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}
	headerInfo := &p2p.ResponseHeaderInfo{
		ToAddr:     GetRootIpAddress(),
		PacketType: packets.PacketType_MasterNetwork,
		SecondType: packets.PacketSecondType_RequestMasterNodeList,
		RequestId:  sq.Master.GetRequestId(),
		PacketData: sendData,
		SqFlag:     true,
	}
	conn.udp.SendUdpPacket(headerInfo, headerInfo.ToAddr)
}

func (conn *ConnectMaster) RegisterPacketHandler(packetFactory *p2p.PacketFactory) {
	conn.packetSqHandler = make(map[packets.PacketSecondType]p2p.FuncPacketHandler)
	conn.packetCqHandler = make(map[packets.PacketSecondType]p2p.FuncPacketHandler)

	conn.packetSqHandler[packets.PacketSecondType_ResponseMasterNodeList] = conn.ResponseMasterNodeListSq

	conn.packetCqHandler[packets.PacketSecondType_GetGhostNetVersion] = conn.GetGhostNetVersionCq
	conn.packetCqHandler[packets.PacketSecondType_ConnectToMasterNode] = conn.ConnectToMasterNodeCq
	conn.packetCqHandler[packets.PacketSecondType_RequestMasterNodeList] = conn.RequestMasterNodeListCq

	packetFactory.RegisterPacketHandler(packets.PacketType_MasterNetwork, conn.packetSqHandler, conn.packetCqHandler)
}

func (conn *ConnectMaster) GetGhostNetVersionCq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	// TODO: Version Check
	return nil
}

func (conn *ConnectMaster) ConnectToMasterNodeCq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	cq := &packets.MasterNodeUserInfoCq{}
	if err := proto.Unmarshal(header.PacketData, cq); err != nil {
		log.Fatal(err)
	}
	conn.masterNode = cq.User
	if conn.eventWait {
		conn.eventChannel <- true
	}
	return nil
}

func (conn *ConnectMaster) RequestMasterNodeListCq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	return nil
}

func (conn *ConnectMaster) ResponseMasterNodeListSq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	sq := &packets.ResponseMasterNodeListSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	if len(sq.User) == 0 {
		sq.User = append(sq.User, &ptypes.GhostUser{
			PubKey:   gcrypto.Translate160ToBase58Addr(store.AdamsAddress()),
			Nickname: "Adam",
			Ip:       GetGhostRootIpAddress(),
		})
	}

	conn.SaveMasterNodeList(sq.User)
	if conn.eventWait {
		conn.eventChannel <- true
	}

	cq := packets.ResponseMasterNodeListCq{
		Master: p2p.MakeMasterPacket(conn.wallet.GetPubAddress(), sq.Master.GetRequestId(), 0, conn.udp.GetLocalIp()),
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
			SecondType: packets.PacketSecondType_ResponseMasterNodeList,
			PacketData: sendData,
			SqFlag:     false,
		},
	}
}
