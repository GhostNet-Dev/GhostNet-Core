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
	fileObj := blockMgr.cloud.DownloadSync(filename, from)
	blockMgr.DownloadTransaction(fileObj, nil)

	//master.blockHandler.SendTransaction(sq.TxId)
	cq := packets.SendTransactionCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
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
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
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
	txFileObj := blockMgr.cloud.DownloadSync(txFilename, from)
	dataTxFileObj := blockMgr.cloud.DownloadSync(dataTxFilename, from)
	blockMgr.DownloadDataTransaction(txFileObj.Buffer, dataTxFileObj.Buffer)

	cq := packets.SendDataTransactionCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
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
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
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
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	newSq := packets.SendTxStatusSq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
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
			SqFlag:     false,
		},
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendTxStatus,
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
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
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
			SqFlag:     false,
		},
	}
}

func (blockMgr *BlockManager) SendTxStatusCq(header *packets.Header, from *net.UDPAddr) {}
