package maincontainer

import (
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
)

type MainContainer struct {
	config *gconfig.GConfig
	udp    *p2p.UdpServer
}

func NewMainContainer(config *gconfig.GConfig) *MainContainer {
	return &MainContainer{config: config}
}

func (main *MainContainer) StartContainer() {
	for {
		<-time.After((time.Second * 3))
	}
}

func (main *MainContainer) StartBootLoading() {
	/*
		netFactory := factory.NewNetworkFactory(main.config)
		bootFactory := factory.NewBootFactory(netFactory.Udp, netFactory.PacketFactory, main.config)
		bootloader := bootloader.NewBootLoader(factory.BootTables, netFactory.Udp, netFactory.PacketFactory, main.config,
		bootFactory.Db, bootFactory.Wallet)
	*/
}

func (main *MainContainer) StartServer() {
	main.udp.Start(nil)
}
