package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"

	cfg "github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
)

const (
	defaultConfigFilename = "global"
	envPrefix             = "GON" // GhOstNet
)

// StartCmdTest test command binding
var StartCmdTest = &cobra.Command{
	Use:   "test",
	Short: "test library",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("hi~ i'm Ghost")
	},
}

// RootCmd root command binding
var RootCmd = &cobra.Command{
	Use:   "GhostNet",
	Short: "GhostNet in MasterNode",
	Long:  `GhostNet Core Server for Distributed BlockChain Network`,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
		return initializeConfig(cmd)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
		var config = cfg.DefaultConfig()
		config.RootPath = rootPath
		config.Ip = host
		config.Port = port
	},
}

var (
	host     string
	port     string
	rootPath string
)

func init() {
	RootCmd.Flags().StringVarP(&host, "ip", "i", cfg.DefaultIp, "Host Address")
	RootCmd.Flags().StringVarP(&port, "port", "p", "50129", "Port Number")
	RootCmd.Flags().StringVarP(&rootPath, "rootPath", "", "", "Home Directory Path")
}

func initializeConfig(cmd *cobra.Command) error {
	v := viper.New()
	v.SetConfigName(defaultConfigFilename)
	v.AddConfigPath(".")
	v.AddConfigPath("../")

	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return err
		}
	}

	v.SetEnvPrefix(envPrefix)
	v.AutomaticEnv()

	bindFlags(cmd, v)
	return nil
}

func bindFlags(cmd *cobra.Command, v *viper.Viper) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			v.BindEnv(f.Name, fmt.Sprintf("%s_%s", envPrefix, envVarSuffix))
		}

		if !f.Changed && v.IsSet(f.Name) {
			val := v.Get(f.Name)
			cmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}
