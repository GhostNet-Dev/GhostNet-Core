package fileserver

import (
	"log"
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"google.golang.org/protobuf/proto"
)

func (fileServer *FileServer) RequestFileSq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
	sq := &packets.RequestFilePacketSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	switch sq.RequestType {
	case packets.FileRequestType_GetFileInfo:
		return []p2p.PacketHeaderInfo{*fileServer.MakeFileInfo(sq.Filename)}
	case packets.FileRequestType_GetFileData:
		cq := &packets.RequestFilePacketCq{
			Master:      p2p.MakeMasterPacket(fileServer.owner.GetPubAddress(), 0, 0, fileServer.localAddr),
			RequestType: sq.RequestType,
			Filename:    sq.Filename,
			StartOffset: 0,
		}

		sendData, err := proto.Marshal(cq)
		if err != nil {
			log.Fatal(err)
		}
		return []p2p.PacketHeaderInfo{
			{
				ToAddr:     from,
				PacketType: packets.PacketType_FileTransfer,
				SecondType: packets.PacketSecondType_RequestFile,
				PacketData: sendData,
				SqFlag:     false,
			},
			*fileServer.SendFileData(sq.Filename, sq.StartOffset, 0, sq.Master.Common.TimeId),
		}
	}
	return nil
}

func (fileServer *FileServer) RequestFileCq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
	cq := &packets.RequestFilePacketCq{}
	if err := proto.Unmarshal(header.PacketData, cq); err != nil {
		log.Fatal(err)
	}

	if cq.Result == true && cq.RequestType == packets.FileRequestType_GetFileInfo {
		fileServer.fileObjManager.CreateFileObj(cq.Filename, nil, cq.FileLength, nil, nil)
		sq := &packets.RequestFilePacketSq{
			Master:      p2p.MakeMasterPacket(fileServer.owner.GetPubAddress(), 0, 0, fileServer.localAddr),
			RequestType: packets.FileRequestType_GetFileData,
			Filename:    cq.Filename,
			StartOffset: 0,
		}

		sendData, err := proto.Marshal(sq)
		if err != nil {
			log.Fatal(err)
		}

		return []p2p.PacketHeaderInfo{
			{
				ToAddr:     from,
				PacketType: packets.PacketType_FileTransfer,
				SecondType: packets.PacketSecondType_RequestFile,
				PacketData: sendData,
				SqFlag:     false,
			},
		}
	}
	return nil
}

func (fileServer *FileServer) ResponseFileSq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
	sq := &packets.ResponseFilePacketSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	cq := &packets.ResponseFilePacketCq{
		Master: p2p.MakeMasterPacket(fileServer.owner.GetPubAddress(), 0, 0, fileServer.localAddr),
		Result: true,
	}

	sendData, err := proto.Marshal(cq)
	if err != nil {
		log.Fatal(err)
	}
	headerInfo := []p2p.PacketHeaderInfo{
		{
			ToAddr:     from,
			PacketType: packets.PacketType_FileTransfer,
			SecondType: packets.PacketSecondType_ResponseFile,
			PacketData: sendData,
			SqFlag:     false,
		},
	}

	if fileServer.SaveToFileObject(sq.Filename, sq.StartPos, sq.BufferSize, sq.FileData, sq.FileLength) == false {
		newSq := &packets.RequestFilePacketSq{
			Master:      p2p.MakeMasterPacket(fileServer.owner.GetPubAddress(), 0, 0, fileServer.localAddr),
			RequestType: packets.FileRequestType_GetFileData,
			Filename:    sq.Filename,
			StartOffset: sq.StartPos + BufferSize,
		}
		newSqData, err := proto.Marshal(newSq)
		if err != nil {
			log.Fatal(err)
		}
		headerInfo = append(headerInfo, p2p.PacketHeaderInfo{
			ToAddr:     from,
			PacketType: packets.PacketType_FileTransfer,
			SecondType: packets.PacketSecondType_RequestFile,
			PacketData: newSqData,
			SqFlag:     true,
		})
	}

	return headerInfo
}

func (fileServer *FileServer) ResponseFileCq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
	return nil
}
