package mainserver

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blocks"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
)

type BlockManager struct {
	con   *consensus.Consensus
	fsm   *consensus.BlockMachine
	block *blocks.Blocks
	user  *gcrypto.GhostAddress
}

func NewBlockManager(con *consensus.Consensus,
	fsm *consensus.BlockMachine,
	block *blocks.Blocks,
	user *gcrypto.GhostAddress) *BlockManager {
	return &BlockManager{
		con:   con,
		fsm:   fsm,
		block: block,
		user:  user,
	}
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
	newPairBlock := blockMgr.block.MakeNewBlock(blockMgr.user, blockMgr.user.Get160PubKey(), triggerTxCount)
	if newPairBlock == nil {
		return
	}
}
