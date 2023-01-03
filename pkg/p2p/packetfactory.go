package p2p

import (
	"net"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
)

type PacketSecondHandler struct {
	packetSqHandler map[packets.PacketSecondType]func(*packets.Header, *net.UDPAddr) []ResponsePacketInfo
	packetCqHandler map[packets.PacketSecondType]func(*packets.Header, *net.UDPAddr)
}

type PacketFactory struct {
	firstLevel map[packets.PacketType]*PacketSecondHandler
}

func NewPacketFactory() *PacketFactory {
	return &PacketFactory{
		firstLevel: make(map[packets.PacketType]*PacketSecondHandler),
	}
}

func MakeMasterPacket(from string, reqId uint32, clientId uint32, fromIp *ptypes.GhostIp) *packets.MasterPacket {
	return &packets.MasterPacket{
		Common: &packets.GhostPacket{
			FromPubKeyAddress: from,
			RequestId:         reqId,
			ClientId:          clientId,
			TimeId:            uint64(time.Now().Unix()),
		},
		RoutingT: packets.RoutingType_None,
		Level:    0,
	}
}

func (pf *PacketFactory) SingleRegisterPacketHandler(firstType packets.PacketType, packetType packets.PacketSecondType,
	sqHandler func(*packets.Header, *net.UDPAddr) []ResponsePacketInfo,
	cqHandler func(*packets.Header, *net.UDPAddr)) {
	pf.firstLevel[firstType] = &PacketSecondHandler{
		packetSqHandler: make(map[packets.PacketSecondType]func(*packets.Header, *net.UDPAddr) []ResponsePacketInfo),
		packetCqHandler: make(map[packets.PacketSecondType]func(*packets.Header, *net.UDPAddr)),
	}
	pf.firstLevel[firstType].packetSqHandler[packetType] = sqHandler
	pf.firstLevel[firstType].packetCqHandler[packetType] = cqHandler
}

func (pf *PacketFactory) RegisterPacketHandler(firstType packets.PacketType,
	sqHandler map[packets.PacketSecondType]func(*packets.Header, *net.UDPAddr) []ResponsePacketInfo,
	cqHandler map[packets.PacketSecondType]func(*packets.Header, *net.UDPAddr)) {
	pf.firstLevel[firstType] = &PacketSecondHandler{
		packetSqHandler: sqHandler,
		packetCqHandler: cqHandler,
	}
}
