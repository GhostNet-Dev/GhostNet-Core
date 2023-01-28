package factory

import (
	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blockmanager"
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

type DefaultFactory struct {
	Con            *consensus.Consensus
	Fsm            *states.BlockMachine
	Block          *blocks.Blocks
	Txs            *txs.TXs
	BlockContainer *store.BlockContainer
	Master         *gnetwork.MasterNetwork
	Account        *gnetwork.GhostAccount
	TTreeMap       *gnetwork.TrieTreeMap
	FileService    *fileservice.FileService
	Cloud          *cloudservice.CloudService
	BlockServer    *blockmanager.BlockManager
	GScript        *gvm.GScript
	Gvm            *gvm.GVM
	Owner          *gcrypto.GhostAddress
	LocalIpAddr    *ptypes.GhostIp

	networkFactory *NetworkFactory
}

type NetworkFactory struct {
	packetFactory *p2p.PacketFactory
	udp           *p2p.UdpServer
}

func NewNetworkFactory(config *gconfig.GConfig) *NetworkFactory {

	return &NetworkFactory{
		packetFactory: p2p.NewPacketFactory(),
		udp:           p2p.NewUdpServer(config.Ip, config.Port),
	}
}

func NewDefaultFactory(networkFactory *NetworkFactory,
	user *gcrypto.Wallet, config *gconfig.GConfig) *DefaultFactory {
	ghostIp := &ptypes.GhostIp{
		Ip:   config.Ip,
		Port: config.Port,
	}

	factory := &DefaultFactory{
		networkFactory: networkFactory,
	}

	factory.GScript = gvm.NewGScript()
	factory.Gvm = gvm.NewGVM()
	factory.BlockContainer = store.NewBlockContainer(config.DbName)

	factory.Account = gnetwork.NewGhostAccount()
	factory.TTreeMap = gnetwork.NewTrieTreeMap(user.GetPubAddress(), factory.Account)
	factory.Master = gnetwork.NewMasterNode(user, ghostIp, config, networkFactory.packetFactory,
		networkFactory.udp, factory.BlockContainer, factory.Account, factory.TTreeMap)

	factory.FileService = fileservice.NewFileServer(networkFactory.udp, networkFactory.packetFactory,
		user.GetGhostAddress(), ghostIp, config.FilePath)
	factory.Cloud = cloudservice.NewCloudService(factory.FileService, factory.TTreeMap)
	factory.Txs = txs.NewTXs(factory.GScript, factory.BlockContainer, factory.Gvm)
	factory.Block = blocks.NewBlocks(factory.BlockContainer, factory.Txs, 1)
	factory.Con = consensus.NewConsensus(factory.BlockContainer, factory.Block)
	factory.Fsm = states.NewBlockMachine(factory.BlockContainer, factory.Con)
	factory.BlockServer = blockmanager.NewBlockManager(factory.Con, factory.Fsm, factory.Block,
		factory.Txs, factory.BlockContainer, factory.Master, factory.FileService, factory.Cloud, user.GetGhostAddress(), ghostIp)

	return factory
}
