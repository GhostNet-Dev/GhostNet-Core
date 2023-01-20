package fileservice

import (
	"log"
	"os"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestRequestFileSq(t *testing.T) {
	testFileInit()
	sq := packets.RequestFilePacketSq{
		Master:      p2p.MakeMasterPacket(owner.GetPubAddress(), 0, 0, ipAddr),
		RequestType: packets.FileRequestType_GetFileInfo,
		Filename:    testfile,
	}

	sendData, err := proto.Marshal(&sq)
	if err != nil {
		log.Fatal(err)
	}

	responseInfos := fileService.RequestFileSq(&packets.Header{PacketData: sendData}, &p2p.RoutingInfo{SourceIp: from})
	cq := &packets.RequestFilePacketCq{}
	if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
		log.Fatal(err)
	}
	info, err := os.Stat(fileService.localFilePath + testfile)

	assert.Equal(t, packets.PacketSecondType_RequestFile, responseInfos[0].SecondType, "packet type is wrong")
	assert.Equal(t, packets.FileRequestType_GetFileInfo, cq.RequestType, "request type is wrong")
	assert.Equal(t, uint64(info.Size()), cq.FileLength, "file length is wrong")
	assert.Equal(t, true, cq.Result, "result is wrong")

}

func TestRequestFileDataSq(t *testing.T) {
	testFileInit()
	fileService.LoadFileToMemory(testfile)
	info, _ := os.Stat(fileService.localFilePath + testfile)
	totalDownloadSize := uint32(0)
	// todo for total file packet
	for offset := uint64(0); offset < uint64(info.Size()); offset += BufferSize {
		sq := packets.RequestFilePacketSq{
			Master:      p2p.MakeMasterPacket(owner.GetPubAddress(), 0, 0, ipAddr),
			RequestType: packets.FileRequestType_GetFileData,
			Filename:    testfile,
			StartOffset: offset,
		}
		sendData, err := proto.Marshal(&sq)
		if err != nil {
			log.Fatal(err)
		}
		responseInfos := fileService.RequestFileSq(&packets.Header{PacketData: sendData}, &p2p.RoutingInfo{SourceIp: from})
		assert.Equal(t, true, responseInfos != nil, "return is nil")

		cq := &packets.RequestFilePacketCq{}
		if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
			log.Fatal(err)
		}
		newSq := &packets.ResponseFilePacketSq{}
		if err := proto.Unmarshal(responseInfos[1].PacketData, newSq); err != nil {
			log.Fatal(err)
		}
		totalDownloadSize += newSq.BufferSize
		assert.Equal(t, packets.PacketSecondType_RequestFile, responseInfos[0].SecondType, "packet0 type is wrong")
		assert.Equal(t, packets.PacketSecondType_ResponseFile, responseInfos[1].SecondType, "packet1 type is wrong")
		assert.Equal(t, packets.FileRequestType_GetFileData, cq.RequestType, "request type is wrong")
		assert.Equal(t, uint64(info.Size()), cq.FileLength, "file length is wrong")
		assert.Equal(t, true, cq.Result, "result is wrong")
	}
	assert.Equal(t, uint32(info.Size()), totalDownloadSize, "total download size is wrong")
}

func TestResponseSqTest(t *testing.T) {
	testFileInit()
	// todo: fileread
	buf := make([]byte, BufferSize)
	fp, err := os.Open(testfile)
	if err != nil {
		log.Fatal(err)
	}
	defer fp.Close()
	fileInfo, _ := fp.Stat()
	for {
		readSize, _ := fp.Read(buf)
		if readSize == 0 {
			break
		}
		offset, _ := fp.Seek(0, os.SEEK_CUR)
		sq := packets.ResponseFilePacketSq{
			Master:     p2p.MakeMasterPacket(owner.GetPubAddress(), 0, 0, ipAddr),
			Filename:   testfile,
			FileData:   buf[:readSize],
			FileLength: uint64(fileInfo.Size()),
			BufferSize: uint32(readSize),
			StartPos:   uint64(offset - int64(readSize)),
		}
		sendData, err := proto.Marshal(&sq)
		if err != nil {
			log.Fatal(err)
		}
		responseInfos := fileService.ResponseFileSq(&packets.Header{PacketData: sendData}, &p2p.RoutingInfo{SourceIp: from})
		cq := &packets.ResponseFilePacketCq{}
		if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
			log.Fatal(err)
		}
		if len(responseInfos) > 1 {
			newSq := &packets.RequestFilePacketSq{}
			if err := proto.Unmarshal(responseInfos[1].PacketData, newSq); err != nil {
				log.Fatal(err)
			}
			assert.Equal(t, packets.PacketSecondType_RequestFile, responseInfos[1].SecondType, "packet1 type is wrong")
			assert.Equal(t, packets.FileRequestType_GetFileData, newSq.RequestType, "request type is wrong")
			assert.Equal(t, uint64(offset), newSq.StartOffset, "file offset is wrong")
		}
		assert.Equal(t, packets.PacketSecondType_ResponseFile, responseInfos[0].SecondType, "packet0 type is wrong")
		assert.Equal(t, true, cq.Result, "result is wrong")
	}
	fileObj, exist := fileService.fileObjManager.GetFileObject(testfile)
	assert.Equal(t, true, exist, "file is not exist")
	assert.Equal(t, true, fileObj.CompleteDone, "file is not completed")

}