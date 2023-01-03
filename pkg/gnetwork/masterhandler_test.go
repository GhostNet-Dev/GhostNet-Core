package gnetwork

import (
	"log"
	"net"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

var (
	ipAddr = &ptypes.GhostIp{
		Ip:   "127.0.0.1",
		Port: "8888",
	}
	udp            = p2p.NewUdpServer(ipAddr.Ip, ipAddr.Port)
	packetFactory  = p2p.NewPacketFactory()
	blockContainer = store.NewBlockContainer()
	config         = gconfig.DefaultConfig()
	from, _        = net.ResolveUDPAddr("udp", ipAddr.Ip+":"+ipAddr.Port)

	owner  = gcrypto.GenerateKeyPair()
	master = NewMasterNode("test", owner, ipAddr, config, packetFactory, udp, blockContainer)
)

func TestGetVersionSq(t *testing.T) {
	sq := packets.VersionInfoSq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr),
	}

	sendData, err := proto.Marshal(&sq)
	if err != nil {
		log.Fatal(err)
	}

	responseInfos := master.GetGhostNetVersionSq(&packets.Header{PacketData: sendData}, from)
	cq := &packets.VersionInfoCq{}
	if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, packets.PacketType_MasterNetwork, responseInfos[0].PacketType, "packet type is wrong")
	assert.Equal(t, packets.PacketSecondType_GetGhostNetVersion, responseInfos[0].SecondType, "packet type is wrong")
	assert.Equal(t, master.config.GhostVersion, cq.Version, "packet response is wrong")
}

func TestNotifyMasterNodeSq(t *testing.T) {
	sq := &packets.MasterNodeUserInfoSq{Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr)}

	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}

	responseInfos := master.NotificationMasterNodeSq(&packets.Header{PacketData: sendData}, from)
	cq := &packets.MasterNodeUserInfoCq{}
	if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, packets.PacketType_MasterNetwork, responseInfos[0].PacketType, "packet type is wrong")
	assert.Equal(t, packets.PacketSecondType_NotificationMasterNode, responseInfos[0].SecondType, "packet type is wrong")
	assert.Equal(t, "test", cq.User.Nickname, "nickname is wrong")
	assert.Equal(t, owner.GetPubAddress(), cq.User.PubKey, "pubkey is wrong")
}

func TestConnectToMasterNodeSq(t *testing.T) {
	sq := &packets.MasterNodeUserInfoSq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr),
		User: &ptypes.GhostUser{
			Nickname: master.nickname,
			PubKey:   owner.GetPubAddress(),
		},
	}

	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}

	responseInfos := master.ConnectToMasterNodeSq(&packets.Header{PacketData: sendData}, from)
	cq := &packets.MasterNodeUserInfoCq{}
	if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, packets.PacketType_MasterNetwork, responseInfos[0].PacketType, "packet type is wrong")
	assert.Equal(t, packets.PacketSecondType_ConnectToMasterNode, responseInfos[0].SecondType, "packet type is wrong")
	assert.Equal(t, "test", cq.User.Nickname, "nickname is wrong")
	assert.Equal(t, owner.GetPubAddress(), cq.User.PubKey, "pubkey is wrong")
}

func TestRequestMasterNodeListSq(t *testing.T) {
	TestConnectToMasterNodeSq(t)

	sq := &packets.RequestMasterNodeListSq{
		Master:     p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr),
		StartIndex: 0,
	}

	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}

	responseInfos := master.RequestMasterNodeListSq(&packets.Header{PacketData: sendData}, from)
	resInfo := responseInfos[1]
	resSq := &packets.ResponseMasterNodeListSq{}
	if err := proto.Unmarshal(resInfo.PacketData, resSq); err != nil {
		log.Fatal(err)
	}
	user := resSq.User[0]
	assert.Equal(t, packets.PacketType_MasterNetwork, responseInfos[0].PacketType, "packet[0] type is wrong")
	assert.Equal(t, packets.PacketSecondType_RequestMasterNodeList, responseInfos[0].SecondType, "packet[0] type is wrong")
	assert.Equal(t, packets.PacketType_MasterNetwork, responseInfos[1].PacketType, "packet[1] type is wrong")
	assert.Equal(t, packets.PacketSecondType_ResponseMasterNodeList, responseInfos[1].SecondType, "packet[1] type is wrong")
	assert.Equal(t, "test", user.Nickname, "nickname is wrong")
	assert.Equal(t, owner.GetPubAddress(), user.PubKey, "pubkey is wrong")
}

func TestResponseMasterNodeListSq(t *testing.T) {
	sq := &packets.ResponseMasterNodeListSq{Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr)}

	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}

	responseInfos := master.ResponseMasterNodeListSq(&packets.Header{PacketData: sendData}, from)
	cq := &packets.ResponseMasterNodeListCq{}
	if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, packets.PacketType_MasterNetwork, responseInfos[0].PacketType, "packet type is wrong")
	assert.Equal(t, packets.PacketSecondType_ResponseMasterNodeList, responseInfos[0].SecondType, "packet type is wrong")
}

func TestSearchMasterPubKey(t *testing.T) {
	TestConnectToMasterNodeSq(t)
	sq := &packets.SearchGhostPubKeySq{
		Master:   p2p.MakeMasterPacket(master.owner.GetPubAddress(), 0, 0, master.ipAddr),
		Nickname: "test",
	}

	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}

	responseInfos := master.SearchMasterPubKeySq(&packets.Header{PacketData: sendData}, from)
	cq := &packets.SearchGhostPubKeyCq{}
	if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, packets.PacketType_MasterNetwork, responseInfos[0].PacketType, "packet type is wrong")
	assert.Equal(t, packets.PacketSecondType_SearchMasterPubKey, responseInfos[0].SecondType, "packet type is wrong")
	assert.Equal(t, owner.GetPubAddress(), cq.User[0].PubKey, "pubkey is wrong")
}
