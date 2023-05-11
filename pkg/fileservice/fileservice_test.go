package fileservice

import (
	"fmt"
	"log"
	"math/rand"
	"net"
	"os"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
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
	toIpAddr = &ptypes.GhostIp{
		Ip:   "127.0.0.2",
		Port: "8888",
	}
	glog          = glogger.NewGLogger(0)
	from, _       = net.ResolveUDPAddr("udp", ipAddr.Ip+":"+ipAddr.Port)
	to, _         = net.ResolveUDPAddr("udp", toIpAddr.Ip+":"+toIpAddr.Port)
	packetFactory = p2p.NewPacketFactory()
	udp           = p2p.NewUdpServer(ipAddr.Ip, ipAddr.Port, packetFactory, glogger.NewGLogger(0))
	owner         = gcrypto.GenerateKeyPair()
	fileService   = NewFileServer(udp, packetFactory, owner, ipAddr, "./", glog)
)

func testFileInit(filename string) {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		if fp, err := os.Create(filename); err == nil {
			defer fp.Close()
			fp.Write(make([]byte, Buffer_Size*16-128))
		} else {
			log.Fatal(err)
		}
	}
}

func TestLoadFile(t *testing.T) {
	// create file -> send/recv -> defer delete file
	testFileInit(fileService.localFilePath + testfile)
	exist := fileService.CheckFileExist(testfile)
	assert.Equal(t, true, exist, "fail to create test file")
	fileObj := fileService.loadFileToMemory(testfile)
	assert.Equal(t, testfile, fileObj.Filename, "wrong file")
	assert.Equal(t, true, fileObj.CompleteDone, "load fail")
}

func TestFileInfoCq(t *testing.T) {
	testFileInit(fileService.localFilePath + testfile)
	info, _ := os.Stat(fileService.localFilePath + testfile)
	assert.Equal(t, true, info.Size() > 0, "get file stat error")
	fileService.loadFileToMemory(testfile)
	header := fileService.makeFileInfo(testfile, toIpAddr.GetUdpAddr())
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
	testFileInit(fileService.localFilePath + testfile)
	fileObj := fileService.loadFileToMemory(testfile)

	for offset := uint64(0); offset < fileObj.FileLength; offset += 1024 {
		header := fileService.sendFileData(testfile, offset, 0, 0, toIpAddr.GetUdpAddr())
		sq := &packets.ResponseFilePacketSq{}
		if err := proto.Unmarshal(header.PacketData, sq); err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, testfile, sq.Filename, "wrong filename")
		assert.Equal(t, fileObj.FileLength, sq.FileLength, "wrong filelength")
		assert.Equal(t, offset, sq.StartPos, "wrong offset")
	}
}

func TestSaveToExistFileSystem(t *testing.T) {
	testFileInit(fileService.localFilePath + testfile)
	fileInfo, _ := os.Stat(fileService.localFilePath + testfile)
	assert.Equal(t, true, fileInfo.Size() > 0, "get file stat error")
	fileSize := uint64(fileInfo.Size())
	bufferSize := uint64(Buffer_Size)
	buf := make([]byte, bufferSize)
	for offset := uint64(0); offset < fileSize; offset += bufferSize {
		if fileSize-offset < bufferSize {
			bufferSize = fileSize - offset
		}

		done := fileService.saveToFileObject(testfile, offset, uint32(bufferSize), buf, fileSize)
		if offset+bufferSize < fileSize && !done {
			assert.Equal(t, false, done, fmt.Sprint("wrong saving, offset=", offset,
				", buffer size =", bufferSize, ", fileSize =", fileSize))
		} else {
			assert.Equal(t, true, done, fmt.Sprint("wrong complete, offset=", offset,
				", buffer size =", bufferSize, ", fileSize =", fileSize))
			break
		}
		bufferSize = Buffer_Size
	}
}

func TestSaveToFileSystem(t *testing.T) {
	randfilename := testfile + fmt.Sprint(rand.Intn(100))
	filepath := fileService.localFilePath + randfilename
	testFileInit(filepath)
	fileInfo, _ := os.Stat(filepath)
	assert.Equal(t, true, fileInfo.Size() > 0, "get file stat error")
	fileSize := uint64(fileInfo.Size())
	bufferSize := uint64(Buffer_Size)
	buf := make([]byte, bufferSize)
	for offset := uint64(0); offset < fileSize; offset += bufferSize {
		if fileSize-offset < bufferSize {
			bufferSize = fileSize - offset
		}

		done := fileService.saveToFileObject(randfilename, offset, uint32(bufferSize), buf, fileSize)
		if offset+bufferSize >= fileSize {
			assert.Equal(t, true, done, fmt.Sprint("wrong complete, offset=", offset,
				", buffer size =", bufferSize, ", fileSize =", fileSize))
		} else {
			assert.Equal(t, false, done, fmt.Sprint("wrong saving, offset=", offset,
				", buffer size =", bufferSize, ", fileSize =", fileSize))
		}
		bufferSize = Buffer_Size
	}
	os.Remove(filepath)

}
