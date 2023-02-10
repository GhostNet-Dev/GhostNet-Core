package bootloader

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
)

var Tables = []string{"nodes", "wallet"}

func BootLoader(udp *p2p.UdpServer, packetFactory *p2p.PacketFactory, config *gconfig.GConfig) *gcrypto.Wallet {
	db := NewLiteStore(config.SqlPath, Tables)
	wallet := NewLoadWallet(Tables[1], db, &ptypes.GhostIp{Ip: config.Ip, Port: config.Port})
	w, err := wallet.OpenWallet(config.Username, config.Password)
	if err != nil {
		w = wallet.CreateWallet(config.Username, config.Password)
		wallet.SaveWallet(w, config.Password)
	}

	if config.StandaloneMode {
		return w
	}

	conn := NewConnectMaster(Tables[0], db, packetFactory, udp, w)
	masterNode := conn.LoadMasterNode()

	if masterNode == nil {
		conn.RequestMasterNodesToAdam()
		if timeout := conn.WaitEvent(); timeout {
			log.Fatal("could not connect to adam node")
			return nil
		}
		masterNode = conn.LoadMasterNode()
		if masterNode == nil {
			log.Fatal("could not found masterNode")
		}
	}

	conn.ConnectToMaster(masterNode)
	if timeout := conn.WaitEvent(); timeout {
		log.Fatal("could not connect to masterNode = ", masterNode.Nickname)
	}

	w.SetMasterNode(masterNode)
	return w
}
