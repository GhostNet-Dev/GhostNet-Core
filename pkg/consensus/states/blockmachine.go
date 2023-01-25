package states

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type IBlockServer interface {
	BroadcastBlockChainNotification()
	MiningStart()
	MiningStop()
	RequestGetBlock(pubKey string, blockIdx uint32)
	RequestGetBlockHash(pubKey string, currBlockHeight uint32)
	MergeErrorNotification(pubKey string, result bool)
	LocalBlockDbValidation() bool
	BlockServerInitStart()
	CheckHeightForRebuild(uint32) bool
}

type BlockMachine struct {
	idle               IBlockState
	miningState        IBlockState
	getHeightestState  IBlockState
	verificationState  IBlockState
	downloadCheckState IBlockState
	blockMergeState    IBlockState
	recheckState       IBlockState

	currentState   IBlockState
	blockContainer *store.BlockContainer
	blockServer    IBlockServer
	con            *consensus.Consensus

	heightestCandidataPool []string
	startBlockId           uint32
	targetHeight           uint32
}

func NewBlockMachine(b *store.BlockContainer, con *consensus.Consensus) *BlockMachine {
	fsm := &BlockMachine{
		blockContainer: b,
		con:            con,
	}

	idleState := &IdleState{blockMachine: fsm}
	miningState := &MiningState{blockMachine: fsm}
	getHeightestState := &GetHeigtestState{blockMachine: fsm}
	verificationState := &VerificationState{blockMachine: fsm}
	downloadCheckState := &DownloadCheckState{blockMachine: fsm}
	blockMergeState := &BlockMergeState{blockMachine: fsm}
	recheckState := &RecheckState{blockMachine: fsm}

	fsm.setState(idleState)
	fsm.idle = idleState
	fsm.miningState = miningState
	fsm.getHeightestState = getHeightestState
	fsm.verificationState = verificationState
	fsm.downloadCheckState = downloadCheckState
	fsm.blockMergeState = blockMergeState
	fsm.recheckState = recheckState

	return fsm
}

func (fsm *BlockMachine) SetServer(blockServer IBlockServer) {
	fsm.blockServer = blockServer
}

func (fsm *BlockMachine) CheckAcceptNewBlock() bool {
	return fsm.currentState == fsm.idle || fsm.currentState == fsm.miningState
}

func (fsm *BlockMachine) setState(s IBlockState) {
	s.Initialize()
	fsm.currentState = s
}

func (fsm *BlockMachine) State() IBlockState {
	return fsm.currentState
}

func (fsm *BlockMachine) SetTargetHeight(targetHeight uint32) {
	fsm.targetHeight = targetHeight
}

func (fsm *BlockMachine) GetTargetHeight() uint32 {
	return fsm.targetHeight
}

func (fsm *BlockMachine) SetStartBlockId(targetBlockId uint32) {
	fsm.startBlockId = targetBlockId
}

func (fsm *BlockMachine) GetStartBlockId() uint32 {
	return fsm.startBlockId
}
func (fsm *BlockMachine) SetHeighestCandidatePool(pool []string) {
	fsm.heightestCandidataPool = pool
}

func (fsm *BlockMachine) GetHeighestCandidatePool() []string {
	return fsm.heightestCandidataPool
}

func (fsm *BlockMachine) LoadHashFromTempDb(blockId uint32) []byte {
	return fsm.con.LoadHashFromTempDb(blockId)
}

func (fsm *BlockMachine) CheckAndSave(candidatePair *types.PairedBlock) bool {
	return fsm.con.CheckAndSave(candidatePair)
}

func (fsm *BlockMachine) MergeExecute() {
	fsm.con.MergeExecute(fsm.startBlockId, fsm.targetHeight)
}

func (fsm *BlockMachine) LocalBlockCheckProcess() bool {
	return fsm.con.LocalBlockCheckProcess()
}
