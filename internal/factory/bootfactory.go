package factory

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/bootloader"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
)

type BootFactory struct {
	Db         *store.LiteStore
	LoadWallet *bootloader.LoadWallet
	Genesis    *bootloader.LoadGenesis
	glog       *glogger.GLogger
}

func NewBootFactory(udp *p2p.UdpServer, packetFactory *p2p.PacketFactory, config *gconfig.GConfig, glog *glogger.GLogger) *BootFactory {
	db := store.NewLiteStore(config.SqlPath, config.LiteStoreFilename, store.DefaultLiteTable[:], 3)
	if err := db.OpenStore(); err != nil {
		log.Fatal(err)
	}

	wallet := bootloader.NewLoadWallet(store.DefaultWalletTable, db, &ptypes.GhostIp{Ip: config.Ip, Port: config.Port})
	genesis := bootloader.NewLoadGenesis(gvm.NewGVM(), "./")

	return &BootFactory{
		Db:         db,
		LoadWallet: wallet,
		Genesis:    genesis,
		glog:       glog,
	}
}
