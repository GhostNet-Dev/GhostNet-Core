package commands

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/cmd/dummy/dummyfactory"
	"github.com/GhostNet-Dev/GhostNet-Core/internal/factory"
	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/bootloader"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/spf13/cobra"
)

var (
	cfg        = gconfig.NewDefaultConfig()
	password   string
	masterIp   string
	masterPort string
)

// RootCmd root command binding
func RootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "test",
		Short: "GhostNet Dummy Test",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.SetFlags(log.LstdFlags | log.Lshortfile)
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			ExecuteContainer()
		},
	}
	cmd.Flags().StringVarP(&cfg.Username, "username", "u", "", "Ghost Account Nickname")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Ghost Account Password")
	cmd.Flags().StringVarP(&cfg.Ip, "ip", "i", gconfig.DefaultIp, "Host Address")
	cmd.Flags().StringVarP(&cfg.Port, "port", "", "50129", "Port Number")
	cmd.Flags().StringVarP(&cfg.GrpcPort, "rpc", "r", "50229", "GRPC Port Number")
	cmd.Flags().StringVarP(&cfg.RootPath, "rootpath", "", "", "Home Directory Path")
	cmd.Flags().StringVarP(&cfg.SqlPath, "sqlpath", "", "", "Sql Db File Directory Path")
	cmd.Flags().Uint32VarP(&cfg.Timeout, "timeout", "t", 3, "rpc connection timeout")

	cmd.Flags().StringVarP(&masterIp, "mip", "", "", "Master Node Ip")
	cmd.Flags().StringVarP(&masterPort, "mport", "", "", "Master Node Port")
	return cmd
}

func ExecuteContainer() {
	log.Println("Initialize Component")
	masterAddr := &ptypes.GhostIp{
		Ip:   masterIp,
		Port: masterPort,
	}
	user := &ptypes.GhostUser{
		Nickname: cfg.Username,
		Ip:       &ptypes.GhostIp{Ip: cfg.Ip, Port: cfg.Port},
	}

	glog := glogger.NewGLogger(0)
	// for encrypt passwd
	cfg.Password = gcrypto.PasswordToSha256(password)

	// network factory initialize
	netFactory := factory.NewNetworkFactory(cfg, glog)

	log.Println("Network Listen Start")
	netFactory.Udp.Start(nil, cfg.Ip, cfg.Port)

	// boot factory initialize
	bootFactory := factory.NewBootFactory(netFactory.Udp, netFactory.PacketFactory, cfg, glog)
	booter := bootloader.NewBootLoader(netFactory.Udp,
		netFactory.PacketFactory, bootFactory.Db,
		bootFactory.LoadWallet, bootFactory.Genesis)

	w := booter.BootLoading(user, cfg.Password)
	log.Print("User: ", w.GetGhostUser())
	log.Print("Master: ", w.GetMasterNode())

	defaultFactory := factory.NewDefaultFactory(netFactory, bootFactory, w, cfg, glog)
	defaultFactory.FactoryOpen()
	go func() { defaultFactory.BlockServer.BlockServer() }()

	dummyFactory := dummyfactory.NewDummyFactory(1, masterAddr, bootFactory, netFactory, defaultFactory, glog)
	log.Println("Worker Prepare Run")
	for _, worker := range dummyFactory.Worker {
		worker.PrepareRun()
	}

	log.Println("Worker Run")
	for {
		checkRunning := 0
		for _, worker := range dummyFactory.Worker {
			if worker.CheckRunning() {
				checkRunning++
			}
			worker.Run()
		}
		if checkRunning == 0 {
			break
		}
	}
}
