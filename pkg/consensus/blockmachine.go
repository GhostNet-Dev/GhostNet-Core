package consensus

import "github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus/states"

type BlockMachine struct {
	idle               states.BlockState
	miningState        states.BlockState
	getHeightestState  states.BlockState
	verificationState  states.BlockState
	downloadCheckState states.BlockState
	getBlockState      states.BlockState
	blockMergeState    states.BlockState
	recheckState       states.BlockState

	currentState BlockState
}

func NewBlockMachine() *BlockMachine {
	fsm := &BlockMachine{}
	idleState := &states.IdleState{blockMachine: fsm}
	miningState := &states.MiningState{blockMachine: fsm}
	getHeightestState := &states.GetHeigtestState{blockMachine: fsm}
	verificationState := &states.VerificationState{blockMachine: fsm}
	downloadCheckState := &states.DownloadCheckState{blockMachine: fsm}
	getBlockState := &states.GetBlockState{blockMachine: fsm}
	blockMergeState := &states.BlockMergeState{blockMachine: fsm}
	recheckState := &states.RecheckState{blockMachine: fsm}

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

func (m *BlockMachine) setState(s BlockState) {
	m.currentState = s
}
