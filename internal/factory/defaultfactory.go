package factory

import (
	"path"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blockmanager"
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

type DefaultFactory struct {
	Con              *consensus.Consensus
	Fsm              *states.BlockMachine
	Block            *blocks.Blocks
	Txs              *txs.TXs
	BlockContainer   *store.BlockContainer
	AccountContainer *store.AccountContainer
	Master           *gnetwork.MasterNetwork
	Account          *gnetwork.GhostAccount
	TTreeMap         *gnetwork.TrieTreeMap
	FileService      *fileservice.FileService
	Cloud            *cloudservice.CloudService
	BlockServer      *blockmanager.BlockManager
	GScript          *gvm.GScript
	Gvm              *gvm.GVM
	UserWallet       *gcrypto.Wallet
	Owner            *gcrypto.GhostAddress
	LocalIpAddr      *ptypes.GhostIp

	networkFactory *NetworkFactory
	config         *gconfig.GConfig
	glog           *glogger.GLogger
}

type NetworkFactory struct {
	PacketFactory *p2p.PacketFactory
	Udp           *p2p.UdpServer
	glog          *glogger.GLogger
}

func NewNetworkFactory(config *gconfig.GConfig, glog *glogger.GLogger) *NetworkFactory {
	packetFactory := p2p.NewPacketFactory()

	return &NetworkFactory{
		PacketFactory: packetFactory,
		Udp:           p2p.NewUdpServer(config.Ip, config.Port, packetFactory, glog),
		glog:          glog,
	}
}

func NewDefaultFactory(networkFactory *NetworkFactory, bootFactory *BootFactory,
	user *gcrypto.Wallet, config *gconfig.GConfig, glog *glogger.GLogger) *DefaultFactory {
	ghostIp := &ptypes.GhostIp{
		Ip:   config.Ip,
		Port: config.Port,
	}

	factory := &DefaultFactory{
		networkFactory: networkFactory,
		config:         config,
		UserWallet:     user,
	}

	factory.glog = glog
	factory.GScript = gvm.NewGScript()
	factory.Gvm = gvm.NewGVM()
	factory.BlockContainer = store.NewBlockContainer(config.DbName)
	factory.AccountContainer = store.NewAccountContainer(bootFactory.Db)

	factory.Account = gnetwork.NewGhostAccount(bootFactory.Db)
	factory.TTreeMap = gnetwork.NewTrieTreeMap(user.GetPubAddress(), factory.Account)
	factory.Master = gnetwork.NewMasterNode(user, ghostIp, config, networkFactory.PacketFactory,
		networkFactory.Udp, factory.BlockContainer, factory.Account, factory.TTreeMap)

	factory.FileService = fileservice.NewFileServer(networkFactory.Udp, networkFactory.PacketFactory,
		user.GetGhostAddress(), ghostIp, config.FilePath)
	factory.Cloud = cloudservice.NewCloudService(factory.FileService, factory.TTreeMap)
	factory.Txs = txs.NewTXs(factory.GScript, factory.BlockContainer, factory.Gvm)
	factory.Block = blocks.NewBlocks(factory.BlockContainer, factory.Txs, 1)
	factory.Con = consensus.NewConsensus(factory.BlockContainer, factory.Block, factory.glog)
	factory.Fsm = states.NewBlockMachine(factory.BlockContainer, factory.Con, factory.glog)
	factory.BlockServer = blockmanager.NewBlockManager(factory.Con, factory.Fsm, factory.Block,
		factory.Txs, factory.BlockContainer, factory.AccountContainer, factory.Master, factory.FileService, factory.Cloud, user.GetGhostAddress(), ghostIp, factory.glog)

	return factory
}

func (factory *DefaultFactory) FactoryOpen() {
	schemeFilePath := path.Join(factory.config.DbSchemePath, factory.config.DbSchemeFile)
	factory.BlockContainer.BlockContainerOpen(schemeFilePath, factory.config.SqlPath)
	factory.BlockContainer.GenesisBlockChecker(store.GenesisBlock())
	factory.Master.RegisterMyMasterNode(factory.UserWallet.GetMasterNode())
}
