package blockmanager

import (
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blocks"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/fileserver"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gnetwork"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gsql"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
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
	nickname = "test"

	gScript        = gvm.NewGScript()
	gVm            = gvm.NewGVM()
	blockContainer = store.NewBlockContainer()
	Txs            = txs.NewTXs(gScript, blockContainer, gVm)
	block          = blocks.NewBlocks(blockContainer, Txs, 1)
	config         = gconfig.DefaultConfig()
	packetFactory  = p2p.NewPacketFactory()
	udp            = p2p.NewUdpServer(ipAddr.Ip, ipAddr.Port)
	con            = consensus.NewConsensus(blockContainer)
	fsm            = consensus.NewBlockMachine(blockContainer)
	tTreeMap       = gnetwork.NewTrieTreeMap(Miner.GetPubAddress(), gsql.NewAccountSql("sqlite3"))
	master         = gnetwork.NewMasterNode(nickname, Miner, ipAddr, config, packetFactory, udp, blockContainer, tTreeMap)
	fileServer     = fileserver.NewFileServer(udp, packetFactory, Miner, ipAddr, "./")
	blockServer    = NewBlockManager(con, fsm, block, blockContainer, master, fileServer, Miner, ipAddr)
)

func init() {
	blockContainer.BlockContainerOpen("../../db.sqlite3.sql", "./")
}

func TestStartServer(t *testing.T) {
	blockServer.BlockSync()

}
