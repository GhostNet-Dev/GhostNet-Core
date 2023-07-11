package blockmanager

import (
	"log"
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/fileservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"google.golang.org/protobuf/proto"
)

func (blockMgr *BlockManager) SendTransactionSq(header *packets.Header, from *net.UDPAddr) []p2p.ResponseHeaderInfo {
	sq := &packets.SendTransactionSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	filename := fileservice.ByteToFilename(sq.TxId)
	fileObj := blockMgr.cloud.ReadFromCloudSync(filename, from)
	txValidate := false

	if fileObj != nil {
		txValidate = blockMgr.DownloadTransaction(fileObj, nil)
		if !txValidate {
			blockMgr.fileService.DeleteFile(filename)
		}
	}

	//master.blockHandler.SendTransaction(sq.TxId)
	cq := packets.SendTransactionCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), sq.Master.GetRequestId(), 0, blockMgr.localIpAddr),
		Result: txValidate,
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendTransaction,
			PacketData: cqData,
			RequestId:  cq.Master.GetRequestId(),
			SqFlag:     false,
		},
	}
}
func (blockMgr *BlockManager) SendTransactionCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) SearchTransactionSq(header *packets.Header, from *net.UDPAddr) []p2p.ResponseHeaderInfo {
	sq := &packets.SearchTransactionSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	//master.blockHandler.SendTransaction(sq.TxId)
	cq := packets.SearchTransactionCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), sq.Master.GetRequestId(), 0, blockMgr.localIpAddr),
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SearchTransaction,
			PacketData: cqData,
			RequestId:  cq.Master.GetRequestId(),
			SqFlag:     false,
		},
	}
}

func (blockMgr *BlockManager) SearchTransactionCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) SendDataTransactionSq(header *packets.Header, from *net.UDPAddr) []p2p.ResponseHeaderInfo {
	sq := &packets.SendDataTransactionSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	txFilename := fileservice.ByteToFilename(sq.TxId)
	dataTxFilename := fileservice.ByteToFilename(sq.DataTxId)
	txFileObj := blockMgr.cloud.ReadFromCloudSync(txFilename, from)
	dataTxFileObj := blockMgr.cloud.ReadFromCloudSync(dataTxFilename, from)
	result := false

	if txFileObj != nil && dataTxFileObj != nil {
		if !blockMgr.DownloadDataTransaction(txFileObj.Buffer, dataTxFileObj.Buffer) {
			blockMgr.fileService.DeleteFile(txFilename)
		} else {
			result = true
		}
	}
	cq := packets.SendDataTransactionCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), sq.Master.GetRequestId(), 0, blockMgr.localIpAddr),
		Result: result,
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendDataTransaction,
			PacketData: cqData,
			RequestId:  cq.Master.GetRequestId(),
			SqFlag:     false,
		},
	}
}
func (blockMgr *BlockManager) SendDataTransactionCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) SearchDataTransactionSq(header *packets.Header, from *net.UDPAddr) []p2p.ResponseHeaderInfo {
	sq := &packets.SearchTransactionSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	cq := packets.SearchTransactionCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), sq.Master.GetRequestId(), 0, blockMgr.localIpAddr),
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SearchDataTransaction,
			PacketData: cqData,
			RequestId:  cq.Master.GetRequestId(),
			SqFlag:     false,
		},
	}
}

func (blockMgr *BlockManager) SearchDataTransactionCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) GetTxStatusSq(header *packets.Header, from *net.UDPAddr) []p2p.ResponseHeaderInfo {
	sq := &packets.GetTxStatusSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	cq := packets.GetTxStatusCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), sq.Master.GetRequestId(), 0, blockMgr.localIpAddr),
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	newSq := packets.SendTxStatusSq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), nil, 0, blockMgr.localIpAddr),
	}

	sendData, err := proto.Marshal(&newSq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_GetTxStatus,
			PacketData: cqData,
			RequestId:  cq.Master.GetRequestId(),
			SqFlag:     false,
		},
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendTxStatus,
			RequestId:  newSq.Master.GetRequestId(),
			PacketData: sendData,
			SqFlag:     true,
		},
	}
}

func (blockMgr *BlockManager) GetTxStatusCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) SendTxStatusSq(header *packets.Header, from *net.UDPAddr) []p2p.ResponseHeaderInfo {
	sq := &packets.SendTxStatusSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	cq := packets.SendTxStatusCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), sq.Master.GetRequestId(), 0, blockMgr.localIpAddr),
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendTxStatus,
			PacketData: cqData,
			RequestId:  cq.Master.GetRequestId(),
			SqFlag:     false,
		},
	}
}

func (blockMgr *BlockManager) SendTxStatusCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) CheckRootFsSq(header *packets.Header, from *net.UDPAddr) []p2p.ResponseHeaderInfo {
	sq := &packets.CheckRootFsSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	exist := blockMgr.blockContainer.TxContainer.CheckExistFsRoot(sq.Nickname)

	cq := packets.CheckRootFsCq{
		Master:   p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), sq.Master.GetRequestId(), 0, blockMgr.localIpAddr),
		Nickname: sq.Nickname,
		Exist:    exist,
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_CheckRootFs,
			PacketData: cqData,
			RequestId:  cq.Master.GetRequestId(),
			SqFlag:     false,
		},
	}
}

func (blockMgr *BlockManager) CheckRootFsCq(header *packets.Header, from *net.UDPAddr) {
	cq := &packets.CheckRootFsCq{}
	if err := proto.Unmarshal(header.PacketData, cq); err != nil {
		log.Fatal(err)
	}
	blockMgr.callback(cq.Exist)
}
