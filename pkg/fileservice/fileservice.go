package fileservice

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

type FileService struct {
	udp            *p2p.UdpServer
	owner          *gcrypto.GhostAddress
	localAddr      *ptypes.GhostIp
	localFilePath  string
	fileObjManager *FileObjManager

	packetSqHandler map[packets.PacketSecondType]p2p.FuncPacketHandler
	packetCqHandler map[packets.PacketSecondType]p2p.FuncPacketHandler
}

func NewFileServer(udp *p2p.UdpServer, packetFactory *p2p.PacketFactory,
	owner *gcrypto.GhostAddress, ipAddr *ptypes.GhostIp, filePath string) *FileService {
	fileService := &FileService{
		udp:            udp,
		owner:          owner,
		localAddr:      ipAddr,
		localFilePath:  filePath,
		fileObjManager: NewFileObjManager(),
	}
	fileService.RegisterHandler(packetFactory)

	return fileService
}

func (fileService *FileService) RegisterHandler(packetFactory *p2p.PacketFactory) {
	fileService.packetSqHandler = make(map[packets.PacketSecondType]p2p.FuncPacketHandler)
	fileService.packetCqHandler = make(map[packets.PacketSecondType]p2p.FuncPacketHandler)

	fileService.packetSqHandler[packets.PacketSecondType_RequestFile] = fileService.RequestFileSq
	fileService.packetSqHandler[packets.PacketSecondType_ResponseFile] = fileService.ResponseFileSq

	fileService.packetCqHandler[packets.PacketSecondType_RequestFile] = fileService.RequestFileCq
	fileService.packetCqHandler[packets.PacketSecondType_ResponseFile] = fileService.ResponseFileCq

	packetFactory.RegisterPacketHandler(packets.PacketType_FileTransfer, fileService.packetSqHandler, fileService.packetCqHandler)
}

func (fileService *FileService) CheckFileExist(filename string) bool {
	if _, err := os.Stat(fileService.localFilePath + filename); os.IsNotExist(err) {
		return false
	}
	return true
}

func (fileService *FileService) CreateFile(filename string, data []byte,
	callback DoneHandler, context interface{}) *FileObject {
	fileFullPath := fileService.localFilePath + filename
	ioutil.WriteFile(fileFullPath, data, os.FileMode(644))
	return fileService.fileObjManager.CreateFileObj(filename, data, uint64(len(data)), callback, context)
}

func (fileService *FileService) commitFile(fileObj *FileObject) {
	fileFullPath := fileService.localFilePath + fileObj.Filename
	ioutil.WriteFile(fileFullPath, fileObj.Buffer, os.FileMode(644))
}

func (fileService *FileService) SendGetFile(filename string, ipAddr *net.UDPAddr, callback DoneHandler, context interface{}) {
	if fileService.CheckFileExist(filename) == false {
		fileService.sendGetFileInfo(filename, ipAddr, callback, context)
	} else {
		fileObj, exist := fileService.fileObjManager.GetFileObject(filename)
		if exist == false {
			fileObj = fileService.LoadFileToMemory(filename)
		}
		if callback != nil {
			callback(fileObj, context)
		}
	}
}

func (fileService *FileService) DeleteFile(filename string) {
	if fileService.CheckFileExist(filename) == true {
		os.Remove(filename)
		fileService.fileObjManager.DeleteObject(filename)
	}
}

// LoadFileToObj -> load to memory
func (fileService *FileService) LoadFileToMemory(filename string) *FileObject {
	fileFullPath := fileService.localFilePath + filename
	fileObj, ok := fileService.fileObjManager.GetFileObject(filename)
	if ok == true {
		return fileObj
	}

	buf, err := ioutil.ReadFile(fileFullPath)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return fileService.fileObjManager.CreateFileObj(filename, buf, uint64(len(buf)), nil, nil)
}

// sendGetFileInfo -> RequestFilePacketSq
func (fileService *FileService) sendGetFileInfo(filename string, ipAddr *net.UDPAddr,
	callback DoneHandler, context interface{}) {
	sq := &packets.RequestFilePacketSq{
		Master:      p2p.MakeMasterPacket(fileService.owner.GetPubAddress(), 0, 0, fileService.localAddr),
		RequestType: packets.FileRequestType_GetFileInfo,
		Filename:    filename,
	}

	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}

	fileService.fileObjManager.CreateFileObj(filename, nil, 0, callback, context)

	fileService.udp.SendUdpPacket(&p2p.PacketHeaderInfo{
		PacketType: packets.PacketType_FileTransfer,
		SecondType: packets.PacketSecondType_RequestFile,
		ThirdType:  packets.PacketThirdType_Reserved1,
		SqFlag:     true,
		PacketData: sendData,
	}, ipAddr)
}

// makeFileInfo -> RequestFilePacketCq
func (fileService *FileService) makeFileInfo(filename string) *p2p.PacketHeaderInfo {
	fileObj, exist := fileService.fileObjManager.GetFileObject(filename)

	if exist == false {
		if fileObj = fileService.LoadFileToMemory(filename); fileObj != nil {
			exist = true
		}
	}

	cq := &packets.RequestFilePacketCq{
		Master:   p2p.MakeMasterPacket(fileService.owner.GetPubAddress(), 0, 0, fileService.localAddr),
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

// sendFileData -> ResponseFileSq
func (fileService *FileService) sendFileData(filename string, startPos uint64, sequenceNum uint32, timeId uint64) *p2p.PacketHeaderInfo {
	var fileSize uint64 = 0
	var buf []byte

	if fileObj, exist := fileService.fileObjManager.GetFileObject(filename); exist == false {
		// TODO: CreateFileObj

		fp, err := os.Open(fileService.localFilePath + filename)
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
		Master:      p2p.MakeMasterPacket(fileService.owner.GetPubAddress(), 0, 0, fileService.localAddr),
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

func (fileService *FileService) saveToFileObject(filename string, startPos uint64, bufSize uint32, buffer []byte, fileLength uint64) bool {
	fileObj, exist := fileService.fileObjManager.GetFileObject(filename)
	if exist == false {
		fileObj = fileService.fileObjManager.CreateFileObj(filename, nil, fileLength, nil, nil)
	}
	copy(fileObj.Buffer[startPos:startPos+uint64(bufSize)], buffer[:])
	fileObj.UpdateFileImage(startPos)
	if fileObj.CheckComplete() == true {
		fileService.commitFile(fileObj)
	}
	return false
}
