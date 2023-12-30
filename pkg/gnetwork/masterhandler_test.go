package gnetwork

import (
	"log"
	"net"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

var (
	ghostIp = &ptypes.GhostIp{
		Ip:   "127.0.0.1",
		Port: "8888",
	}
	TestTables     = []string{"nodes", "wallet"}
	liteStore      = store.NewLiteStore("./", "litestore.db", TestTables, 3)
	gAccount       = NewGhostAccount(liteStore)
	packetFactory  = p2p.NewPacketFactory()
	udp            = p2p.NewUdpServer(ghostIp.Ip, ghostIp.Port, packetFactory, glogger.NewGLogger(0, glogger.GetFullLogger()))
	blockContainer = store.NewBlockContainer("sqlite3")
	from, _        = net.ResolveUDPAddr("udp", ghostIp.Ip+":"+ghostIp.Port)
	owner          = gcrypto.GenerateKeyPair()

	w      = gcrypto.NewWallet("test", owner, ghostIp, &ptypes.GhostUser{PubKey: "masterpubkey", Nickname: "master"})
	master = NewMasterNode(1, w, ghostIp, packetFactory, udp, blockContainer, gAccount, NewTrieTreeMap(owner.GetPubAddress(), gAccount))
)

func TestGetVersionSq(t *testing.T) {
	sq := packets.VersionInfoSq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), nil, 0, master.localGhostIp),
	}

	sendData, err := proto.Marshal(&sq)
	if err != nil {
		log.Fatal(err)
	}

	responseInfos := master.GetGhostNetVersionSq(&p2p.RequestHeaderInfo{Header: &packets.Header{Source: master.localGhostIp, PacketData: sendData}})
	cq := &packets.VersionInfoCq{}
	if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, packets.PacketSecondType_GetGhostNetVersion, responseInfos[0].SecondType, "packet type is wrong")
	assert.Equal(t, master.ghostNetVersion, cq.Version, "packet response is wrong")
}

func TestNotifyMasterNodeSq(t *testing.T) {
	sq := &packets.MasterNodeUserInfoSq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), nil, 0, master.localGhostIp),
		User: &ptypes.GhostUser{
			Nickname: master.nickname,
			PubKey:   owner.GetPubAddress(),
		},
	}

	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}

	responseInfos := master.NotificationMasterNodeSq(&p2p.RequestHeaderInfo{Header: &packets.Header{Source: master.localGhostIp, PacketData: sendData}})
	cq := &packets.MasterNodeUserInfoCq{}
	if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, packets.PacketSecondType_NotificationMasterNode, responseInfos[0].SecondType, "packet type is wrong")
	assert.Equal(t, "test", cq.User.Nickname, "nickname is wrong")
	assert.Equal(t, owner.GetPubAddress(), cq.User.PubKey, "pubkey is wrong")
}

func TestConnectToMasterNodeSq(t *testing.T) {
	sq := &packets.MasterNodeUserInfoSq{
		Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), nil, 0, master.localGhostIp),
		User: &ptypes.GhostUser{
			Nickname: master.nickname,
			PubKey:   owner.GetPubAddress(),
		},
	}

	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}

	responseInfos := master.ConnectToMasterNodeSq(&p2p.RequestHeaderInfo{Header: &packets.Header{Source: master.localGhostIp, PacketData: sendData}})
	cq := &packets.MasterNodeUserInfoCq{}
	if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, packets.PacketSecondType_ConnectToMasterNode, responseInfos[0].SecondType, "packet type is wrong")
	assert.Equal(t, "test", cq.User.Nickname, "nickname is wrong")
	assert.Equal(t, owner.GetPubAddress(), cq.User.PubKey, "pubkey is wrong")
}

func TestRequestMasterNodeListSq(t *testing.T) {
	TestNotifyMasterNodeSq(t)

	sq := &packets.RequestMasterNodeListSq{
		Master:     p2p.MakeMasterPacket(master.owner.GetPubAddress(), nil, 0, master.localGhostIp),
		StartIndex: 0,
	}

	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}

	responseInfos := master.RequestMasterNodeListSq(&p2p.RequestHeaderInfo{Header: &packets.Header{Source: master.localGhostIp, PacketData: sendData}})
	resInfo := responseInfos[1]
	resSq := &packets.ResponseMasterNodeListSq{}
	if err := proto.Unmarshal(resInfo.PacketData, resSq); err != nil {
		log.Fatal(err)
	}
	user := resSq.User[0]
	assert.Equal(t, "test", user.Nickname, "nickname is wrong")
	assert.Equal(t, owner.GetPubAddress(), user.PubKey, "pubkey is wrong")
	assert.Equal(t, packets.PacketSecondType_RequestMasterNodeList, responseInfos[0].SecondType, "packet[0] type is wrong")
	assert.Equal(t, packets.PacketSecondType_ResponseMasterNodeList, responseInfos[1].SecondType, "packet[1] type is wrong")
}

func TestResponseMasterNodeListSq(t *testing.T) {
	sq := &packets.ResponseMasterNodeListSq{Master: p2p.MakeMasterPacket(master.owner.GetPubAddress(), nil, 0, master.localGhostIp)}

	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}

	responseInfos := master.ResponseMasterNodeListSq(&p2p.RequestHeaderInfo{Header: &packets.Header{Source: master.localGhostIp, PacketData: sendData}})
	cq := &packets.ResponseMasterNodeListCq{}
	if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, packets.PacketSecondType_ResponseMasterNodeList, responseInfos[0].SecondType, "packet type is wrong")
}

func TestSearchMasterPubKey(t *testing.T) {
	TestNotifyMasterNodeSq(t)
	sq := &packets.SearchGhostPubKeySq{
		Master:   p2p.MakeMasterPacket(master.owner.GetPubAddress(), nil, 0, master.localGhostIp),
		Nickname: "test",
		PubKey:   master.owner.GetPubAddress(),
	}

	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}

	responseInfos := master.SearchMasterPubKeySq(&p2p.RequestHeaderInfo{Header: &packets.Header{Source: master.localGhostIp, PacketData: sendData}})
	cq := &packets.SearchGhostPubKeyCq{}
	if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, packets.PacketSecondType_SearchUserInfoByPubKey, responseInfos[0].SecondType, "packet type is wrong")
	assert.Equal(t, owner.GetPubAddress(), cq.User[0].PubKey, "pubkey is wrong")
}
