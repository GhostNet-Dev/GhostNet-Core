package blockmanager

import (
	"log"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blocks"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/cloudservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/consensus/states"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/fileservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
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

	TestTables = []string{"nodes", "wallet"}
	liteStore  = store.NewLiteStore("./", "litestore.db", TestTables)

	glog             = glogger.NewGLogger(0)
	config           = gconfig.NewDefaultConfig()
	gScript          = gvm.NewGScript()
	gVm              = gvm.NewGVM()
	blockContainer   = store.NewBlockContainer("sqlite3")
	accountContainer = store.NewBcAccountContainer(liteStore)
	packetFactory    = p2p.NewPacketFactory()
	udp              = p2p.NewUdpServer(ghostIp.Ip, ghostIp.Port, packetFactory, glog)

	account  = gnetwork.NewGhostAccount(liteStore)
	tTreeMap = gnetwork.NewTrieTreeMap(Miner.GetPubAddress(), account)
	master   = gnetwork.NewMasterNode(1, w, ghostIp, packetFactory, udp, blockContainer, account, tTreeMap)

	fileService = fileservice.NewFileServer(udp, packetFactory, Miner, ghostIp, "./", glog)
	cloud       = cloudservice.NewCloudService(fileService, tTreeMap, glog)
	tXs         = txs.NewTXs(gScript, blockContainer, gVm)
	block       = blocks.NewBlocks(blockContainer, tXs, 1)
	con         = consensus.NewConsensus(blockContainer, block, glog)
	fsm         = states.NewBlockMachine(blockContainer, con, glog)
	blockServer = NewBlockManager(8, con, fsm, block, tXs, blockContainer, accountContainer, master, fileService, cloud, w, glog)
)

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	fsm.SetServer(blockServer)
}

func TestStartServer(t *testing.T) {
	blockContainer.BlockContainerOpen("./")
	defer blockContainer.Close()
	blockServer.BlockSync()
}
