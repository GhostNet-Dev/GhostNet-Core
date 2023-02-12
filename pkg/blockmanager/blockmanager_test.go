package blockmanager

import (
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blocks"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/cloudservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus/states"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/fileservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gnetwork"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/txs"
)

var (
	Miner   = gcrypto.GenerateKeyPair()
	ghostIp = &ptypes.GhostIp{
		Ip:   "127.0.0.1",
		Port: "8888",
	}
	nickname = "test"
	w        = gcrypto.NewWallet(nickname, Miner, ghostIp, nil)

	config         = gconfig.NewDefaultConfig()
	gScript        = gvm.NewGScript()
	gVm            = gvm.NewGVM()
	blockContainer = store.NewBlockContainer("sqlite3")
	packetFactory  = p2p.NewPacketFactory()
	udp            = p2p.NewUdpServer(ghostIp.Ip, ghostIp.Port)

	account  = gnetwork.NewGhostAccount()
	tTreeMap = gnetwork.NewTrieTreeMap(Miner.GetPubAddress(), account)
	master   = gnetwork.NewMasterNode(w, ghostIp, config, packetFactory, udp, blockContainer, account, tTreeMap)

	fileService = fileservice.NewFileServer(udp, packetFactory, Miner, ghostIp, "./")
	cloud       = cloudservice.NewCloudService(fileService, tTreeMap)
	tXs         = txs.NewTXs(gScript, blockContainer, gVm)
	block       = blocks.NewBlocks(blockContainer, tXs, 1)
	con         = consensus.NewConsensus(blockContainer, block)
	fsm         = states.NewBlockMachine(blockContainer, con)
	blockServer = NewBlockManager(con, fsm, block, tXs, blockContainer, master, fileService, cloud, Miner, ghostIp)
)

func init() {
	fsm.SetServer(blockServer)
	blockContainer.BlockContainerOpen("../../db.sqlite3.sql", "./")
}

func TestStartServer(t *testing.T) {
	blockServer.BlockSync()

}
