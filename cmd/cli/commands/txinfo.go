package commands

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"github.com/spf13/cobra"
)

var readFilename string

func ReadTxFileCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "readtx",
		Short: `read transaction file`,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			// You can bind cobra and viper in a few locations, but PersistencePreRunE on the root command works well
			return initializeConfig(cmd)
		},
		Run: func(cmd *cobra.Command, args []string) {
			// Working with OutOrStdout/OutOrStderr allows us to unit test our command easier
			data, err := ioutil.ReadFile(readFilename)
			if err != nil {
				log.Fatal(err)
			}
			tx := &types.GhostTransaction{}
			if !tx.Deserialize(bytes.NewBuffer(data)).Result() {
				return
			}
			if b, err := json.Marshal(tx); err == nil {
				log.Print(string(b))
			}
		},
	}
	cmd.Flags().StringVarP(&readFilename, "tx", "t", "", "decode tx file")
	return cmd
}
