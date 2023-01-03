package fileserver

import (
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
)

type FileServer struct {
	udp           *p2p.UdpServer
	owner         *gcrypto.GhostAddress
	localAddr     *ptypes.GhostIp
	localFilePath string

	packetSqHandler map[packets.PacketSecondType]func(*packets.Header, *net.UDPAddr) []p2p.ResponsePacketInfo
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
	fileServer.packetSqHandler = make(map[packets.PacketSecondType]func(*packets.Header, *net.UDPAddr) []p2p.ResponsePacketInfo)
	fileServer.packetCqHandler = make(map[packets.PacketSecondType]func(*packets.Header, *net.UDPAddr))

	fileServer.packetSqHandler[packets.PacketSecondType_RequestFile] = fileServer.RequestFileSq
	fileServer.packetSqHandler[packets.PacketSecondType_ResponseFile] = fileServer.ResponseFileSq

	fileServer.packetCqHandler[packets.PacketSecondType_RequestFile] = fileServer.RequestFileCq
	fileServer.packetCqHandler[packets.PacketSecondType_ResponseFile] = fileServer.ResponseFileCq

	packetFactory.RegisterPacketHandler(packets.PacketType_FileTransfer, fileServer.packetSqHandler, fileServer.packetCqHandler)
}

func (fileServer *FileServer) CheckFileExist(filename string) {

}

func (fileServer *FileServer) CreateFile(filename string, data []byte) {

}

// LoadFileToObj -> load to memory
func (fileServer *FileServer) LoadFileToObj(filename string, data []byte) {
}

// SendGetFileInfo -> RequestFilePacketSq
func (fileServer *FileServer) SendGetFileInfo(filename string, data []byte) {
}

// SendFileInfo -> RequestFilePacketCq
func (fileServer *FileServer) SendFileInfo(filename string, data []byte) {
}

// SendFileData -> ResponseFileSq
func (fileServer *FileServer) SendFileData(filename string, data []byte) {
}
