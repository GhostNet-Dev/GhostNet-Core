package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/factory"
	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/internal/maincontainer"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
)

var (
	cfg      = gconfig.NewDefaultConfig()
	password string
)

// RootCmd root command binding
func RootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ghostnet",
		Short: "GhostNet in MasterNode",
		Long:  `GhostNet Core Server for Distributed BlockChain Network`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			cfg.Password = gcrypto.PasswordToSha256(password)

			fmt.Printf("Start GhostNet Node Addr = %s:%s", cfg.Ip, cfg.Port)
			glog := glogger.NewGLogger(cfg.Id, glogger.GetFullLogger())
			networkFactory := factory.NewNetworkFactory(cfg, glog)
			bootFactory := factory.NewBootFactory(networkFactory.Udp, networkFactory.PacketFactory, cfg, glog)
			container := maincontainer.NewMainContainer(networkFactory, bootFactory, cfg)
			container.StartContainer()
		},
	}
	cmd.Flags().StringVarP(&cfg.Username, "username", "u", "", "Ghost Account Nickname")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Ghost Account Password")
	cmd.Flags().StringVarP(&cfg.Ip, "ip", "i", gconfig.DefaultIp, "Host Address")
	cmd.Flags().StringVarP(&cfg.Port, "port", "", "50129", "Port Number")
	cmd.Flags().StringVarP(&cfg.RootPath, "rootpath", "", "", "Home Directory Path")
	cmd.Flags().StringVarP(&cfg.SqlPath, "sqlpath", "", "", "Sql Db File Directory Path")
	cmd.Flags().StringVarP(&cfg.FilePath, "filepath", "", "", "Download File Directory Path")
	cmd.Flags().BoolVarP(&cfg.StandaloneMode, "standalonemode", "", false, "Single Node Mode")
	cmd.Flags().Uint32VarP(&cfg.Id, "id", "", 0, "Container Id")

	cmd.MarkFlagRequired("username")
	return cmd
}
func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()
	v.SetConfigName(cfg.DefaultConfigFilename)
	v.AddConfigPath(".")
	v.AddConfigPath("../")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.SetEnvPrefix(cfg.EnvPrefix)
	v.AutomaticEnv()

	bindFlags(cmd, v)
	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", cfg.EnvPrefix, envVarSuffix))
		}

		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
