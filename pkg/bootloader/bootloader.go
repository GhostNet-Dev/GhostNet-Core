package bootloader

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
)

type BootLoader struct {
	udp           *p2p.UdpServer
	packetFactory *p2p.PacketFactory
	config        *gconfig.GConfig
	db            *LiteStore
	wallet        *LoadWallet
	conn          *ConnectMaster
}

var Tables []string

func NewBootLoader(tables []string, udp *p2p.UdpServer, packetFactory *p2p.PacketFactory, config *gconfig.GConfig,
	db *LiteStore, wallet *LoadWallet) *BootLoader {
	Tables = tables

	return &BootLoader{
		udp:           udp,
		packetFactory: packetFactory,
		config:        config,
		db:            db,
		wallet:        wallet,
	}
}

func (b *BootLoader) BootLoading(config *gconfig.GConfig) *gcrypto.Wallet {
	w, err := b.wallet.OpenWallet(config.Username, config.Password)
	if err != nil {
		w = b.wallet.CreateWallet(config.Username, config.Password)
		b.wallet.SaveWallet(w, config.Password)
	}

	b.conn = NewConnectMaster(Tables[0], b.db, b.packetFactory, b.udp, w)

	if config.StandaloneMode {
		return w
	}

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
