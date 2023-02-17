package bootloader

import (
	"errors"
	"io/ioutil"
	"log"
	"path"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/txs"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"google.golang.org/protobuf/proto"
)

type LoadGenesis struct {
	pairedBlock *types.PairedBlock
	gVm         *gvm.GVM
	keyPath     string
	creator     map[string]*ptypes.GhostUser
}

func NewLoadGenesis(gVm *gvm.GVM, keyPath string) *LoadGenesis {
	return &LoadGenesis{
		pairedBlock: store.GenesisBlock(),
		gVm:         gVm,
		keyPath:     keyPath,
		creator:     make(map[string]*ptypes.GhostUser),
	}
}

func (load *LoadGenesis) CreatorList() map[string]*ptypes.GhostUser {
	for _, tx := range load.pairedBlock.Block.Transaction {
		copyTx := tx.TxCopy()
		load.gVm.Clear()
		gFuncParam, toAddr := txs.ExtractFSRootGParam(&copyTx)
		if verify := load.gVm.ExecuteGFunction(copyTx.SerializeToByte(), gFuncParam); !verify {
			log.Fatal("wrong script")
		}
		nickname := string(load.gVm.Regs.Stack.Pop().([]byte))
		load.creator[nickname] = &ptypes.GhostUser{
			Nickname: nickname, PubKey: gcrypto.Translate160ToBase58Addr(toAddr)}
	}
	return load.creator
}

func (load *LoadGenesis) LoadCreatorKeyFile(nickname, pubKey string, password []byte) (*gcrypto.GhostAddress, error) {
	filename := nickname + "@" + pubKey + ".ghost"
	data, err := ioutil.ReadFile(path.Join(load.keyPath, filename))
	if err != nil {
		log.Print(err)
		return nil, err
	}
	der := gcrypto.Decryption(password, data)
	if der == nil {
		return nil, errors.New("password is wrong")
	}
	keyPair := &ptypes.KeyPair{}
	if err := proto.Unmarshal(der, keyPair); err != nil {
		log.Fatal(err)
	}

	ghostAddr := &gcrypto.GhostAddress{}
	ghostAddr.PrivateKeyDeserialize(keyPair.PrivateKey)
	if ghostAddr.GetPubAddress() != keyPair.PubKey {
		return nil, errors.New("password is wrong")
	}
	return ghostAddr, nil
}
