package maincontainer

import (
	"fmt"
	"log"
	"strconv"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/factory"
	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/bootloader"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gapi"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/grpc"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
)

type MainContainer struct {
	config         *gconfig.GConfig
	udp            *p2p.UdpServer
	networkFactory *factory.NetworkFactory
	bootFactory    *factory.BootFactory
	defaultFactory *factory.DefaultFactory
	grpcServer     *grpc.GrpcServer
	ghostApi       *gapi.GhostApi
}

func NewMainContainer(config *gconfig.GConfig) *MainContainer {
	networkFactory := factory.NewNetworkFactory(config)
	bootFactory := factory.NewBootFactory(networkFactory.Udp, networkFactory.PacketFactory, config)
	return &MainContainer{
		config:         config,
		networkFactory: networkFactory,
		bootFactory:    bootFactory,
	}
}

func (main *MainContainer) StartContainer() {
	if port, err := strconv.Atoi(main.config.Port); err != nil {
		log.Fatal(err)
	} else {
		main.config.GrpcPort = fmt.Sprint(port + 1)
	}

	main.grpcServer = grpc.NewGrpcServer()
	log.Println("Start Grpc Server")

	main.grpcServer.LoginContainerHandler = func(passwdHash []byte, username, ip, port string) bool {
		if w, _ := main.bootFactory.LoadWallet.OpenWallet(username, passwdHash); w == nil {
			log.Println("Login fail user = ", username)
			return false
		}
		main.config.Username = username
		main.config.Password = passwdHash
		main.StartBootLoading()
		return true
	}

	if err := main.grpcServer.ServeGRPC(main.config); err != nil {
		log.Fatal(err)
	}
}

func (main *MainContainer) StartBootLoading() {
	booter := bootloader.NewBootLoader(factory.BootTables, main.networkFactory.Udp,
		main.networkFactory.PacketFactory, main.config, main.bootFactory.Db,
		main.bootFactory.LoadWallet)
	w := booter.BootLoading(main.config)
	if w == nil {
		return
	}

	main.defaultFactory = factory.NewDefaultFactory(main.networkFactory, w, main.config)
	main.ghostApi = gapi.NewGhostApi(main.grpcServer, main.defaultFactory.Block, main.defaultFactory.BlockContainer,
		main.bootFactory.LoadWallet, main.config)
	main.StartServer()
}

func (main *MainContainer) StartServer() {
	main.udp.Start(nil)
}
