package blockmanager

import (
	"log"
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gnetwork"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"google.golang.org/protobuf/proto"
)

func (blockMgr *BlockManager) InitHandler(master *gnetwork.MasterNetwork) {
	blockMgr.packetSqHandler[packets.PacketThirdType_GetHeightestBlock] = blockMgr.GetHeightestBlockSq
	blockMgr.packetSqHandler[packets.PacketThirdType_NewBlock] = blockMgr.NewBlockSq
	blockMgr.packetSqHandler[packets.PacketThirdType_GetBlock] = blockMgr.GetBlockSq
	blockMgr.packetSqHandler[packets.PacketThirdType_SendBlock] = blockMgr.SendBlockSq
	blockMgr.packetSqHandler[packets.PacketThirdType_SendTransaction] = blockMgr.SendTransactionSq
	blockMgr.packetSqHandler[packets.PacketThirdType_SearchTransaction] = blockMgr.SearchTransactionSq
	blockMgr.packetSqHandler[packets.PacketThirdType_SendDataTransaction] = blockMgr.SendDataTransactionSq
	blockMgr.packetSqHandler[packets.PacketThirdType_SearchDataTransaction] = blockMgr.SearchDataTransactionSq
	blockMgr.packetSqHandler[packets.PacketThirdType_GetBlockHash] = blockMgr.GetBlockHashSq
	blockMgr.packetSqHandler[packets.PacketThirdType_SendBlockHash] = blockMgr.SendBlockHashSq
	blockMgr.packetSqHandler[packets.PacketThirdType_GetTxStatus] = blockMgr.GetTxStatusSq
	blockMgr.packetSqHandler[packets.PacketThirdType_SendTxStatus] = blockMgr.SendTxStatusSq

	blockMgr.packetCqHandler[packets.PacketThirdType_GetHeightestBlock] = blockMgr.GetHeightestBlockCq
	blockMgr.packetCqHandler[packets.PacketThirdType_NewBlock] = blockMgr.NewBlockCq
	blockMgr.packetCqHandler[packets.PacketThirdType_GetBlock] = blockMgr.GetBlockCq
	blockMgr.packetCqHandler[packets.PacketThirdType_SendBlock] = blockMgr.SendBlockCq
	blockMgr.packetCqHandler[packets.PacketThirdType_SendTransaction] = blockMgr.SendTransactionCq
	blockMgr.packetCqHandler[packets.PacketThirdType_SearchTransaction] = blockMgr.SearchTransactionCq
	blockMgr.packetCqHandler[packets.PacketThirdType_SendDataTransaction] = blockMgr.SendDataTransactionCq
	blockMgr.packetCqHandler[packets.PacketThirdType_SearchDataTransaction] = blockMgr.SearchDataTransactionCq
	blockMgr.packetCqHandler[packets.PacketThirdType_GetBlockHash] = blockMgr.GetBlockHashCq
	blockMgr.packetCqHandler[packets.PacketThirdType_SendBlockHash] = blockMgr.SendBlockHashCq
	blockMgr.packetCqHandler[packets.PacketThirdType_GetTxStatus] = blockMgr.GetTxStatusCq
	blockMgr.packetCqHandler[packets.PacketThirdType_SendTxStatus] = blockMgr.SendTxStatusCq

	master.RegisterBlockHandler(blockMgr.BlockHandlerSq, blockMgr.BlockHandlerCq)
}

func (blockMgr *BlockManager) BlockHandlerSq(header *packets.Header, from *net.UDPAddr) []p2p.ResponseHeaderInfo {
	return blockMgr.packetSqHandler[header.ThirdType](header, from)
}

func (blockMgr *BlockManager) BlockHandlerCq(header *packets.Header, from *net.UDPAddr) {
	blockMgr.packetCqHandler[header.ThirdType](header, from)
}

func (blockMgr *BlockManager) GetHeightestBlockSq(header *packets.Header, from *net.UDPAddr) []p2p.ResponseHeaderInfo {
	sq := &packets.GetHeightestBlockSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	cq := packets.GetHeightestBlockCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
		Height: blockMgr.blockContainer.BlockHeight(),
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_GetHeightestBlock,
			PacketData: sendData,
			SqFlag:     false,
		},
	}
}

