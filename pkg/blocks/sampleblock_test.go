package blocks

import (
	"bytes"
	"io/ioutil"
	"log"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"github.com/btcsuite/btcutil/base58"
	"github.com/stretchr/testify/assert"
)

func TestMakeGenesis(t *testing.T) {
	creator := []string{
		"Alice", "Bob", "Carol", "Dave", "Eve", "Frank", "Grace",
		"Heidi", "Ivan", "Judy", "Michael", "Niaj", "Oscar", "Peggy",
		"Rupert", "Sybil", "Theo", "Terry", "Victor", "Walter", "Wendy",
	}
	genesis, accountFile := blocks.MakeGenesisBlock(creator)
	blockFilename := "1@" + base58.Encode(genesis.GetHashKey()) + ".ghost"
	if err := ioutil.WriteFile(blockFilename, genesis.SerializeToByte(), 0); err != nil {
		log.Fatal(err)
	}

	for filename, privateKey := range accountFile {
		if err := ioutil.WriteFile(filename, privateKey, 0); err != nil {
			log.Fatal(err)
		}
	}

	loadedGenesis, err := ioutil.ReadFile(blockFilename)
	if err != nil {
		log.Fatal(err)
	}

	newBlock := types.GhostNetBlock{}
	byteBuf := bytes.NewBuffer(loadedGenesis)
	newBlock.Deserialize(byteBuf)
	newBlockByte := newBlock.SerializeToByte()
	result := bytes.Compare(loadedGenesis, newBlockByte)

	assert.Equal(t, 0, result, "binary is different.")
}
