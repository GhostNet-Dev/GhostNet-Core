package commands

import (
	"fmt"
	"strings"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	ghostDeamonName = "ghostd"
	port            string
	host            string
	username        string
	password        string
	id              uint32
	rpcPort         string
)

// RootCmd root command binding
func NewRootCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "ghost",
		Short: "terminal in GhostNet",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			ExecuteNode()
		},
	}
	cmd.Flags().StringVarP(&host, "ip", "i", "", "Host Address")
	cmd.Flags().StringVarP(&port, "port", "", "50129", "Port Number")
	cmd.Flags().StringVarP(&username, "username", "u", "", "Ghost Account Nickname")
	cmd.Flags().StringVarP(&password, "password", "p", "", "Ghost Account Password")

	return cmd
}

func initializeConfig(cmd *cobra.Command) error {
	var defaultCfg = gconfig.DefaultConfig()
	v := viper.New()
	v.SetConfigName(defaultCfg.DefaultConfigFilename)
	v.AddConfigPath(".")
	v.AddConfigPath("../")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.SetEnvPrefix(defaultCfg.EnvPrefix)
	v.AutomaticEnv()

	bindFlags(cmd, v)
	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			var defaultCfg = gconfig.DefaultConfig()
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", defaultCfg.EnvPrefix, envVarSuffix))
		}

		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
