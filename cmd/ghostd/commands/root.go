package commands

import (
	"fmt"
	"log"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

var (
	port string
	host string
	id   uint32 = 0
)

// RootCmd root command binding
func RootCommand() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "ghostd",
		Short: "GhostNet Deamon",
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			log.Printf("%v", args)
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			log.Printf("%v", args)
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			ExecuteContainer()
		},
	}
	rootCmd.Flags().StringVarP(&port, "port", "p", "50129", "Port Number")
	rootCmd.Flags().StringVarP(&host, "ip", "i", "", "Host Address")
	return rootCmd
}

func ExecuteContainer() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	//netChannel := make(chan []byte)
	id++
	go func(cid uint32) {
		fmt.Printf("execute node addr = %s:%s\n", host, port)
		args := []string{"-p=" + port}
		if host != "" {
			args = append(args, host)
		}
		cmd := exec.Command("ghostnet", args...)
		out, err := cmd.StdoutPipe()
		cmd.Stderr = cmd.Stdout
		if err != nil {
			log.Fatal(err)
		}
		if err := cmd.Start(); err != nil {
			log.Fatal(err)
		}
		log.Println(cmd.Process.Pid)

		outBuf := make([]byte, 128)
		for {
			_, err := out.Read(outBuf)
			log.Printf("[%d] %s", cid, string(outBuf))
			if err != nil {
				log.Fatal(err)
			}
		}
	}(id)
	for {
		<-time.After((time.Second * 3))
	}
}
