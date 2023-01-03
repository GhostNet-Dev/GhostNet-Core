package blockmanager

import (
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blocks"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/txs"
)

var (
	Miner  = gcrypto.GenerateKeyPair()
	ipAddr = &ptypes.GhostIp{
		Ip:   "127.0.0.1",
		Port: "8888",
	}

	gScript        = gvm.NewGScript()
	gVm            = gvm.NewGVM()
	blockContainer = store.NewBlockContainer()
	Txs            = txs.NewTXs(gScript, blockContainer, gVm)
	block          = blocks.NewBlocks(blockContainer, Txs, 1)
	con            = consensus.NewConsensus(blockContainer)
	fsm            = consensus.NewBlockMachine(blockContainer)
	blockServer    = NewBlockManager(con, fsm, block, blockContainer, Miner, ipAddr)
)

func init() {
	blockContainer.BlockContainerOpen("../../db.sqlite3.sql", "./")
}

func TestStartServer(t *testing.T) {
	blockServer.BlockSync()

}
