package fileservice

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"google.golang.org/protobuf/proto"
)

func (fileService *FileService) RequestFileSq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
	sq := &packets.RequestFilePacketSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	switch sq.RequestType {
	case packets.FileRequestType_GetFileInfo:
		return []p2p.ResponseHeaderInfo{*fileService.makeFileInfo(sq.Filename, header.Source.GetUdpAddr())}
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
		return []p2p.ResponseHeaderInfo{
			{
				ToAddr:     header.Source.GetUdpAddr(),
				PacketType: packets.PacketType_FileTransfer,
				SecondType: packets.PacketSecondType_RequestFile,
				PacketData: sendData,
				SqFlag:     false,
			},
			*fileService.sendFileData(sq.Filename, sq.StartOffset, 0, sq.Master.Common.TimeId, header.Source.GetUdpAddr()),
		}
	}
	return nil
}

func (fileService *FileService) RequestFileCq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
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

		return []p2p.ResponseHeaderInfo{
			{
				ToAddr:     header.Source.GetUdpAddr(),
				PacketType: packets.PacketType_FileTransfer,
				SecondType: packets.PacketSecondType_RequestFile,
				PacketData: sendData,
				SqFlag:     true,
			},
		}
	}
	return nil
}

// ResponseFileSq
func (fileService *FileService) ResponseFileSq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	header := requestHeaderInfo.Header
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
	headerInfo := []p2p.ResponseHeaderInfo{
		{
			ToAddr:     header.Source.GetUdpAddr(),
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
			StartOffset: sq.StartPos + Buffer_Size,
		}
		newSqData, err := proto.Marshal(newSq)
		if err != nil {
			log.Fatal(err)
		}
		headerInfo = append(headerInfo, p2p.ResponseHeaderInfo{
			ToAddr:     header.Source.GetUdpAddr(),
			PacketType: packets.PacketType_FileTransfer,
			SecondType: packets.PacketSecondType_RequestFile,
			PacketData: newSqData,
			SqFlag:     true,
		})
	}

	return headerInfo
}

func (fileService *FileService) ResponseFileCq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	return nil
}
