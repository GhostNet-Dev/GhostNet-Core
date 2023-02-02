package maincontainer

import (
	"github.com/GhostNet-Dev/GhostNet-Core/internal/factory"
	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/bootloader"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
)

type MainContainer struct {
	config *gconfig.GConfig
	udp    *p2p.UdpServer
}

func NewMainContainer(config *gconfig.GConfig) *MainContainer {
	return &MainContainer{config: config}
}

func (main *MainContainer) StartBootLoading() {
	netFactory := factory.NewNetworkFactory(main.config)
	bootloader.BootLoader(netFactory.Udp, netFactory.PacketFactory, main.config)
}

func (main *MainContainer) StartServer() {
	main.udp.Start(nil)
}
