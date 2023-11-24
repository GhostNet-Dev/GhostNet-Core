package blockmanager

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blocks"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/cloudservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus/states"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/fileservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gnetwork"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/packets"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/txs"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"google.golang.org/protobuf/proto"
)

type BlockManager struct {
	BlockTick        int
	consensus        *consensus.Consensus
	fsm              *states.BlockMachine
	block            *blocks.Blocks
	tXs              *txs.TXs
	blockContainer   *store.BlockContainer
	accountContainer *store.AccountContainer
	master           *gnetwork.MasterNetwork
	fileService      *fileservice.FileService
	cloud            *cloudservice.CloudService
	owner            *gcrypto.Wallet
	localIpAddr      *ptypes.GhostIp
	glog             *glogger.GLogger

	packetSqHandler map[packets.PacketThirdType]func(*packets.Header, *net.UDPAddr) []p2p.ResponseHeaderInfo
	packetCqHandler map[packets.PacketThirdType]func(*packets.Header, *net.UDPAddr)
	callback        func(bool)
	newBlockTrigger bool
	blockControl    bool
}

func NewBlockManager(blockTick int, con *consensus.Consensus,
	fsm *states.BlockMachine,
	block *blocks.Blocks,
	tXs *txs.TXs,
	blockContainer *store.BlockContainer,
	accountContainer *store.AccountContainer,
	master *gnetwork.MasterNetwork,
	fileService *fileservice.FileService,
	cloud *cloudservice.CloudService,
	user *gcrypto.Wallet,
	glog *glogger.GLogger) *BlockManager {

	blockMgr := &BlockManager{
		BlockTick:        blockTick,
		consensus:        con,
		fsm:              fsm,
		block:            block,
		tXs:              tXs,
		blockContainer:   blockContainer,
		accountContainer: accountContainer,
		master:           master,
		fileService:      fileService,
		cloud:            cloud,
		owner:            user,
		glog:             glog,
		packetSqHandler:  make(map[packets.PacketThirdType]func(*packets.Header, *net.UDPAddr) []p2p.ResponseHeaderInfo),
		packetCqHandler:  make(map[packets.PacketThirdType]func(*packets.Header, *net.UDPAddr)),
		newBlockTrigger:  false,
	}

	fsm.BlockServer = blockMgr
	blockMgr.InitHandler(master)
	blockContainer.RegisterBlockEvent(
		func(pairedBlock *types.PairedBlock) {
			blockMgr.SaveExtraInformation(pairedBlock)
		},
		func(pairedBlock *types.PairedBlock) {
			blockMgr.UnsaveExtraInformation(pairedBlock)
		})

	blockMgr.BlockPlay()
	return blockMgr
}

func (blockMgr *BlockManager) BlockServer() {
	blockMgr.glog.DebugOutput(blockMgr, "Block Server Start.", glogger.Default)
	blockMgr.fsm.CheckBlock()
	blockMgr.BlockSync()
	ticker := time.NewTicker(time.Second * time.Duration(blockMgr.BlockTick))
	defer ticker.Stop()

	for range ticker.C {
		blockMgr.BlockSync()
		if !blockMgr.blockControl {
			return
		}
	}
}

func (blockMgr *BlockManager) BlockPlay() bool {
	blockMgr.blockControl = true
	return true
}

func (blockMgr *BlockManager) BlockStop() bool {
	blockMgr.blockControl = false
	return true
}

func (blockMgr *BlockManager) BlockPause() bool {
	blockMgr.blockControl = false
	return true
}

func (blockMgr *BlockManager) BlockSync() bool {
	if !blockMgr.fsm.CheckAcceptNewBlock() {
		return true
	}

	if result, _ := blockMgr.consensus.CheckTriggerNewBlock(); result {
		blockMgr.TriggerNewBlock()
	} else {
		blockMgr.BroadcastBlockChainNotification()
	}
	return true
}

