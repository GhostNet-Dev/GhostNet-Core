package maincontainer

import (
	"fmt"
	"log"
	"strconv"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/factory"
	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/bootloader"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gapi"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/grpc"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/p2p"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
)

type MainContainer struct {
	config         *gconfig.GConfig
	udp            *p2p.UdpServer
	networkFactory *factory.NetworkFactory
	bootFactory    *factory.BootFactory
	defaultFactory *factory.DefaultFactory
	grpcServer     *grpc.GrpcServer
	ghostApi       *gapi.GhostContainerApi
	glog           *glogger.GLogger

	user *ptypes.GhostUser
}

func NewMainContainer(networkFactory *factory.NetworkFactory, bootFactory *factory.BootFactory,
	config *gconfig.GConfig) *MainContainer {
	glog := glogger.NewGLogger(config.Id)

	return &MainContainer{
		config:         config,
		bootFactory:    bootFactory,
		networkFactory: networkFactory,
		udp:            networkFactory.Udp,
		glog:           glog,
	}
}

func (main *MainContainer) StartContainer() {
	if port, err := strconv.Atoi(main.config.Port); err != nil {
		log.Fatal(err)
	} else {
		main.config.GrpcPort = fmt.Sprint(port + 1)
	}

	main.grpcServer = grpc.NewGrpcServer()
	main.glog.DebugOutput(main, "Start Grpc Server", glogger.Default)

	main.grpcServer.LoginContainerHandler = func(id uint32, passwdHash []byte, username, ip, port string) bool {
		creators := main.bootFactory.Genesis.CreatorList()
		if _, exist := creators[username]; !exist {
			if w, err := main.bootFactory.LoadWallet.OpenWallet(username, passwdHash); w == nil {
				main.glog.DebugOutput(main, fmt.Sprint("Login fail user = ", username, ", err = ", err), glogger.Default)
				return false
			}
		}
		main.config.Username = username
		main.config.Password = passwdHash
		main.config.Ip = ip
		main.config.Port = port

		main.user = &ptypes.GhostUser{
			Nickname: username,
			Ip: &ptypes.GhostIp{
				Ip:   ip,
				Port: port,
			},
		}

		main.glog.DebugOutput(main, "Net Open", glogger.Default)
		main.udp.Start(nil, ip, port)

		main.glog.DebugOutput(main, "Start Bootloading", glogger.Default)
		main.startBootLoading()
		return true
	}

	if err := main.grpcServer.ServeGRPC(main.config.GrpcPort); err != nil {
		log.Fatal(err)
	}
}

func (main *MainContainer) startBootLoading() {
	booter := bootloader.NewBootLoader(main.networkFactory.Udp,
		main.networkFactory.PacketFactory, main.bootFactory.Db,
		main.bootFactory.LoadWallet, main.bootFactory.Genesis)
	w := booter.BootLoading(main.user, main.config.Password)
	if w == nil {
		return
	}

	main.defaultFactory = factory.NewDefaultFactory(main.networkFactory, main.bootFactory, w, main.config, main.glog)
	main.defaultFactory.FactoryOpen()
	main.ghostApi = gapi.NewGhostContainerApi(main.grpcServer, main.defaultFactory.Block, main.defaultFactory.BlockContainer,
		main.bootFactory.LoadWallet)
	main.glog.DebugOutput(main, "Initialize complete", glogger.Default)
	main.glog.DebugOutput(main, "Start Mainserver", glogger.Default)
	go main.StartServer()
}

func (main *MainContainer) StartServer() {
	main.defaultFactory.BlockServer.BlockServer()
}
