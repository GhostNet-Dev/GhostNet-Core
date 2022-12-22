package mainserver

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blocks"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
)

type BlockServer struct {
	con   *consensus.Consensus
	fsm   *consensus.BlockMachine
	block *blocks.Blocks
	user  *gcrypto.GhostAddress
}

func NewBlockServer(con *consensus.Consensus,
	fsm *consensus.BlockMachine,
	block *blocks.Blocks,
	user *gcrypto.GhostAddress) *BlockServer {
	return &BlockServer{
		con:   con,
		fsm:   fsm,
		block: block,
		user:  user,
	}
}

func (server *BlockServer) BlockSync() bool {
	if server.fsm.CheckAcceptNewBlock() == false {
		return true
	}

	if result, _ := server.con.CheckTriggerNewBlock(); result == true {
		server.TriggerNewBlock()
	} else {
		server.BroadcastBlockChainNotification()
	}
	return true
}

func (server *BlockServer) TriggerNewBlock() {
	result, triggerTxCount := server.con.CheckTriggerNewBlock()
	if server.fsm.CheckAcceptNewBlock() == false || result == false {
		return
	}
	// miner와 creator는 동일하게 한다. 즉 creator만 mining을 할 수 있다.
	newPairBlock := server.block.MakeNewBlock(server.user, server.user.Get160PubKey(), triggerTxCount)
	if newPairBlock == nil {
		return
	}
}