func (blockMgr *BlockManager) TriggerNewBlock() {
	result, triggerTxCount := blockMgr.consensus.CheckTriggerNewBlock()
	if !blockMgr.fsm.CheckAcceptNewBlock() || !result {
		return
	}
	if !blockMgr.newBlockTrigger {
		blockMgr.newBlockTrigger = true
		defer func() { blockMgr.newBlockTrigger = false }()
	} else {
		return
	}
	blockMgr.glog.DebugOutput(blockMgr, "Trigger New Block", glogger.BlockConsensus)
	// miner와 creator는 동일하게 한다. 즉 creator만 mining을 할 수 있다.
	newPairBlock := blockMgr.block.MakeNewBlock(blockMgr.owner.GetGhostAddress(), blockMgr.owner.GetGhostAddress().Get160PubKey(), triggerTxCount)
	if newPairBlock == nil {
		blockMgr.glog.DebugOutput(blockMgr, "Fail to Make New Block", glogger.BlockConsensus)
		return
	}
	blockFilename := newPairBlock.GetBlockFilename()
	blockMgr.fileService.CreateFile(blockFilename, newPairBlock.SerializeToByte(), nil, nil)
	blockMgr.glog.DebugOutput(blockMgr, fmt.Sprint("Create Block Id = ", newPairBlock.BlockId()), glogger.BlockConsensus)

	sq := packets.NewBlockSq{
		Master:        p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), nil, 0, blockMgr.localIpAddr),
		BlockFilename: blockFilename,
	}
	sendData, err := proto.Marshal(&sq)
	if err != nil {
		log.Fatal(err)
	}
	headerInfo := &p2p.ResponseHeaderInfo{
		PacketType: packets.PacketType_MasterNetwork,
		SecondType: packets.PacketSecondType_BlockChain,
		ThirdType:  packets.PacketThirdType_NewBlock,
		RequestId:  sq.Master.GetRequestId(),
		PacketData: sendData,
		SqFlag:     true,
	}
	blockMgr.master.SendToMasterNodeGrpSq(packets.RoutingType_Flooding, gnetwork.DefaultTreeLevel, headerInfo)
}

func (blockMgr *BlockManager) GetHeightestBlock() uint32 {
	return blockMgr.blockContainer.BlockHeight()
}

func (blockMgr *BlockManager) PrepareSendBlock(blockId uint32) (string, bool) {
	pairedBlock := blockMgr.blockContainer.GetBlock(blockId)
	if pairedBlock == nil {
		return "", false
	}

	blockFilename := pairedBlock.GetBlockFilename()
	blockMgr.fileService.CreateFile(blockFilename, pairedBlock.SerializeToByte(), nil, nil)

	return blockFilename, true
}

func (blockMgr *BlockManager) DownloadDataTransaction(txByte []byte, dataTxByte []byte) bool {
	blockMgr.glog.DebugOutput(blockMgr, "Tx with data Download Complete", glogger.Default)
	tx := &types.GhostTransaction{}
	if !tx.Deserialize(bytes.NewBuffer(txByte)).Result() {
		return false
	}
	dataTx := &types.GhostDataTransaction{}
	if !dataTx.Deserialize(bytes.NewBuffer(dataTxByte)).Result() {
		return false
	}
	if blockMgr.blockContainer.TxContainer.CheckExistCandidateTxId(tx.TxId) {
		blockMgr.glog.DebugOutput(blockMgr, "Already candidate Tx", glogger.Default)
		return false
	}
	if !blockMgr.tXs.TransactionValidation(tx, dataTx, blockMgr.blockContainer.TxContainer).Result() {
		blockMgr.glog.DebugOutput(blockMgr, "Tx with data Validation Fail", glogger.Default)
		return false
	}
	blockMgr.blockContainer.TxContainer.SaveCandidateTx(tx)
	blockMgr.blockContainer.TxContainer.SaveCandidateDataTx(dataTx)
	blockMgr.TriggerNewBlock()
	return true
}

