package blockmanager

import (
	"log"
	"net"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

var (
	from, _ = net.ResolveUDPAddr("udp", ghostIp.Ip+":"+ghostIp.Port)
)

func TestGetHeightestBlock(t *testing.T) {
	sq := &packets.MasterNodeUserInfoSq{
		Master: p2p.MakeMasterPacket(blockServer.owner.GetPubAddress(), 0, 0, blockServer.localIpAddr),
		User: &ptypes.GhostUser{
			Nickname: nickname,
			PubKey:   blockServer.owner.GetPubAddress(),
		},
	}

	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}

	responseInfos := blockServer.GetHeightestBlockSq(&packets.Header{PacketData: sendData}, from)
	cq := &packets.MasterNodeUserInfoCq{}
	if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, packets.PacketSecondType_ConnectToMasterNode, responseInfos[0].SecondType, "packet type is wrong")
	assert.Equal(t, "test", cq.User.Nickname, "nickname is wrong")
	assert.Equal(t, Miner.GetPubAddress(), cq.User.PubKey, "pubkey is wrong")
}
