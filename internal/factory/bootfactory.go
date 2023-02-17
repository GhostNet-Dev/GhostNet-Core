package factory

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/bootloader"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
)

type BootFactory struct {
	Db         *bootloader.LiteStore
	LoadWallet *bootloader.LoadWallet
	Genesis    *bootloader.LoadGenesis
}

var BootTables = []string{"nodes", "wallet"}

func NewBootFactory(udp *p2p.UdpServer, packetFactory *p2p.PacketFactory, config *gconfig.GConfig) *BootFactory {
	db := bootloader.NewLiteStore(config.SqlPath, BootTables)
	if err := db.OpenStore(); err != nil {
		log.Fatal(err)
	}

	wallet := bootloader.NewLoadWallet(BootTables[1], db, &ptypes.GhostIp{Ip: config.Ip, Port: config.Port})
	genesis := bootloader.NewLoadGenesis(gvm.NewGVM(), "./")

	return &BootFactory{
		Db:         db,
		LoadWallet: wallet,
		Genesis:    genesis,
	}
}
