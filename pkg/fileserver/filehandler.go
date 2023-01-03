package fileserver

import (
	"log"
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"google.golang.org/protobuf/proto"
)

func (fileServer *FileServer) RequestFileSq(header *packets.Header, from *net.UDPAddr) []p2p.ResponsePacketInfo {
	sq := &packets.RequestFilePacketSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	cq := &packets.RequestFilePacketCq{
		Master: p2p.MakeMasterPacket(fileServer.owner.GetPubAddress(), 0, 0, fileServer.localAddr),
	}

	sendData, err := proto.Marshal(cq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.ResponsePacketInfo{
		{
			ToAddr:     from,
			SecondType: packets.PacketSecondType_RequestFile,
			PacketData: sendData,
			SqFlag:     false,
		},
	}
}

func (fileServer *FileServer) RequestFileCq(header *packets.Header, from *net.UDPAddr) {
	cq := &packets.RequestFilePacketCq{}
	if err := proto.Unmarshal(header.PacketData, cq); err != nil {
		log.Fatal(err)
	}

	if cq.Result == true && cq.RequestType == packets.FileRequestType_GetFileInfo {

	}
}

func (fileServer *FileServer) ResponseFileSq(header *packets.Header, from *net.UDPAddr) []p2p.ResponsePacketInfo {

	return []p2p.ResponsePacketInfo{
		{
			ToAddr:     from,
			SecondType: packets.PacketSecondType_ResponseFile,
			SqFlag:     false,
		},
	}
}

func (fileServer *FileServer) ResponseFileCq(header *packets.Header, from *net.UDPAddr) {

}
