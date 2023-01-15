package blockmanager

import (
	"bytes"
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blocks"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/fileserver"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gnetwork"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type BlockManager struct {
	consensus      *consensus.Consensus
	fsm            *consensus.BlockMachine
	block          *blocks.Blocks
	blockContainer *store.BlockContainer
	fileServer     *fileserver.FileServer
	owner          *gcrypto.GhostAddress
	localIpAddr    *ptypes.GhostIp

	packetSqHandler map[packets.PacketThirdType]func(*packets.Header, *net.UDPAddr) []p2p.PacketHeaderInfo
	packetCqHandler map[packets.PacketThirdType]func(*packets.Header, *net.UDPAddr)
}

func NewBlockManager(con *consensus.Consensus,
	fsm *consensus.BlockMachine,
	block *blocks.Blocks,
	blockContainer *store.BlockContainer,
	master *gnetwork.MasterNetwork,
	fileServer *fileserver.FileServer,
	user *gcrypto.GhostAddress,
	myIpAddr *ptypes.GhostIp) *BlockManager {

	blockMgr := &BlockManager{
		consensus:       con,
		fsm:             fsm,
		block:           block,
		blockContainer:  blockContainer,
		fileServer:      fileServer,
		owner:           user,
		localIpAddr:     myIpAddr,
		packetSqHandler: make(map[packets.PacketThirdType]func(*packets.Header, *net.UDPAddr) []p2p.PacketHeaderInfo),
		packetCqHandler: make(map[packets.PacketThirdType]func(*packets.Header, *net.UDPAddr)),
	}
	blockMgr.InitHandler(master)
	return blockMgr
}

func (blockMgr *BlockManager) BlockSync() bool {
	if blockMgr.fsm.CheckAcceptNewBlock() == false {
		return true
	}

	if result, _ := blockMgr.consensus.CheckTriggerNewBlock(); result == true {
		blockMgr.TriggerNewBlock()
	} else {
		blockMgr.BroadcastBlockChainNotification()
	}
	return true
}

func (blockMgr *BlockManager) TriggerNewBlock() {
	result, triggerTxCount := blockMgr.consensus.CheckTriggerNewBlock()
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

func (blockMgr *BlockManager) NewBlockEvent(filename string, addr *ptypes.GhostIp) {
	if blockMgr.fileServer.CheckFileExist(filename) == false {
		blockMgr.fileServer.SendGetFileInfo(addr, filename, blockMgr.DownloadNewBlock, nil)
	}
}

func (blockMgr *BlockManager) DownloadNewBlock(obj *fileserver.FileObject, context interface{}) {
	byteBuf := bytes.NewBuffer(obj.Buffer)
	newPair := types.PairedBlock{}
	if newPair.Deserialize(byteBuf) == false {
		return
	}

	if blockMgr.TryAddMyBlockChain(&newPair) == true {

	}
}

func (blockMgr *BlockManager) TryAddMyBlockChain(pairedBlock *types.PairedBlock) bool {
	localHeight := blockMgr.GetHeightestBlock()
	if localHeight+1 == pairedBlock.BlockId() {
		if blockMgr.block.BlockValidation(pairedBlock, nil) == true &&
			blockMgr.consensus.CheckMinimumTxCount(pairedBlock) == true {
			blockMgr.blockContainer.InsertBlock(pairedBlock)
			return true
		}
	} else if localHeight+1 < pairedBlock.BlockId() {
		// trigger get neighbor block
	}

	return false
}
