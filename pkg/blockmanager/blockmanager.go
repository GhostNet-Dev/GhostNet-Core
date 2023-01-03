package blockmanager

import (
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blocks"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gnetwork"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
)

type BlockManager struct {
	con            *consensus.Consensus
	fsm            *consensus.BlockMachine
	block          *blocks.Blocks
	blockContainer *store.BlockContainer
	owner          *gcrypto.GhostAddress
	localIpAddr    *ptypes.GhostIp

	packetSqHandler map[packets.PacketThirdType]func(*packets.Header, *net.UDPAddr) []p2p.ResponsePacketInfo
	packetCqHandler map[packets.PacketThirdType]func(*packets.Header, *net.UDPAddr)
}

func NewBlockManager(con *consensus.Consensus,
	fsm *consensus.BlockMachine,
	block *blocks.Blocks,
	blockContainer *store.BlockContainer,
	master *gnetwork.MasterNetwork,
	user *gcrypto.GhostAddress,
	myIpAddr *ptypes.GhostIp) *BlockManager {

	blockMgr := &BlockManager{
		con:            con,
		fsm:            fsm,
		block:          block,
		blockContainer: blockContainer,
		owner:          user,
		localIpAddr:    myIpAddr,
	}
	blockMgr.InitHandler(master)
	return blockMgr
}

func (blockMgr *BlockManager) BlockSync() bool {
	if blockMgr.fsm.CheckAcceptNewBlock() == false {
		return true
	}

	if result, _ := blockMgr.con.CheckTriggerNewBlock(); result == true {
		blockMgr.TriggerNewBlock()
	} else {
		blockMgr.BroadcastBlockChainNotification()
	}
	return true
}

func (blockMgr *BlockManager) TriggerNewBlock() {
	result, triggerTxCount := blockMgr.con.CheckTriggerNewBlock()
	if blockMgr.fsm.CheckAcceptNewBlock() == false || result == false {
		return
	}
	// miner와 creator는 동일하게 한다. 즉 creator만 mining을 할 수 있다.
	newPairBlock := blockMgr.block.MakeNewBlock(blockMgr.owner, blockMgr.owner.Get160PubKey(), triggerTxCount)
	if newPairBlock == nil {
		return
	}
}

func (blockMgr *BlockManager) GetHeightestBlock() uint32 {
	return blockMgr.blockContainer.BlockHeight()
}
