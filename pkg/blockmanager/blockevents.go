package blockmanager

import (
	"fmt"
	"log"
	"sync"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gnetwork"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"google.golang.org/protobuf/proto"
)

func (blockMgr *BlockManager) BroadcastBlockChainNotification() {
	sq := packets.GetHeightestBlockSq{
		Master: p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), nil, 0, blockMgr.localIpAddr),
	}
	sendData, err := proto.Marshal(&sq)
	if err != nil {
		log.Fatal(err)
	}
	headerInfo := &p2p.ResponseHeaderInfo{
		PacketType: packets.PacketType_MasterNetwork,
		SecondType: packets.PacketSecondType_BlockChain,
		ThirdType:  packets.PacketThirdType_GetHeightestBlock,
		RequestId:  sq.Master.GetRequestId(),
		PacketData: sendData,
		SqFlag:     true,
	}
	blockMgr.master.SendToMasterNodeGrpSq(packets.RoutingType_BroadCastingLevelZero, gnetwork.DefaultTreeLevel, headerInfo)
}

func (blockMgr *BlockManager) MiningStart() {
	blockMgr.block.MinerStart()
}

func (blockMgr *BlockManager) MiningStop() {
	blockMgr.block.MinerStop()
}

func (blockMgr *BlockManager) RequestGetBlock(pubKey string, blockIdx uint32) {
	sq := packets.GetBlockSq{
		Master:  p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), nil, 0, blockMgr.localIpAddr),
		BlockId: blockIdx,
	}
	sendData, err := proto.Marshal(&sq)
	if err != nil {
		log.Fatal(err)
	}
	blockMgr.master.SendToMasterNodeSq(packets.PacketThirdType_GetBlock, pubKey, sendData, sq.Master.GetRequestId())
}

func (blockMgr *BlockManager) RequestGetBlockHash(pubKey string, blockIdx uint32) {
	sq := packets.GetBlockHashSq{
		Master:  p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), nil, 0, blockMgr.localIpAddr),
		BlockId: blockIdx,
	}
	sendData, err := proto.Marshal(&sq)
	if err != nil {
		log.Fatal(err)
	}
	blockMgr.master.SendToMasterNodeSq(packets.PacketThirdType_GetBlockHash, pubKey, sendData, sq.Master.GetRequestId())
}

func (blockMgr *BlockManager) MergeErrorNotification(pubKey string, result bool) {
	blockMgr.glog.DebugOutput(blockMgr, fmt.Sprint("Merge Error = ", pubKey), glogger.BlockConsensus)
}

func (blockMgr *BlockManager) BlockServerInitStart() {
	blockMgr.consensus.Clear()
	blockMgr.MiningStart()
	blockMgr.glog.DebugOutput(blockMgr, "Mining Start", glogger.BlockConsensus)
	//TODO it needs to more clear!

}

func (blockMgr *BlockManager) CheckHeightForRebuild(neighborHeight uint32) bool {
	currHeight := blockMgr.blockContainer.BlockHeight()

	return currHeight < neighborHeight
}

func (blockMgr *BlockManager) CheckValidNode(candidatePool *sync.Map, maxHeight uint32) (string, []string, uint32) {
	trycount := 0
	candidatePool.Range(func(k, v interface{}) bool {
		trycount++
		return true
	})
	i := 0
	for target := maxHeight; i < trycount; target-- {
		if candiList, exist := candidatePool.Load(target); exist {
			for _, pubKey := range candiList.([]string) {
				if blockMgr.master.CheckNodeInfo(pubKey) {
					return pubKey, candiList.([]string), target
				}
				i++
			}
		}
	}
	return "", nil, 0
}
