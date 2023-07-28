package commands

import (
	"io/ioutil"
	"log"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/factory"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/bootloader"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/spf13/cobra"
)

var (
	executeScript = false
	codeFilepath  = "./sample.gs"
)

func ScriptCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "script",
		Short: `Register script in blockchain`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			if registerScriptCommand(username, password) {
				log.Println("Regist Complete")
			}
		},
	}
	cmd.Flags().StringVarP(&host, "ip", "i", "", "Host Address")
	cmd.Flags().StringVarP(&port, "port", "", "50129", "Port Number")
	cmd.Flags().StringVarP(&rpcPort, "rpc", "r", "50229", "GRPC Port Number")
	cmd.Flags().StringVarP(&username, "username", "u", "", "Ghost Account Nickname")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Ghost Account Password")
	cmd.Flags().StringVarP(&codeFilepath, "code", "c", "", "script file path")
	cmd.Flags().Uint32VarP(&id, "id", "", 0, "Container Id, if not select, show all container")
	cmd.Flags().Uint32VarP(&timeout, "timeout", "t", 3, "rpc connection timeout")
	cmd.Flags().BoolVarP(&executeScript, "exe", "e", false, "execute script")

	cmd.MarkFlagRequired("username")
	cmd.MarkFlagRequired("password")
	return cmd
}

func registerScriptCommand(username, password string) bool {
	cfg.Username = username
	// for encrypt passwd
	cfg.Password = gcrypto.PasswordToSha256(password)
	sampleCode, err := ioutil.ReadFile(codeFilepath)
	if err != nil {
		log.Fatal(err)
	}

	user := &ptypes.GhostUser{
		Nickname: cfg.Username,
		Ip:       &ptypes.GhostIp{Ip: cfg.Ip, Port: cfg.Port},
	}

	glog := glogger.NewGLogger(0)
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

	var txId []byte
	ok := false
	for {
		if txId == nil {
			if txId, ok = defaultFactory.ScriptIo.CreateScript(w, "workload", string(sampleCode)); !ok {
				time.Sleep(time.Second * 3)
				continue
			}
		}
		if scriptIoHandler := defaultFactory.ScriptIo.OpenScript(txId); scriptIoHandler == nil {
			time.Sleep(time.Second * 3)
			continue
		} else if executeScript {
			scriptIoHandler.ExecuteScript()
			break
		} else {
			break
		}
	}
	return true
}
