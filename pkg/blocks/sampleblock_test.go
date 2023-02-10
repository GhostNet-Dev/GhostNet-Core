package blocks

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"github.com/btcsuite/btcutil/base58"
	"github.com/stretchr/testify/assert"
	mems "github.com/traherom/memstream"
)

func TestMakeGenesisHeader(t *testing.T) {
	//todo block sign을 만들어야함
	genesis := blocks.MakeGenesisBlock(nil)
	block := genesis.Block
	header := block.Header
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
	genesis := blocks.MakeGenesisBlock(nil)
	blockByte := genesis.SerializeToByte()

	newBlock := types.PairedBlock{}
	byteBuf := bytes.NewBuffer(blockByte)
	newBlock.Deserialize(byteBuf)
	newBlockByte := newBlock.SerializeToByte()

	result := bytes.Compare(blockByte, newBlockByte)

	assert.Equal(t, 0, result, "binary is different.")
}

func TestMakeGenesisFileIo(t *testing.T) {
	if err := os.RemoveAll("./samples"); err != nil {
		log.Fatal(err)
	}
	if err := os.Mkdir("./samples", os.ModePerm); err != nil {
		log.Fatal(err)
	}
	accountFile := map[string]*gcrypto.GhostAddress{}

	genesis := blocks.MakeGenesisBlock(func(name string, address *gcrypto.GhostAddress) {
		accountFile[name+"@"+address.GetPubAddress()+".ghost"] = address
	})

	blockFilename := "./samples/1@" + base58.Encode(genesis.Block.GetHashKey()) + ".ghost"
	genesisBuf := genesis.SerializeToByte()
	if err := ioutil.WriteFile(blockFilename, genesisBuf, 0); err != nil {
		log.Fatal(err)
	}

	for name, ghostAddr := range accountFile {
		filename := name + "@" + ghostAddr.GetPubAddress() + ".ghost"
		if err := ioutil.WriteFile("./samples/"+filename, ghostAddr.PrivateKeySerialize(), 0); err != nil {
			log.Fatal(err)
		}
	}

	loadedGenesis, err := ioutil.ReadFile(blockFilename)
	if err != nil {
		log.Fatal(err)
	}

	newBlock := types.PairedBlock{}
	byteBuf := bytes.NewBuffer(loadedGenesis)
	newBlock.Deserialize(byteBuf)
	newBlockByte := newBlock.SerializeToByte()

	result := bytes.Compare(loadedGenesis, genesisBuf)
	assert.Equal(t, 0, result, "original binary is different.")

	result = bytes.Compare(loadedGenesis, newBlockByte)
	assert.Equal(t, 0, result, "binary is different.")
}
