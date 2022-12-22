package mainserver

import (
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blocks"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"google.golang.org/protobuf/proto"
)

type MainServer struct {
	con   *consensus.Consensus
	fsm   *consensus.BlockMachine
	block *blocks.Blocks
	user  *gcrypto.GhostAddress
	udp   *p2p.UdpServer
}

func NewMainServer(con *consensus.Consensus,
	fsm *consensus.BlockMachine,
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
	netChannel := make(chan []byte)

	main.udp.Start(netChannel)

	for {
		select {
		case packetByte := <-netChannel:
			recvPacket := packets.Any{}
			_ := proto.Unmarshal(packetByte, &recvPacket)
		case <-time.After(1000 * time.Nanosecond):

		}
	}
}
