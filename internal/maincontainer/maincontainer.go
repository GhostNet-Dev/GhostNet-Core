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

func NewMainContainer(networkFactory *factory.NetworkFactory, bootFactory *factory.BootFactory,
	config *gconfig.GConfig) *MainContainer {
	if networkFactory == nil {
		networkFactory = factory.NewNetworkFactory(config)
	}
	if bootFactory == nil {
		bootFactory = factory.NewBootFactory(networkFactory.Udp, networkFactory.PacketFactory, config)
	}
	return &MainContainer{
		config:         config,
		bootFactory:    bootFactory,
		networkFactory: networkFactory,
		udp:            networkFactory.Udp,
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
		creators := main.bootFactory.Genesis.CreatorList()
		if _, exist := creators[username]; !exist {
			if w, _ := main.bootFactory.LoadWallet.OpenWallet(username, passwdHash); w == nil {
				log.Println("Login fail user = ", username)
				return false
			}
		}
		main.config.Username = username
		main.config.Password = passwdHash
		main.config.Ip = ip
		main.config.Port = port

		log.Println("Net Open")
		main.udp.Start(nil, ip, port)

		log.Println("Start Bootloading")
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
		main.bootFactory.LoadWallet, main.bootFactory.Genesis)
	w := booter.BootLoading(main.config)
	if w == nil {
		return
	}

	main.defaultFactory = factory.NewDefaultFactory(main.networkFactory, w, main.config)
	main.defaultFactory.FactoryOpen()
	main.ghostApi = gapi.NewGhostApi(main.grpcServer, main.defaultFactory.Block, main.defaultFactory.BlockContainer,
		main.bootFactory.LoadWallet, main.config)
	go main.StartServer()
}

func (main *MainContainer) StartServer() {
	main.defaultFactory.BlockServer.BlockServer()
}
