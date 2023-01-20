package mainserver

import (
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blocks"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus/states"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"google.golang.org/protobuf/proto"
)

type MainServer struct {
	con   *consensus.Consensus
	fsm   *states.BlockMachine
	block *blocks.Blocks
	user  *gcrypto.GhostAddress
	udp   *p2p.UdpServer
}

func NewMainServer(con *consensus.Consensus,
	fsm *states.BlockMachine,
	block *blocks.Blocks,
	user *gcrypto.GhostAddress,
	udp *p2p.UdpServer) *MainServer {

	return &MainServer{
		con:   con,
		fsm:   fsm,
		block: block,
		user:  user,
		udp:   udp,
	}
}

func (main *MainServer) StartServer() {
	netChannel := make(chan p2p.RequestPacketInfo)

	main.udp.Start(netChannel)

	for {
		select {
		case packetInfo := <-netChannel:
			packetByte := packetInfo.PacketByte
			recvPacket := packets.Header{}
			if err := proto.Unmarshal(packetByte, &recvPacket); err != nil {
				// masternode layer로 전송한다.

			}
		case <-time.After(1000 * time.Nanosecond):

		}
	}
}
