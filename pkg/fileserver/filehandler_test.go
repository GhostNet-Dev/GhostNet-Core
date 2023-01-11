package fileserver

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

	responseInfos := fileServer.RequestFileSq(&packets.Header{PacketData: sendData}, from)
	cq := &packets.RequestFilePacketCq{}
	if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
		log.Fatal(err)
	}
	info, err := os.Stat(fileServer.localFilePath + testfile)

	assert.Equal(t, packets.PacketSecondType_RequestFile, responseInfos[0].SecondType, "packet type is wrong")
	assert.Equal(t, packets.FileRequestType_GetFileInfo, cq.RequestType, "request type is wrong")
	assert.Equal(t, uint64(info.Size()), cq.FileLength, "file length is wrong")
	assert.Equal(t, true, cq.Result, "result is wrong")

}

func TestRequestFileDataSq(t *testing.T) {
	testFileInit()
	// todo for total file packet
	sq := packets.RequestFilePacketSq{
		Master:      p2p.MakeMasterPacket(owner.GetPubAddress(), 0, 0, ipAddr),
		RequestType: packets.FileRequestType_GetFileData,
		Filename:    testfile,
	}
	sendData, err := proto.Marshal(&sq)
	if err != nil {
		log.Fatal(err)
	}
	responseInfos := fileServer.RequestFileSq(&packets.Header{PacketData: sendData}, from)
	cq := &packets.RequestFilePacketCq{}
	if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
		log.Fatal(err)
	}
	newSq := &packets.ResponseFilePacketSq{}
	if err := proto.Unmarshal(responseInfos[1].PacketData, newSq); err != nil {
		log.Fatal(err)
	}
	info, err := os.Stat(fileServer.localFilePath + testfile)
	assert.Equal(t, packets.PacketSecondType_RequestFile, responseInfos[0].SecondType, "packet0 type is wrong")
	assert.Equal(t, packets.PacketSecondType_ResponseFile, responseInfos[1].SecondType, "packet1 type is wrong")
	assert.Equal(t, packets.FileRequestType_GetFileData, cq.RequestType, "request type is wrong")
	assert.Equal(t, uint64(info.Size()), cq.FileLength, "file length is wrong")
	assert.Equal(t, true, cq.Result, "result is wrong")

}

func TestResponseSqTest(t *testing.T) {
	testFileInit()
	// todo: fileread
	sq := packets.ResponseFilePacketSq{
		Master:   p2p.MakeMasterPacket(owner.GetPubAddress(), 0, 0, ipAddr),
		Filename: testfile,
	}
	sendData, err := proto.Marshal(&sq)
	if err != nil {
		log.Fatal(err)
	}
	responseInfos := fileServer.ResponseFileSq(&packets.Header{PacketData: sendData}, from)
	cq := &packets.ResponseFilePacketCq{}
	if err := proto.Unmarshal(responseInfos[0].PacketData, cq); err != nil {
		log.Fatal(err)
	}
}
