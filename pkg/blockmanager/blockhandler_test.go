package blockmanager

import (
	"log"
	"net"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

var (
	from, _ = net.ResolveUDPAddr("udp", ghostIp.Ip+":"+ghostIp.Port)
)

func TestGetHeightestBlock(t *testing.T) {
	blockContainer.BlockContainerOpen("../../db.sqlite3.sql", "./")
	defer blockContainer.Close()
	sq := &packets.GetHeightestBlockSq{
		Master: p2p.MakeMasterPacket(blockServer.owner.GetPubAddress(), nil, 0, blockServer.localIpAddr),
	}

	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}

	responseInfos := blockServer.GetHeightestBlockSq(&packets.Header{PacketData: sendData}, from)
	cq := &packets.GetHeightestBlockCq{}
	if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, packets.PacketThirdType_GetHeightestBlock, responseInfos[0].ThirdType, "packet type is wrong")
}
