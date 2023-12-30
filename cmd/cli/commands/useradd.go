package commands

import (
	"log"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/factory"
	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/bootloader"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/grpc"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/spf13/cobra"
)

var (
	cfg = gconfig.NewDefaultConfig()
)

func UserAddCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "useradd",
		Short: `stop a GhostNet Server`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			if createUserAccountCommand(username, password) {
				log.Printf("Create User: %s", username)
				userAddExecuteCommand(username, password)
			}
		},
	}
	cmd.Flags().StringVarP(&host, "ip", "i", "", "Host Address")
	cmd.Flags().StringVarP(&port, "port", "", "50129", "Port Number")
	cmd.Flags().StringVarP(&rpcPort, "rpc", "r", "50229", "GRPC Port Number")
	cmd.Flags().StringVarP(&username, "username", "u", "", "Ghost Account Nickname")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Ghost Account Password")
	cmd.Flags().Uint32VarP(&id, "id", "", 0, "Container Id, if not select, show all container")
	cmd.Flags().Uint32VarP(&timeout, "timeout", "t", 3, "rpc connection timeout")

	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")
	return cmd
}

func userAddExecuteCommand(username, password string) {
	grpcClient := grpc.NewGrpcClient(host, rpcPort, timeout)
	grpcClient.ConnectServer()
	defer grpcClient.CloseServer()

	ret := grpcClient.CreateAccount(username, password)
	log.Printf("Create User: %s= %t", username, ret)
}

func createUserAccountCommand(username, password string) bool {
	cfg.Username = username
	// for encrypt passwd
	cfg.Password = gcrypto.PasswordToSha256(password)
	cfg.Ip = host
	cfg.Port = port
	cfg.GrpcPort = rpcPort

	user := &ptypes.GhostUser{
		Nickname: cfg.Username,
		Ip:       &ptypes.GhostIp{Ip: cfg.Ip, Port: cfg.Port},
	}

	glog := glogger.NewGLogger(0, glogger.GetFullLogger())
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

	if blockIoHandler := defaultFactory.BlockIo.OpenFilesystem(w); blockIoHandler == nil {
		defaultFactory.BlockIo.CreateFilesystem(w)
	} else {
		log.Println("Already Exist Username = ", username)
		return false
	}

	for {
		if blockIoHandler := defaultFactory.BlockIo.OpenFilesystem(w); blockIoHandler != nil {
			break
		}
		<-time.After(time.Second * time.Duration(timeout))
	}
	return true
}
