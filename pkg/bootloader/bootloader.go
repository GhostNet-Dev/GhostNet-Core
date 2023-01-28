package bootloader

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
)

var tables = []string{"nodes", "wallet"}

func BootLoader(udp *p2p.UdpServer, packetFactory *p2p.PacketFactory, config *gconfig.GConfig, nickname string) *gcrypto.Wallet {
	db := NewLiteStore(config.DbPath, tables)
	wallet := NewLoadWallet(tables[1], db, &ptypes.GhostIp{Ip: config.Ip, Port: config.Port})
	w, err := wallet.OpenWallet(nickname)
	if err != nil {
		w = wallet.CreateWallet(nickname)
		wallet.SaveWallet(w)
	}

	if config.StandaloneMode {
		return w
	}

	conn := NewConnectMaster(tables[0], db, packetFactory, udp, w)
	masterNode := conn.LoadMasterNode()

	if masterNode == nil {
		conn.RequestMasterNodesToAdam()
		conn.WaitEvent()
		masterNode = conn.LoadMasterNode()
		if masterNode == nil {
			log.Fatal("could not found masterNode")
		}
	}

	conn.ConnectToMaster(masterNode)
	conn.WaitEvent()

	w.SetMasterNode(masterNode)
	return w
}
