package p2p

import (
	"net"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/rs/xid"
)

type RoutingInfo struct {
	RoutingType packets.RoutingType
	SourceIp    *net.UDPAddr
	Level       int
	Context     interface{}
}

type FuncPacketHandler func(*RequestHeaderInfo) []ResponseHeaderInfo

type PacketSecondHandler struct {
	packetSqHandler map[packets.PacketSecondType]FuncPacketHandler
	packetCqHandler map[packets.PacketSecondType]FuncPacketHandler
}

type PacketFactory struct {
	firstLevel map[packets.PacketType]*PacketSecondHandler
}

func NewPacketFactory() *PacketFactory {
	return &PacketFactory{
		firstLevel: make(map[packets.PacketType]*PacketSecondHandler),
	}
}

func MakeMasterPacket(from string, reqId []byte, clientId uint32, fromIp *ptypes.GhostIp) *packets.MasterPacket {
	if reqId == nil {
		reqId = xid.New().Bytes()
	}
	return &packets.MasterPacket{
		Common: &packets.GhostPacket{
			FromPubKeyAddress: from,
			RequestId:         reqId,
			ClientId:          clientId,
			TimeId:            uint64(time.Now().Unix()),
		},
		RoutingT: packets.RoutingType_None,
		Level:    1,
	}
}

func (pf *PacketFactory) SingleRegisterPacketHandler(firstType packets.PacketType, packetType packets.PacketSecondType,
	sqHandler FuncPacketHandler, cqHandler FuncPacketHandler) {
	pf.firstLevel[firstType] = &PacketSecondHandler{
		packetSqHandler: make(map[packets.PacketSecondType]FuncPacketHandler),
		packetCqHandler: make(map[packets.PacketSecondType]FuncPacketHandler),
	}
	pf.firstLevel[firstType].packetSqHandler[packetType] = sqHandler
	pf.firstLevel[firstType].packetCqHandler[packetType] = cqHandler
}

func (pf *PacketFactory) RegisterPacketHandler(firstType packets.PacketType,
	sqHandler map[packets.PacketSecondType]FuncPacketHandler,
	cqHandler map[packets.PacketSecondType]FuncPacketHandler) {
	pf.firstLevel[firstType] = &PacketSecondHandler{
		packetSqHandler: sqHandler,
		packetCqHandler: cqHandler,
	}
}

func (pf *PacketFactory) UpdatePacketHandler(firstType packets.PacketType,
	sqHandler map[packets.PacketSecondType]FuncPacketHandler,
	cqHandler map[packets.PacketSecondType]FuncPacketHandler) {
	if pf.firstLevel[firstType] == nil {
		return
	}
	for key, value := range sqHandler {
		pf.firstLevel[firstType].packetSqHandler[key] = value
	}
	for key, value := range cqHandler {
		pf.firstLevel[firstType].packetCqHandler[key] = value
	}
}
