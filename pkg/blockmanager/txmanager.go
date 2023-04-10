package blockmanager

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/fileservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gnetwork"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"google.golang.org/protobuf/proto"
)

func (blockMgr *BlockManager) SendTx(tx *types.GhostTransaction) {
	filename := fileservice.ByteToFilename(tx.TxId)
	if exist := blockMgr.fileService.CheckFileExist(filename); !exist {
		blockMgr.fileService.CreateFile(filename, tx.SerializeToByte(), nil, nil)
	}
	sq := &packets.SendTransactionSq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
		TxId:   tx.TxId,
	}
	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}
	headerInfo := &p2p.ResponseHeaderInfo{
		PacketType: packets.PacketType_MasterNetwork,
		SecondType: packets.PacketSecondType_BlockChain,
		ThirdType:  packets.PacketThirdType_SendTransaction,
		PacketData: sendData,
		SqFlag:     true,
	}
	blockMgr.master.SendToMasterNodeGrpSq(packets.RoutingType_BroadCastingLevelZero, gnetwork.DefaultTreeLevel, headerInfo)
}

func (blockMgr *BlockManager) SendDataTx(tx *types.GhostTransaction, dataTx *types.GhostDataTransaction) {
	filename := fileservice.ByteToFilename(tx.TxId)
	dataFilename := fileservice.ByteToFilename(dataTx.TxId)
	if exist := blockMgr.fileService.CheckFileExist(filename); !exist {
		blockMgr.fileService.CreateFile(filename, tx.SerializeToByte(), nil, nil)
	}
	if exist := blockMgr.fileService.CheckFileExist(dataFilename); !exist {
		blockMgr.fileService.CreateFile(dataFilename, dataTx.SerializeToByte(), nil, nil)
	}
	sq := &packets.SendDataTransactionSq{
		Master:   p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
		TxId:     tx.TxId,
		DataTxId: dataTx.TxId,
	}
	sendData, err := proto.Marshal(sq)
	if err != nil {
		log.Fatal(err)
	}
	headerInfo := &p2p.ResponseHeaderInfo{
		PacketType: packets.PacketType_MasterNetwork,
		SecondType: packets.PacketSecondType_BlockChain,
		ThirdType:  packets.PacketThirdType_SendDataTransaction,
		PacketData: sendData,
		SqFlag:     true,
	}
	blockMgr.master.SendToMasterNodeGrpSq(packets.RoutingType_BroadCastingLevelZero, gnetwork.DefaultTreeLevel, headerInfo)
}
