package p2p

import (
	"net"
	"sync"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestUdpDefault(t *testing.T) {
	var wg sync.WaitGroup
	wg.Add(1)
	netChannel := make(chan RequestPacketInfo)
	srv := NewUdpServer("127.0.0.1", "8888", NewPacketFactory(), glogger.NewGLogger(0, glogger.GetFullLogger()))
	srv.Start(netChannel, "", "")

	go func() {
		packetInfo := <-netChannel
		packetByte := packetInfo.PacketByte
		recvPacket := packets.Header{}
		if err := proto.Unmarshal(packetByte, &recvPacket); err != nil {
			// packet type별로 callback handler를 만들어야한다.
			t.Error(err)
		}
		infoSq := &packets.MasterNodeUserInfoSq{}
		proto.Unmarshal(recvPacket.PacketData, infoSq)
		assert.Equal(t, "test", infoSq.User.Nickname, "packet내용이 맞지 않습니다.")
		wg.Done()
	}()

	testPacket := packets.MasterNodeUserInfoSq{User: &ptypes.GhostUser{Nickname: "test"}}
	data, err := proto.Marshal(&testPacket)
	if err != nil {
		t.Error(err)
	}

	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:8888")
	srv.SendResponse(&ResponseHeaderInfo{
		ToAddr: addr, PacketType: packets.PacketType_MasterNetwork,
		SecondType: packets.PacketSecondType_NotificationMasterNode,
		PacketData: data, SqFlag: true})

	wg.Wait()
}

var (
	GlobalWg   sync.WaitGroup
	TestResult bool
)

func PacketSqHandler(requestHeaderInfo *RequestHeaderInfo) []ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	packetByte := header.PacketData
	infoSq := &packets.MasterNodeUserInfoSq{}
	proto.Unmarshal(packetByte, infoSq)
	TestResult = infoSq.User.Nickname == "test"
	GlobalWg.Done()

	return nil
}

func TestPacketHandler(t *testing.T) {
	GlobalWg.Add(1)
	TestResult = false

	srv := NewUdpServer("127.0.0.1", "8889", NewPacketFactory(), glogger.NewGLogger(0, glogger.GetFullLogger()))
	srv.Start(nil, "", "")
	srv.Pf.SingleRegisterPacketHandler(packets.PacketType_MasterNetwork,
		packets.PacketSecondType_NotificationMasterNode,
		PacketSqHandler, nil)

	testPacket := packets.MasterNodeUserInfoSq{User: &ptypes.GhostUser{Nickname: "test"}}
	data, err := proto.Marshal(&testPacket)
	if err != nil {
		t.Error(err)
	}

	addr, _ := net.ResolveUDPAddr("udp", "127.0.0.1:8889")
	srv.SendResponse(&ResponseHeaderInfo{
		ToAddr: addr, PacketType: packets.PacketType_MasterNetwork,
		SecondType: packets.PacketSecondType_NotificationMasterNode,
		PacketData: data, SqFlag: true})
	GlobalWg.Wait()
	assert.Equal(t, true, TestResult, "packet내용이 맞지 않습니다.")
}

// buffer size overflow => packet header + buf(1024) = 1024?