func (blockMgr *BlockManager) DownloadTransaction(obj *fileservice.FileObject, context interface{}) bool {
	blockMgr.glog.DebugOutput(blockMgr, "Tx Download Complete", glogger.Default)
	tx := &types.GhostTransaction{}
	if !tx.Deserialize(bytes.NewBuffer(obj.Buffer)).Result() {
		return false
	}
	if blockMgr.blockContainer.TxContainer.CheckExistCandidateTxId(tx.TxId) {
		blockMgr.glog.DebugOutput(blockMgr, "Already candidate Tx", glogger.Default)
		return false
	}
	if !blockMgr.tXs.TransactionValidation(tx, nil, blockMgr.blockContainer.TxContainer).Result() {
		blockMgr.glog.DebugOutput(blockMgr, "Tx Validation Fail", glogger.Default)
		return false
	}

	if !blockMgr.checkExistFSRoot(tx) {
		return false
	}

	blockMgr.blockContainer.TxContainer.SaveCandidateTx(tx)
	blockMgr.TriggerNewBlock()
	return true
}

func (blockMgr *BlockManager) checkExistFSRoot(tx *types.GhostTransaction) bool {
	// TODO: Check from Block.db
	for _, output := range tx.Body.Vout {
		if output.Type == types.TxTypeFSRoot {
			nick := output.ScriptEx
			if blockMgr.blockContainer.TxContainer.CheckExistFsRoot(nick) {
				return false
			}
		}
	}
	return true
}

func (blockMgr *BlockManager) RequestCheckExistFsRoot(nickname []byte, callback func(bool)) {
	sq := packets.CheckRootFsSq{
		Master:   p2p.MakeMasterPacket(blockMgr.owner.GetPubAddress(), nil, 0, blockMgr.localIpAddr),
		Nickname: nickname,
	}
	sendData, err := proto.Marshal(&sq)
	if err != nil {
		log.Fatal(err)
	}
	blockMgr.callback = callback
	blockMgr.master.SendToMasterNodeSq(packets.PacketThirdType_CheckRootFs,
		blockMgr.owner.GetMasterNode().PubKey, sendData, sq.Master.GetRequestId())
}

func (blockMgr *BlockManager) SaveExtraInformation(pairedBlock *types.PairedBlock) bool {
	return true
}

func (blockMgr *BlockManager) UnsaveExtraInformation(pairedBlock *types.PairedBlock) {
}

func (blockMgr *BlockManager) DownloadBlock(obj *fileservice.FileObject, pubKey string) bool {
	pair := &types.PairedBlock{}
	if !pair.Deserialize(bytes.NewBuffer(obj.Buffer)) {
		blockMgr.fileService.DeleteFile(obj.Filename)
		return false
	}
	blockMgr.fsm.State().RecvBlock(pair, pubKey)

	return true
}

func (blockMgr *BlockManager) DownloadNewBlock(obj *fileservice.FileObject, context interface{}) {
	byteBuf := bytes.NewBuffer(obj.Buffer)
	newPair := types.PairedBlock{}
	if !newPair.Deserialize(byteBuf) {
		blockMgr.fileService.DeleteFile(obj.Filename)
		return
	}

	blockMgr.TryAddMyBlockChain(&newPair)
}

func (blockMgr *BlockManager) TryAddMyBlockChain(pairedBlock *types.PairedBlock) bool {
	localHeight := blockMgr.GetHeightestBlock()
	if localHeight+1 == pairedBlock.BlockId() {
		if blockMgr.block.BlockValidation(pairedBlock, nil) &&
			blockMgr.consensus.CheckMinimumTxCount(pairedBlock) {
			blockMgr.blockContainer.InsertBlock(pairedBlock)
			return true
		}
	} else if localHeight+1 < pairedBlock.BlockId() {
		// trigger get neighbor block
		blockMgr.fsm.State().Rebuild()
		return true
	}

	return false
}
