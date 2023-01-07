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
	blockMgr.packetSqHandler[packets.PacketThirdType_GetBlockPrevHash] = blockMgr.GetBlockPrevHashSq
	blockMgr.packetSqHandler[packets.PacketThirdType_SendBlockPrevHash] = blockMgr.SendBlockPrevHashSq
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
	blockMgr.packetCqHandler[packets.PacketThirdType_GetBlockPrevHash] = blockMgr.GetBlockPrevHashCq
	blockMgr.packetCqHandler[packets.PacketThirdType_SendBlockPrevHash] = blockMgr.SendBlockPrevHashCq
	blockMgr.packetCqHandler[packets.PacketThirdType_GetTxStatus] = blockMgr.GetTxStatusCq
	blockMgr.packetCqHandler[packets.PacketThirdType_SendTxStatus] = blockMgr.SendTxStatusCq

	master.RegisterBlockHandler(blockMgr.BlockHandlerSq, blockMgr.BlockHandlerCq)
}

func (blockMgr *BlockManager) BlockHandlerSq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
	return blockMgr.packetSqHandler[header.ThirdType](header, from)
}

func (blockMgr *BlockManager) BlockHandlerCq(header *packets.Header, from *net.UDPAddr) {
	blockMgr.packetCqHandler[header.ThirdType](header, from)
}

func (blockMgr *BlockManager) GetHeightestBlockSq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
	sq := &packets.GetHeightestBlockSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	cq := packets.GetHeightestBlockCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.PacketHeaderInfo{
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
}

func (blockMgr *BlockManager) NewBlockSq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
	sq := &packets.NewBlockSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	cq := packets.NewBlockCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
	}

	sendData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}
	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_NewBlock,
			PacketData: sendData,
			SqFlag:     false,
		},
	}
}

func (blockMgr *BlockManager) NewBlockCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) GetBlockSq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
	sq := &packets.GetBlockSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	//blockname := master.blockHandler.GetBlock(sq.BlockId)

	cq := packets.GetBlockCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	newSq := &packets.SendBlockSq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
		//BlockFilename: blockname,
	}

	sendData, err := proto.Marshal(newSq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_GetBlock,
			PacketData: cqData,
			SqFlag:     false,
		},
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendBlock,
			PacketData: sendData,
			SqFlag:     true,
		},
	}
}

func (blockMgr *BlockManager) GetBlockCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) SendBlockSq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
	sq := &packets.SendBlockSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	//master.blockHandler.SendBlock(sq.BlockFilename)
	cq := packets.SendBlockCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendBlock,
			PacketData: cqData,
			SqFlag:     false,
		},
	}
}

func (blockMgr *BlockManager) SendBlockCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) SendTransactionSq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
	sq := &packets.SendTransactionSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	//master.blockHandler.SendTransaction(sq.TxId)
	cq := packets.SendTransactionCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendTransaction,
			PacketData: cqData,
			SqFlag:     false,
		},
	}
}
func (blockMgr *BlockManager) SendTransactionCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) SearchTransactionSq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
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

	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SearchTransaction,
			PacketData: cqData,
			SqFlag:     false,
		},
	}
}

func (blockMgr *BlockManager) SearchTransactionCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) SendDataTransactionSq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
	sq := &packets.SendDataTransactionSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	cq := packets.SendDataTransactionCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendDataTransaction,
			PacketData: cqData,
			SqFlag:     false,
		},
	}
}
func (blockMgr *BlockManager) SendDataTransactionCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) SearchDataTransactionSq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
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

	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SearchDataTransaction,
			PacketData: cqData,
			SqFlag:     false,
		},
	}
}

func (blockMgr *BlockManager) SearchDataTransactionCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) GetBlockHashSq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
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

	newSq := packets.SendBlockHashSq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
	}

	sendData, err := proto.Marshal(&newSq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_GetBlockHash,
			PacketData: cqData,
			SqFlag:     false,
		},
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendBlockHash,
			PacketData: sendData,
			SqFlag:     true,
		},
	}
}

func (blockMgr *BlockManager) GetBlockHashCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) SendBlockHashSq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
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

	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendBlockHash,
			PacketData: cqData,
			SqFlag:     false,
		},
	}
}

func (blockMgr *BlockManager) SendBlockHashCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) GetBlockPrevHashSq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
	sq := &packets.GetBlockPrevHashSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	cq := packets.GetBlockPrevHashCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	newSq := packets.SendBlockPrevHashSq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
	}

	sendData, err := proto.Marshal(&newSq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_GetBlockPrevHash,
			PacketData: cqData,
			SqFlag:     false,
		},
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendBlockPrevHash,
			PacketData: sendData,
			SqFlag:     true,
		},
	}
}

func (blockMgr *BlockManager) GetBlockPrevHashCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) SendBlockPrevHashSq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
	sq := &packets.SendBlockPrevHashSq{}
	if err := proto.Unmarshal(header.PacketData, sq); err != nil {
		log.Fatal(err)
	}

	cq := packets.SendBlockPrevHashCq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), 0, 0, blockMgr.localIpAddr),
	}

	cqData, err := proto.Marshal(&cq)
	if err != nil {
		log.Fatal(err)
	}

	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendBlockPrevHash,
			PacketData: cqData,
			SqFlag:     false,
		},
	}
}

func (blockMgr *BlockManager) SendBlockPrevHashCq(header *packets.Header, from *net.UDPAddr) {}

func (blockMgr *BlockManager) GetTxStatusSq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
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

	return []p2p.PacketHeaderInfo{
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

func (blockMgr *BlockManager) SendTxStatusSq(header *packets.Header, from *net.UDPAddr) []p2p.PacketHeaderInfo {
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

	return []p2p.PacketHeaderInfo{
		{
			ToAddr:     from,
			ThirdType:  packets.PacketThirdType_SendTxStatus,
			PacketData: cqData,
			SqFlag:     false,
		},
	}
}

func (blockMgr *BlockManager) SendTxStatusCq(header *packets.Header, from *net.UDPAddr) {}
