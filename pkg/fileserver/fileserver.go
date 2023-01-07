package fileserver

import (
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"google.golang.org/protobuf/proto"
)

const BufferSize = 1024

type FileServer struct {
	udp            *p2p.UdpServer
	owner          *gcrypto.GhostAddress
	localAddr      *ptypes.GhostIp
	localFilePath  string
	fileObjManager *FileObjManager

	packetSqHandler map[packets.PacketSecondType]func(*packets.Header, *net.UDPAddr) []p2p.PacketHeaderInfo
	packetCqHandler map[packets.PacketSecondType]func(*packets.Header, *net.UDPAddr)
}

func NewFileServer(udp *p2p.UdpServer, packetFactory *p2p.PacketFactory,
	owner *gcrypto.GhostAddress, ipAddr *ptypes.GhostIp, filePath string) *FileServer {
	fileServer := &FileServer{
		udp:           udp,
		owner:         owner,
		localAddr:     ipAddr,
		localFilePath: filePath,
	}
	fileServer.RegisterHandler(packetFactory)

	return fileServer
}

func (fileServer *FileServer) RegisterHandler(packetFactory *p2p.PacketFactory) {
	fileServer.packetSqHandler = make(map[packets.PacketSecondType]func(*packets.Header, *net.UDPAddr) []p2p.PacketHeaderInfo)
	fileServer.packetCqHandler = make(map[packets.PacketSecondType]func(*packets.Header, *net.UDPAddr))

	fileServer.packetSqHandler[packets.PacketSecondType_RequestFile] = fileServer.RequestFileSq
	fileServer.packetSqHandler[packets.PacketSecondType_ResponseFile] = fileServer.ResponseFileSq

	fileServer.packetCqHandler[packets.PacketSecondType_RequestFile] = fileServer.RequestFileCq
	fileServer.packetCqHandler[packets.PacketSecondType_ResponseFile] = fileServer.ResponseFileCq

	packetFactory.RegisterPacketHandler(packets.PacketType_FileTransfer, fileServer.packetSqHandler, fileServer.packetCqHandler)
}

func (fileServer *FileServer) CheckFileExist(filename string) bool {
	if _, err := os.Stat(fileServer.localFilePath + filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func (fileServer *FileServer) CreateFile(filename string, data []byte,
	callback DoneHandler, context interface{}) *FileObject {
	fileFullPath := fileServer.localFilePath + filename
	ioutil.WriteFile(fileFullPath, data, os.FileMode(644))
	return fileServer.fileObjManager.CreateFileObj(filename, data, callback, context)
}

// LoadFileToObj -> load to memory
func (fileServer *FileServer) LoadFileToMemory(filename string) *FileObject {
	fileFullPath := fileServer.localFilePath + filename
	fileObj, ok := fileServer.fileObjManager.GetFileObject(filename)
	if ok == true {
		return fileObj
	}

	buf, err := ioutil.ReadFile(fileFullPath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return fileServer.fileObjManager.CreateFileObj(filename, buf, nil, nil)
}

// SendGetFileInfo -> RequestFilePacketSq
func (fileServer *FileServer) SendGetFileInfo(ipAddr *ptypes.GhostIp, filename string, offset uint64) {
	sq := &packets.RequestFilePacketSq{
		Master:      p2p.MakeMasterPacket(fileServer.owner.GetPubAddress(), 0, 0, fileServer.localAddr),
		Filename:    filename,
		StartOffset: offset,
	}

	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}

	fileServer.udp.SendPacket(&p2p.PacketHeaderInfo{
		PacketType: packets.PacketType_FileTransfer,
		SecondType: packets.PacketSecondType_RequestFile,
		ThirdType:  packets.PacketThirdType_Reserved1,
		SqFlag:     true,
		PacketData: sendData,
	}, ipAddr)
}

// SendFileInfo -> RequestFilePacketCq
func (fileServer *FileServer) SendFileInfo(ipAddr *ptypes.GhostIp, filename string, fileLength uint64, exist bool) {
	cq := &packets.RequestFilePacketCq{
		Master:     p2p.MakeMasterPacket(fileServer.owner.GetPubAddress(), 0, 0, fileServer.localAddr),
		Filename:   filename,
		FileLength: fileLength,
		Result:     exist,
	}

	sendData, err := proto.Marshal(cq)
	if err != nil {
		log.Fatal(err)
	}

	fileServer.udp.SendPacket(&p2p.PacketHeaderInfo{
		PacketType: packets.PacketType_FileTransfer,
		SecondType: packets.PacketSecondType_RequestFile,
		ThirdType:  packets.PacketThirdType_Reserved1,
		SqFlag:     false,
		PacketData: sendData,
	}, ipAddr)
}

// SendFileData -> ResponseFileSq
func (fileServer *FileServer) SendFileData(ipAddr *ptypes.GhostIp, filename string, startPos uint64, sequenceNum uint32, timeId uint64) {
	// TODO: Read file
	sq := &packets.ResponseFilePacketSq{
		Master:      p2p.MakeMasterPacket(fileServer.owner.GetPubAddress(), 0, 0, fileServer.localAddr),
		Filename:    filename,
		StartPos:    startPos,
		FileData:    nil,
		BufferSize:  BufferSize,
		FileLength:  0,
		SequenceNum: sequenceNum,
	}

	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}

	fileServer.udp.SendPacket(&p2p.PacketHeaderInfo{
		PacketType: packets.PacketType_FileTransfer,
		SecondType: packets.PacketSecondType_RequestFile,
		ThirdType:  packets.PacketThirdType_Reserved1,
		SqFlag:     true,
		PacketData: sendData,
	}, ipAddr)
}