func (blockMgr *BlockManager) GetHeightestBlockCq(header *packets.Header, from *net.UDPAddr) {
	cq := &packets.GetHeightestBlockCq{}
	if err := proto.Unmarshal(header.PacketData, cq); err != nil {
		log.Fatal(err)
	}
	blockMgr.fsm.State().RecvBlockHeight(cq.Height, cq.Master.Common.FromPubKeyAddress)
}

func (blockMgr *BlockManager) NewBlockSq(header *packets.Header, from *net.UDPAddr) []p2p.ResponseHeaderInfo {
	sq := &packets.NewBlockSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	if blockMgr.fsm.CheckAcceptNewBlock() {
		fileObj := blockMgr.cloud.ReadFromCloudSync(sq.BlockFilename, from)
		//blockMgr.RequestBlockChainFile(sq.BlockFilename, from, blockMgr.DownloadNewBlock, nil)
		blockMgr.DownloadNewBlock(fileObj, nil)
	}

	cq := packets.NewBlockCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_NewBlock,
			PacketData: sendData,
			SqFlag:     false,
		},
	}
}

func (blockMgr *BlockManager) NewBlockCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) GetBlockSq(header *packets.Header, from *net.UDPAddr) []p2p.ResponseHeaderInfo {
	sq := &packets.GetBlockSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	blockFilename, exist := blockMgr.PrepareSendBlock(sq.BlockId)

	cq := packets.GetBlockCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
		Result: exist,
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	responseInfo := []p2p.ResponseHeaderInfo{{
		ToAddr:     from,
		ThirdType:  packets.PacketThirdType_GetBlock,
		PacketData: cqData,
		SqFlag:     false,
	}}

	if exist {
		newSq := &packets.SendBlockSq{
			Master:        p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
			BlockFilename: blockFilename,
		}

		sendData, err := proto.Marshal(newSq)
		if err != nil {
			log.Fatal(err)
		}
		responseInfo = append(responseInfo, p2p.ResponseHeaderInfo{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendBlock,
			PacketData: sendData,
			SqFlag:     true,
		})
	}

	return responseInfo
}

func (blockMgr *BlockManager) GetBlockCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) SendBlockSq(header *packets.Header, from *net.UDPAddr) []p2p.ResponseHeaderInfo {
	sq := &packets.SendBlockSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	if blockMgr.fsm.CheckAcceptDownloadBlock() {
		fileObj := blockMgr.cloud.ReadFromCloudSync(sq.BlockFilename, from)
		blockMgr.DownloadBlock(fileObj, sq.GetMaster().Common.FromPubKeyAddress)
	}
	//blockMgr.RequestBlockChainFile(sq.BlockFilename, from, blockMgr.DownloadBlock, nil)
	//master.blockHandler.SendBlock(sq.BlockFilename)
	cq := packets.SendBlockCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendBlock,
			PacketData: cqData,
			SqFlag:     false,
		},
	}
}

func (blockMgr *BlockManager) SendBlockCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) GetBlockHashSq(header *packets.Header, from *net.UDPAddr) []p2p.ResponseHeaderInfo {
	sq := &packets.GetBlockHashSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	cq := packets.GetBlockHashCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	responseInfos := []p2p.ResponseHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_GetBlockHash,
			PacketData: cqData,
			SqFlag:     false,
		},
	}

	pairedBlock := blockMgr.blockContainer.GetBlock(sq.BlockId)
	if pairedBlock != nil {
		hash := pairedBlock.Block.GetHashKey()

		newSq := packets.SendBlockHashSq{
			Master:  p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
			Hash:    hash,
			BlockId: sq.BlockId,
		}

		sendData, err := proto.Marshal(&newSq)
		if err != nil {
			log.Fatal(err)
		}
		responseInfos = append(responseInfos, p2p.ResponseHeaderInfo{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendBlockHash,
			PacketData: sendData,
			SqFlag:     true,
		})
	}

	return responseInfos
}

func (blockMgr *BlockManager) GetBlockHashCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) SendBlockHashSq(header *packets.Header, from *net.UDPAddr) []p2p.ResponseHeaderInfo {
	sq := &packets.SendBlockHashSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	cq := packets.SendBlockHashCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	blockMgr.fsm.State().RecvBlockHash(sq.Master.Common.FromPubKeyAddress, sq.Hash, sq.BlockId)

	return []p2p.ResponseHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendBlockHash,
			PacketData: cqData,
			SqFlag:     false,
		},
	}
}

func (blockMgr *BlockManager) SendBlockHashCq(header *packets.Header, from *net.UDPAddr) {}
