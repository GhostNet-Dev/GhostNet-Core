package consensus

import "github.com/GhostNet-Dev/GhostNet-Core/pkg/store"

type BlockMachine struct {
	idle               BlockState
	miningState        BlockState
	getHeightestState  BlockState
	verificationState  BlockState
	downloadCheckState BlockState
	getBlockState      BlockState
	blockMergeState    BlockState
	recheckState       BlockState

	currentState   BlockState
	blockContainer *store.BlockContainer
}

func NewBlockMachine(b *store.BlockContainer) *BlockMachine {
	fsm := &BlockMachine{}
	fsm.blockContainer = b

	idleState := &IdleState{blockMachine: fsm}
	miningState := &MiningState{blockMachine: fsm}
	getHeightestState := &GetHeigtestState{blockMachine: fsm}
	verificationState := &VerificationState{blockMachine: fsm}
	downloadCheckState := &DownloadCheckState{blockMachine: fsm}
	getBlockState := &GetBlockState{blockMachine: fsm}
	blockMergeState := &BlockMergeState{blockMachine: fsm}
	recheckState := &RecheckState{blockMachine: fsm}

	fsm.setState(idleState)
	fsm.idle = idleState
	fsm.miningState = miningState
	fsm.getHeightestState = getHeightestState
	fsm.verificationState = verificationState
	fsm.downloadCheckState = downloadCheckState
	fsm.getBlockState = getBlockState
	fsm.blockMergeState = blockMergeState
	fsm.recheckState = recheckState

	return fsm
}

func (fsm *BlockMachine) CheckAcceptNewBlock() bool {
	return fsm.currentState == fsm.idle || fsm.currentState == fsm.miningState
}

func (fsm *BlockMachine) setState(s BlockState) {
	fsm.currentState = s
}
