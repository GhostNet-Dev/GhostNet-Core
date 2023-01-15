package fileserver

import (
	"io/ioutil"
	"log"
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

	packetSqHandler map[packets.PacketSecondType]p2p.FuncPacketHandler
	packetCqHandler map[packets.PacketSecondType]p2p.FuncPacketHandler
}

func NewFileServer(udp *p2p.UdpServer, packetFactory *p2p.PacketFactory,
	owner *gcrypto.GhostAddress, ipAddr *ptypes.GhostIp, filePath string) *FileServer {
	fileServer := &FileServer{
		udp:            udp,
		owner:          owner,
		localAddr:      ipAddr,
		localFilePath:  filePath,
		fileObjManager: NewFileObjManager(),
	}
	fileServer.RegisterHandler(packetFactory)

	return fileServer
}

func (fileServer *FileServer) RegisterHandler(packetFactory *p2p.PacketFactory) {
	fileServer.packetSqHandler = make(map[packets.PacketSecondType]p2p.FuncPacketHandler)
	fileServer.packetCqHandler = make(map[packets.PacketSecondType]p2p.FuncPacketHandler)

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
	return fileServer.fileObjManager.CreateFileObj(filename, data, uint64(len(data)), callback, context)
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
	return fileServer.fileObjManager.CreateFileObj(filename, buf, uint64(len(buf)), nil, nil)
}

// SendGetFileInfo -> RequestFilePacketSq
func (fileServer *FileServer) SendGetFileInfo(ipAddr *ptypes.GhostIp, filename string,
	callback DoneHandler, context interface{}) {
	sq := &packets.RequestFilePacketSq{
		Master:      p2p.MakeMasterPacket(fileServer.owner.GetPubAddress(), 0, 0, fileServer.localAddr),
		RequestType: packets.FileRequestType_GetFileInfo,
		Filename:    filename,
	}

	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}

	fileServer.fileObjManager.CreateFileObj(filename, nil, 0, callback, context)

	fileServer.udp.SendPacket(&p2p.PacketHeaderInfo{
		PacketType: packets.PacketType_FileTransfer,
		SecondType: packets.PacketSecondType_RequestFile,
		ThirdType:  packets.PacketThirdType_Reserved1,
		SqFlag:     true,
		PacketData: sendData,
	}, ipAddr)
}

// MakeFileInfo -> RequestFilePacketCq
func (fileServer *FileServer) MakeFileInfo(filename string) *p2p.PacketHeaderInfo {
	fileObj, exist := fileServer.fileObjManager.GetFileObject(filename)

	if exist == false {
		if fileObj = fileServer.LoadFileToMemory(filename); fileObj != nil {
			exist = true
		}
	}

	cq := &packets.RequestFilePacketCq{
		Master:   p2p.MakeMasterPacket(fileServer.owner.GetPubAddress(), 0, 0, fileServer.localAddr),
		Filename: filename,
		Result:   exist,
	}

	if exist == true {
		cq.FileLength = fileObj.FileLength
	}

	sendData, err := proto.Marshal(cq)
	if err != nil {
		log.Fatal(err)
	}

	return &p2p.PacketHeaderInfo{
		PacketType: packets.PacketType_FileTransfer,
		SecondType: packets.PacketSecondType_RequestFile,
		ThirdType:  packets.PacketThirdType_Reserved1,
		SqFlag:     false,
		PacketData: sendData,
	}
}

func FileErrorCheck(err error) bool {
	if err != nil {
		log.Fatal(err)
		return false
	}
	return true
}

// SendFileData -> ResponseFileSq
func (fileServer *FileServer) SendFileData(filename string, startPos uint64, sequenceNum uint32, timeId uint64) *p2p.PacketHeaderInfo {
	var fileSize uint64 = 0
	var buf []byte

	if fileObj, exist := fileServer.fileObjManager.GetFileObject(filename); exist == false {
		// TODO: CreateFileObj

		fp, err := os.Open(fileServer.localFilePath + filename)
		if ret := FileErrorCheck(err); ret == false {
			return nil
		}
		defer fp.Close()

		fileInfo, err := fp.Stat()
		if ret := FileErrorCheck(err); ret == false {
			return nil
		}
		fileSize = uint64(fileInfo.Size())

		if startPos > uint64(fileInfo.Size()) {
			return nil
		}

		fp.Seek(int64(startPos), 0)

		readSize := uint64(fileInfo.Size()) - startPos
		if readSize > BufferSize {
			readSize = BufferSize
		}
		buf = make([]byte, readSize)
		fp.Read(buf)
	} else {
		fileSize = fileObj.FileLength
		readSize := fileSize - startPos
		if readSize > BufferSize {
			readSize = BufferSize
		}
		buf = make([]byte, readSize)
		copy(buf, fileObj.Buffer[startPos:startPos+readSize])
	}

	sq := &packets.ResponseFilePacketSq{
		Master:      p2p.MakeMasterPacket(fileServer.owner.GetPubAddress(), 0, 0, fileServer.localAddr),
		Filename:    filename,
		StartPos:    startPos,
		FileData:    buf,
		BufferSize:  uint32(len(buf)),
		FileLength:  fileSize,
		SequenceNum: sequenceNum,
	}

	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}

	return &p2p.PacketHeaderInfo{
		PacketType: packets.PacketType_FileTransfer,
		SecondType: packets.PacketSecondType_ResponseFile,
		ThirdType:  packets.PacketThirdType_Reserved1,
		SqFlag:     true,
		PacketData: sendData,
	}
}

func (fileServer *FileServer) SaveToFileObject(filename string, startPos uint64, bufSize uint32, buffer []byte, fileLength uint64) bool {
	fileObj, exist := fileServer.fileObjManager.GetFileObject(filename)
	if exist == false {
		fileObj = fileServer.fileObjManager.CreateFileObj(filename, nil, fileLength, nil, nil)
	}
	copy(fileObj.Buffer[startPos:startPos+uint64(bufSize)], buffer[:])
	fileObj.UpdateFileImage(startPos)
	return fileObj.CheckComplete()
}
