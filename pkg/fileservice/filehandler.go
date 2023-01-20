package fileservice

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"google.golang.org/protobuf/proto"
)

func (fileService *FileService) RequestFileSq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	sq := &packets.RequestFilePacketSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	switch sq.RequestType {
	case packets.FileRequestType_GetFileInfo:
		return []p2p.PacketHeaderInfo{*fileService.makeFileInfo(sq.Filename)}
	case packets.FileRequestType_GetFileData:
		fileObj, exist := fileService.fileObjManager.GetFileObject(sq.Filename)
		if exist == false {
			return nil
		}

		cq := &packets.RequestFilePacketCq{
			Master:      p2p.MakeMasterPacket(fileService.owner.GetPubAddress(), 0, 0, fileService.localAddr),
			RequestType: sq.RequestType,
			Filename:    sq.Filename,
			StartOffset: sq.StartOffset,
			FileLength:  fileObj.FileLength,
			Result:      exist,
		}

		sendData, err := proto.Marshal(cq)
		if err != nil {
			log.Fatal(err)
		}
		return []p2p.PacketHeaderInfo{
			{
				ToAddr:     routingInfo.SourceIp,
				PacketType: packets.PacketType_FileTransfer,
				SecondType: packets.PacketSecondType_RequestFile,
				PacketData: sendData,
				SqFlag:     false,
			},
			*fileService.sendFileData(sq.Filename, sq.StartOffset, 0, sq.Master.Common.TimeId),
		}
	}
	return nil
}

func (fileService *FileService) RequestFileCq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	cq := &packets.RequestFilePacketCq{}
	if err := proto.Unmarshal(header.PacketData, cq); err != nil {
		log.Fatal(err)
	}

	if cq.Result == true && cq.RequestType == packets.FileRequestType_GetFileInfo {
		fileService.fileObjManager.CreateFileObj(cq.Filename, nil, cq.FileLength, nil, nil)
		sq := &packets.RequestFilePacketSq{
			Master:      p2p.MakeMasterPacket(fileService.owner.GetPubAddress(), 0, 0, fileService.localAddr),
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
				ToAddr:     routingInfo.SourceIp,
				PacketType: packets.PacketType_FileTransfer,
				SecondType: packets.PacketSecondType_RequestFile,
				PacketData: sendData,
				SqFlag:     false,
			},
		}
	}
	return nil
}

// ResponseFileSq
func (fileService *FileService) ResponseFileSq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	sq := &packets.ResponseFilePacketSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	cq := &packets.ResponseFilePacketCq{
		Master: p2p.MakeMasterPacket(fileService.owner.GetPubAddress(), 0, 0, fileService.localAddr),
		Result: true,
	}

	sendData, err := proto.Marshal(cq)
	if err != nil {
		log.Fatal(err)
	}
	headerInfo := []p2p.PacketHeaderInfo{
		{
			ToAddr:     routingInfo.SourceIp,
			PacketType: packets.PacketType_FileTransfer,
			SecondType: packets.PacketSecondType_ResponseFile,
			PacketData: sendData,
			SqFlag:     false,
		},
	}

	if fileService.saveToFileObject(sq.Filename, sq.StartPos, sq.BufferSize, sq.FileData, sq.FileLength) == false {
		newSq := &packets.RequestFilePacketSq{
			Master:      p2p.MakeMasterPacket(fileService.owner.GetPubAddress(), 0, 0, fileService.localAddr),
			RequestType: packets.FileRequestType_GetFileData,
			Filename:    sq.Filename,
			StartOffset: sq.StartPos + BufferSize,
		}
		newSqData, err := proto.Marshal(newSq)
		if err != nil {
			log.Fatal(err)
		}
		headerInfo = append(headerInfo, p2p.PacketHeaderInfo{
			ToAddr:     routingInfo.SourceIp,
			PacketType: packets.PacketType_FileTransfer,
			SecondType: packets.PacketSecondType_RequestFile,
			PacketData: newSqData,
			SqFlag:     true,
		})
	}

	return headerInfo
}

func (fileService *FileService) ResponseFileCq(header *packets.Header, routingInfo *p2p.RoutingInfo) []p2p.PacketHeaderInfo {
	return nil
}