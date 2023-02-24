package blocks

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"github.com/btcsuite/btcutil/base58"
	"github.com/stretchr/testify/assert"
	mems "github.com/traherom/memstream"
	"google.golang.org/protobuf/proto"
)

var password = "test"

func TestMakeGenesisHeader(t *testing.T) {
	//todo block sign을 만들어야함
	accountFile := map[string]*gcrypto.GhostAddress{}
	genesis := blocks.MakeGenesisBlock(func(name string, address *gcrypto.GhostAddress) {
		accountFile[name+"@"+address.GetPubAddress()+".ghost"] = address
	})
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
	accountFile := map[string]*gcrypto.GhostAddress{}
	genesis := blocks.MakeGenesisBlock(func(name string, address *gcrypto.GhostAddress) {
		accountFile[name+"@"+address.GetPubAddress()+".ghost"] = address
	})
	blockByte := genesis.SerializeToByte()

	newBlock := types.PairedBlock{}
	byteBuf := bytes.NewBuffer(blockByte)
	newBlock.Deserialize(byteBuf)
	newBlockByte := newBlock.SerializeToByte()

	result := bytes.Compare(blockByte, newBlockByte)

	assert.Equal(t, 0, result, "binary is different.")
}

func TestMakeGenesisFileIo(t *testing.T) {

	accountFile := map[string]*gcrypto.GhostAddress{}

	genesis := blocks.MakeGenesisBlock(func(name string, address *gcrypto.GhostAddress) {
		accountFile[name+"@"+address.GetPubAddress()+".ghost"] = address
	})

	blockFilename := "1@" + base58.Encode(genesis.Block.GetHashKey()) + ".ghost"
	blockFilePath := path.Join("./samples/", blockFilename)
	genesisBuf := genesis.SerializeToByte()

	if err := os.RemoveAll("./samples"); err != nil {
		log.Fatal(err)
	}
	if err := os.Mkdir("./samples", os.ModePerm); err != nil {
		log.Fatal(err)
	}
	if err := ioutil.WriteFile(blockFilePath, genesisBuf, 0); err != nil {
		log.Fatal(err)
	}
	/*
		blockCopyPath := "../store/genesisblock"
		if err := os.Remove(blockCopyPath); err != nil {
			log.Fatal(err)
		}
		if err := ioutil.WriteFile(blockCopyPath, genesisBuf, 0); err != nil {
			log.Fatal(err)
		}
	*/

	for filename, ghostAddr := range accountFile {
		keyPair := &ptypes.KeyPair{
			PubKey:     ghostAddr.GetPubAddress(),
			PrivateKey: ghostAddr.PrivateKeySerialize(),
		}
		data, err := proto.Marshal(keyPair)
		if err != nil {
			log.Fatal(err)
		}
		cipherPivateKey := gcrypto.Encryption(gcrypto.PasswordToSha256(password), data)
		if err := ioutil.WriteFile("./samples/"+filename, cipherPivateKey, 0); err != nil {
			log.Fatal(err)
		}
	}

	loadedGenesis, err := ioutil.ReadFile(blockFilePath)
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
