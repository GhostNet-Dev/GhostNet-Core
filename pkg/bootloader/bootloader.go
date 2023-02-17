package bootloader

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
)

type BootLoader struct {
	udp           *p2p.UdpServer
	packetFactory *p2p.PacketFactory
	config        *gconfig.GConfig
	db            *LiteStore
	wallet        *LoadWallet
	conn          *ConnectMaster
	genesis       *LoadGenesis
}

var Tables []string

func NewBootLoader(tables []string, udp *p2p.UdpServer, packetFactory *p2p.PacketFactory, config *gconfig.GConfig,
	db *LiteStore, wallet *LoadWallet, gen *LoadGenesis) *BootLoader {
	Tables = tables

	return &BootLoader{
		udp:           udp,
		packetFactory: packetFactory,
		config:        config,
		db:            db,
		wallet:        wallet,
		genesis:       gen,
	}
}

func (b *BootLoader) BootLoading(config *gconfig.GConfig) *gcrypto.Wallet {
	// Load Wallet
	w, err := b.wallet.OpenWallet(config.Username, config.Password)
	if err != nil {
		w = b.wallet.CreateWallet(config.Username, config.Password)
		b.wallet.SaveWallet(w, config.Password)
	}

	// if Creator, need not to connect other master
	// Load Creator
	creators := b.genesis.CreatorList()
	if creator, exist := creators[config.Username]; exist {
		if creatorAddr, err := b.genesis.LoadCreatorKeyFile(creator.Nickname,
			creator.PubKey, config.Password); err == nil {
			w = gcrypto.NewWallet(creator.Nickname, creatorAddr, &ptypes.GhostIp{Ip: config.Ip, Port: config.Port}, nil)
		}
		return w
	}

	// connect to master Node
	b.conn = NewConnectMaster(Tables[0], b.db, b.packetFactory, b.udp, w)
	masterNode := b.conn.LoadMasterNode()

	if masterNode == nil {
		b.conn.RequestMasterNodesToAdam()
		if timeout := b.conn.WaitEvent(); timeout {
			log.Println("could not connect to adam node")
			return nil
		}
		masterNode = b.conn.LoadMasterNode()
		if masterNode == nil {
			log.Println("could not found masterNode")
			return nil
		}
	}

	b.conn.ConnectToMaster(masterNode)
	if timeout := b.conn.WaitEvent(); timeout {
		log.Println("could not connect to masterNode = ", masterNode.Nickname)
	}

	w.SetMasterNode(masterNode)
	return w
}
