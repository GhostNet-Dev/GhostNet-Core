package commands

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/cmd/ghostd/manager"
	"github.com/GhostNet-Dev/GhostNet-Core/internal/factory"
	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/grpc"
	"github.com/spf13/cobra"
)

var (
	cfg      = gconfig.NewDefaultConfig()
	password string
)

// RootCmd root command binding
func RootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ghostd",
		Short: "GhostNet Deamon",
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
	return cmd
}

func ExecuteContainer() {
	log.Println("Initialize Component")
	glog := glogger.NewGLogger(0, glogger.GetFullLogger())
	// for encrypt passwd
	cfg.Password = gcrypto.PasswordToSha256(password)

	// network factory initialize
	netFactory := factory.NewNetworkFactory(cfg, glog)

	// boot factory initialize
	bootFactory := factory.NewBootFactory(netFactory.Udp, netFactory.PacketFactory, cfg, glog)

	// container management initialize
	containers := manager.NewContainers(netFactory, bootFactory, cfg)

	// grpc initiaize
	server := grpc.NewGrpcServer()
	manager.NewGrpcDeamonHandler(bootFactory.LoadWallet, bootFactory.Genesis, containers, server, cfg)

	log.Println("Start Grpc Server")
	if err := server.ServeGRPC(cfg.GrpcPort); err != nil {
		log.Fatal(err)
	}
}
