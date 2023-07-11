package bootloader

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
)

type BootLoader struct {
	udp           *p2p.UdpServer
	packetFactory *p2p.PacketFactory
	db            *store.LiteStore
	wallet        *LoadWallet
	conn          *ConnectMaster
	genesis       *LoadGenesis
}

func NewBootLoader(udp *p2p.UdpServer, packetFactory *p2p.PacketFactory,
	db *store.LiteStore, wallet *LoadWallet, gen *LoadGenesis) *BootLoader {

	return &BootLoader{
		udp:           udp,
		packetFactory: packetFactory,
		db:            db,
		wallet:        wallet,
		genesis:       gen,
	}
}

func (b *BootLoader) BootLoading(user *ptypes.GhostUser, passwdHash []byte) *gcrypto.Wallet {
	// Load Wallet
	w, err := b.wallet.OpenWallet(user.Nickname, passwdHash)
	if w == nil || err != nil {
		w = b.wallet.CreateWallet(user.Nickname, passwdHash)
		b.wallet.SaveWallet(w, passwdHash)
	}

	// if Creator, need not to connect other master
	// Load Creator
	creators := b.genesis.CreatorList()
	if creator, exist := creators[user.Nickname]; exist {
		if creatorAddr, err := b.genesis.LoadCreatorKeyFile(creator.Nickname,
			creator.PubKey, passwdHash); err == nil {
			w = gcrypto.NewWallet(creator.Nickname, creatorAddr, &ptypes.GhostIp{Ip: user.Ip.Ip, Port: user.Ip.Port}, nil)
		} else {
			log.Println("Load Creator Key File Fail..")
			return nil
		}
		w.SetMasterNode(w.GetGhostUser())
		return w
	}

	// connect to master Node
	b.conn = NewConnectMaster(store.DefaultNodeTable, b.db, b.packetFactory, b.udp, w)
	masterNode := b.conn.LoadMasterNode()

	if masterNode == nil {
		b.conn.RequestMasterNodesToAdam()
		if timeout := b.conn.WaitEvent(); timeout {
			log.Fatal("could not connect to adam node")
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
