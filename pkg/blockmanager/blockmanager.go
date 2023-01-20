package blockmanager

import (
	"bytes"
	"fmt"
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blocks"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/cloudservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus/states"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/fileservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gnetwork"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/txs"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"github.com/btcsuite/btcutil/base58"
)

type BlockManager struct {
	consensus      *consensus.Consensus
	fsm            *states.BlockMachine
	block          *blocks.Blocks
	tXs            *txs.TXs
	blockContainer *store.BlockContainer
	fileService    *fileservice.FileService
	cloud          *cloudservice.CloudService
	owner          *gcrypto.GhostAddress
	localIpAddr    *ptypes.GhostIp

	packetSqHandler map[packets.PacketThirdType]func(*packets.Header, *net.UDPAddr) []p2p.PacketHeaderInfo
	packetCqHandler map[packets.PacketThirdType]func(*packets.Header, *net.UDPAddr)
}

func NewBlockManager(con *consensus.Consensus,
	fsm *states.BlockMachine,
	block *blocks.Blocks,
	tXs *txs.TXs,
	blockContainer *store.BlockContainer,
	master *gnetwork.MasterNetwork,
	fileService *fileservice.FileService,
	cloud *cloudservice.CloudService,
	user *gcrypto.GhostAddress,
	myIpAddr *ptypes.GhostIp) *BlockManager {

	blockMgr := &BlockManager{
		consensus:       con,
		fsm:             fsm,
		block:           block,
		tXs:             tXs,
		blockContainer:  blockContainer,
		fileService:     fileService,
		cloud:           cloud,
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

func (blockMgr *BlockManager) PrepareSendBlock(blockId uint32) (string, bool) {
	pairedBlock := blockMgr.blockContainer.GetBlock(blockId)
	if pairedBlock == nil {
		return "", false
	}

	blockFilename := fmt.Sprint(blockId, "@", base58.Encode(pairedBlock.Block.GetHashKey()), ".ghost")
	blockMgr.fileService.CreateFile(blockFilename, pairedBlock.SerializeToByte(), nil, nil)

	return blockFilename, true
}

func (blockMgr *BlockManager) DownloadDataTransaction(obj *fileservice.FileObject, context interface{}) {
}

func (blockMgr *BlockManager) DownloadTransaction(obj *fileservice.FileObject, context interface{}) bool {
	tx := &types.GhostTransaction{}
	if tx.Deserialize(bytes.NewBuffer(obj.Buffer)).Result() == false {
		return false
	}
	if blockMgr.tXs.TransactionValidation(tx, nil, blockMgr.blockContainer.TxContainer).Result() == false {
		return false
	}
	blockMgr.blockContainer.TxContainer.SaveCandidateTx(tx)
	return true
}

func (blockMgr *BlockManager) DownloadBlock(obj *fileservice.FileObject, context interface{}) {
}

func (blockMgr *BlockManager) DownloadNewBlock(obj *fileservice.FileObject, context interface{}) {
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
