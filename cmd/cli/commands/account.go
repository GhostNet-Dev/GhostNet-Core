package commands

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/grpc"
	"github.com/spf13/cobra"
)

func GetAccountListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "userlist",
		Short: `Ghost Account List`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			grpcClient := grpc.NewGrpcClient(host, rpcPort, timeout)
			grpcClient.ConnectServer()
			defer grpcClient.CloseServer()
			users := grpcClient.GetAccount(id)
			if users != nil {
				for _, user := range users {
					log.Printf("[%d] %s\n", id, user.Nickname)
				}
			}
		},
	}
	cmd.Flags().StringVarP(&host, "ip", "i", "", "Host Address")
	cmd.Flags().StringVarP(&port, "port", "", "50129", "Port Number")
	cmd.Flags().StringVarP(&rpcPort, "rpc", "r", "50229", "GRPC Port Number")
	cmd.Flags().Uint32VarP(&id, "id", "", 0, "Container Id, if not select, show all container")
	cmd.Flags().Uint32VarP(&timeout, "timeout", "t", 3, "rpc connection timeout")
	return cmd
}

func LoginCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "login",
		Short: `login GhostNet Server`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			grpcClient := grpc.NewGrpcClient(host, rpcPort, timeout)
			grpcClient.ConnectServer()
			defer grpcClient.CloseServer()
			ret := grpcClient.LoginContainer(id, gcrypto.PasswordToSha256(password), username, host, port)
			fmt.Print("login = ", ret)
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
func GetPrivateKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "getkey",
		Short: `suspend a GhostNet Server`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			grpcClient := grpc.NewGrpcClient(host, rpcPort, timeout)
			grpcClient.ConnectServer()
			defer grpcClient.CloseServer()
			key, ret := grpcClient.GetPrivateKey(id, username, password)
			if ret {
				err := ioutil.WriteFile(fmt.Sprint("PrivateKey@", username), key, os.FileMode(0644))
				if err != nil {
					log.Fatal(err)
				}
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
	return cmd
}
