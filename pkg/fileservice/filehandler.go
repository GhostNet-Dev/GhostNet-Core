package fileservice

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
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
		return []p2p.ResponseHeaderInfo{*fileService.makeFileInfo(sq, header.Source.GetUdpAddr())}
	case packets.FileRequestType_GetFileData:
		fileObj, exist := fileService.fileObjManager.GetFileObject(sq.Filename)
		if !exist {
			return nil
		}

		cq := &packets.RequestFilePacketCq{
			Master:      p2p.MakeMasterPacket(fileService.owner.GetPubAddress(), sq.Master.GetRequestId(), 0, fileService.localAddr),
			RequestType: sq.RequestType,
			Filename:    sq.Filename,
			StartOffset: sq.StartOffset,
			FileLength:  fileObj.FileLength,
			Result:      exist,
		}

		fileService.FileLogger(cq)
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
				RequestId:  cq.Master.GetRequestId(),
				SqFlag:     false,
			},
			*fileService.sendFileData(sq.Filename, sq.StartOffset, 0, sq.Master.GetTimeId(), header.Source.GetUdpAddr()),
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

	if cq.Result && cq.RequestType == packets.FileRequestType_GetFileInfo {
		if fileObj := fileService.fileObjManager.AllocBuffer(cq.Filename, cq.FileLength); fileObj == nil {
			fileService.glog.DebugOutput(fileService, "wrong filename", glogger.Default)
			return nil
		}
		sq := &packets.RequestFilePacketSq{
			Master:      p2p.MakeMasterPacket(fileService.owner.GetPubAddress(), nil, 0, fileService.localAddr),
			RequestType: packets.FileRequestType_GetFileData,
			Filename:    cq.Filename,
			StartOffset: 0,
		}

		fileService.FileLogger(sq)
		sendData, err := proto.Marshal(sq)
		if err != nil {
			log.Fatal(err)
		}

		return []p2p.ResponseHeaderInfo{
			{
				ToAddr:     header.Source.GetUdpAddr(),
				PacketType: packets.PacketType_FileTransfer,
				SecondType: packets.PacketSecondType_RequestFile,
				RequestId:  sq.Master.GetRequestId(),
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
		Master: p2p.MakeMasterPacket(fileService.owner.GetPubAddress(), sq.Master.GetRequestId(), 0, fileService.localAddr),
		Result: true,
	}

	fileService.FileLogger(cq)
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
			RequestId:  cq.Master.GetRequestId(),
			SqFlag:     false,
		},
	}

	if fileService.saveToFileObject(sq.Filename, sq.StartPos, sq.BufferSize, sq.FileData, sq.FileLength) == false {
		newSq := &packets.RequestFilePacketSq{
			Master:      p2p.MakeMasterPacket(fileService.owner.GetPubAddress(), nil, 0, fileService.localAddr),
			RequestType: packets.FileRequestType_GetFileData,
			Filename:    sq.Filename,
			StartOffset: sq.StartPos + Buffer_Size,
		}
		fileService.FileLogger(newSq)
		newSqData, err := proto.Marshal(newSq)
		if err != nil {
			log.Fatal(err)
		}
		headerInfo = append(headerInfo, p2p.ResponseHeaderInfo{
			ToAddr:     header.Source.GetUdpAddr(),
			PacketType: packets.PacketType_FileTransfer,
			SecondType: packets.PacketSecondType_RequestFile,
			RequestId:  newSq.Master.GetRequestId(),
			PacketData: newSqData,
			SqFlag:     true,
		})
	}

	return headerInfo
}

func (fileService *FileService) ResponseFileCq(requestHeaderInfo *p2p.RequestHeaderInfo) []p2p.ResponseHeaderInfo {
	return nil
}
