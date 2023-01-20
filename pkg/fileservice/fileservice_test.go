package fileservice

import (
	"log"
	"net"
	"os"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

var (
	testfile = "test.txt"
	ipAddr   = &ptypes.GhostIp{
		Ip:   "127.0.0.1",
		Port: "8888",
	}
	from, _       = net.ResolveUDPAddr("udp", ipAddr.Ip+":"+ipAddr.Port)
	udp           = p2p.NewUdpServer(ipAddr.Ip, ipAddr.Port)
	packetFactory = p2p.NewPacketFactory()
	owner         = gcrypto.GenerateKeyPair()
	fileService    = NewFileServer(udp, packetFactory, owner, ipAddr, "./")
)

func testFileInit() {
	if _, err := os.Stat(fileService.localFilePath + testfile); os.IsNotExist(err) {
		if fp, err := os.Create(fileService.localFilePath + testfile); err == nil {
			defer fp.Close()
			fp.Write(make([]byte, BufferSize*16-128))
		} else {
			log.Fatal(err)
		}
	}
}

func TestLoadFile(t *testing.T) {
	// create file -> send/recv -> defer delete file
	testFileInit()
	exist := fileService.CheckFileExist(testfile)
	assert.Equal(t, true, exist, "fail to create test file")
	fileObj := fileService.LoadFileToMemory(testfile)
	assert.Equal(t, testfile, fileObj.Filename, "wrong file")
	assert.Equal(t, true, fileObj.CompleteDone, "load fail")
}

func TestFileInfoCq(t *testing.T) {
	testFileInit()
	info, err := os.Stat(fileService.localFilePath + testfile)
	fileService.LoadFileToMemory(testfile)
	assert.Equal(t, nil, err, "errrrr")
	header := fileService.makeFileInfo(testfile)
	assert.Equal(t, packets.PacketSecondType_RequestFile, header.SecondType, "wrong second type")
	assert.Equal(t, false, header.SqFlag, "wrong SqFlag")

	cq := &packets.RequestFilePacketCq{}
	if err := proto.Unmarshal(header.PacketData, cq); err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, testfile, cq.Filename, "wrong filename")
	assert.Equal(t, uint64(info.Size()), cq.FileLength, "wrong filesize")
}

func TestSendFileData(t *testing.T) {
	testFileInit()
	fileObj := fileService.LoadFileToMemory(testfile)

	for offset := uint64(0); offset < fileObj.FileLength; offset += 1024 {
		header := fileService.sendFileData(testfile, offset, 0, 0)
		sq := &packets.ResponseFilePacketSq{}
		if err := proto.Unmarshal(header.PacketData, sq); err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, testfile, sq.Filename, "wrong filename")
		assert.Equal(t, fileObj.FileLength, sq.FileLength, "wrong filelength")
		assert.Equal(t, offset, sq.StartPos, "wrong offset")
	}
}

func TestSaveToFileSystem(t *testing.T) {
	testFileInit()
	fileInfo, _ := os.Stat(fileService.localFilePath + testfile)
	fileSize := uint64(fileInfo.Size())
	bufferSize := uint64(BufferSize)
	buf := make([]byte, bufferSize)
	for offset := uint64(0); offset < fileSize; offset += bufferSize {
		if fileSize-offset < bufferSize {
			bufferSize = fileSize - offset
		}

		done := fileService.saveToFileObject(testfile, offset, uint32(bufferSize), buf, fileSize)
		if offset+BufferSize >= fileSize {
			assert.Equal(t, true, done, "wrong complete")
		} else {
			assert.Equal(t, false, done, "wrong saving")
		}
		bufferSize = BufferSize
	}
}
