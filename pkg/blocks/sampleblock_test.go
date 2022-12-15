package blocks

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"github.com/btcsuite/btcutil/base58"
	"github.com/stretchr/testify/assert"
	mems "github.com/traherom/memstream"
)

var (
	creator = []string{
		"Alice", "Bob", "Carol", "Dave", "Eve", "Frank", "Grace",
		"Heidi", "Ivan", "Judy", "Michael", "Niaj", "Oscar", "Peggy",
		"Rupert", "Sybil", "Theo", "Terry", "Victor", "Walter", "Wendy",
	}
)

func TestMakeGenesisHeader(t *testing.T) {
	//todo block sign을 만들어야함
	genesis, _ := blocks.MakeGenesisBlock(creator)
	header := genesis.Header
	size := header.Size()
	stream := mems.NewCapacity(int(size))
	header.Serialize(stream)

	oriBuf := stream.Bytes()
	byteBuf := bytes.NewBuffer(oriBuf)
	newHeader := types.GhostNetBlockHeader{}
	newHeader.Deserialize(byteBuf)

	size = newHeader.Size()
	stream = mems.NewCapacity(int(size))
	newHeader.Serialize(stream)
	newByteBuf := stream.Bytes()
	result := bytes.Compare(oriBuf, newByteBuf)
	assert.Equal(t, 0, result, "binary is different.")
}

func TestMakeGenesis(t *testing.T) {
	genesis, _ := blocks.MakeGenesisBlock(creator)
	blockByte := genesis.SerializeToByte()

	newBlock := types.GhostNetBlock{}
	byteBuf := bytes.NewBuffer(blockByte)
	newBlock.Deserialize(byteBuf)
	newBlockByte := newBlock.SerializeToByte()

	result := bytes.Compare(blockByte, newBlockByte)

	assert.Equal(t, 0, result, "binary is different.")
}

func TestMakeGenesisFileIo(t *testing.T) {
	if err := os.Mkdir("./samples", os.ModePerm); err != nil {
		log.Fatal(err)
	}

	genesis, accountFile := blocks.MakeGenesisBlock(creator)
	blockFilename := "./samples/1@" + base58.Encode(genesis.GetHashKey()) + ".ghost"
	genesisBuf := genesis.SerializeToByte()
	if err := ioutil.WriteFile(blockFilename, genesisBuf, 0); err != nil {
		log.Fatal(err)
	}

	for filename, privateKey := range accountFile {
		if err := ioutil.WriteFile("./samples/"+filename, privateKey, 0); err != nil {
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

	result := bytes.Compare(loadedGenesis, genesisBuf)
	assert.Equal(t, 0, result, "original binary is different.")

	result = bytes.Compare(loadedGenesis, newBlockByte)
	assert.Equal(t, 0, result, "binary is different.")
}
