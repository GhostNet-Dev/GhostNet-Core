package workload

import (
	"math/rand"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/cmd/dummy/common"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blockfilesystem"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/bootloader"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type WorkloadFs struct {
	workerName     string
	loadWallet     *bootloader.LoadWallet
	wallet         *gcrypto.Wallet
	conn           *common.ConnectMaster
	blockIo        *blockfilesystem.BlockIo
	blockIoHandler *blockfilesystem.BlockIoHandler
	glog           *glogger.GLogger
	Running        bool
}

func NewWorkloadFs(workerName string, loadWallet *bootloader.LoadWallet,
	blockIo *blockfilesystem.BlockIo, conn *common.ConnectMaster, user *gcrypto.Wallet,
	glog *glogger.GLogger) IWorkload {
	return &WorkloadFs{
		workerName: workerName,
		loadWallet: loadWallet,
		conn:       conn,
		blockIo:    blockIo,
		wallet:     user,
		glog:       glog,
		Running:    false,
	}
}

func (worker *WorkloadFs) CheckRunning() bool {
	return worker.Running
}

func (worker *WorkloadFs) LoadWorker(masterNode *ptypes.GhostUser) {
	username := worker.workerName
	cipherText := gcrypto.PasswordToSha256(password)

	w, err := worker.loadWallet.OpenWallet(username, cipherText)
	if err != nil {
		w = worker.loadWallet.CreateWallet(username, cipherText)
		worker.loadWallet.SaveWallet(w, cipherText)
	}
	worker.wallet = w
	worker.wallet.SetMasterNode(masterNode)
}

func (w *WorkloadFs) PrepareRun() {
	w.Running = true
	if exist, err := w.CheckAccountTx(); !exist && err == nil {
		if w.blockIoHandler = w.blockIo.OpenFilesystem(w.wallet); w.blockIoHandler == nil {
			w.blockIoHandler = w.blockIo.CreateFilesystem(w.wallet)
		}
	} else {
		w.blockIoHandler.WriteData([]byte("worker1"), w.MakeDummyFile())
	}
}

func (w *WorkloadFs) Run() {
	//w.Running = false
	time.Sleep(time.Second * 3)
	w.blockIoHandler.WriteData([]byte("worker1"), w.MakeDummyFile())
}

func (w *WorkloadFs) CheckAccountTx() (bool, error) {
	if exist, err := w.conn.CheckNickname(w.workerName); exist && err == nil {
		return true, nil
	} else if !exist && err == nil {
		return false, nil
	} else {
		return false, err
	}
}

func (w *WorkloadFs) MakeDummyTransaction(wallet *gcrypto.Wallet, to []byte, broker []byte) (*types.GhostTransaction, *types.GhostDataTransaction) {
	return nil, nil
}

func (w *WorkloadFs) MakeDummyFile() []byte {
	rand.Seed(time.Now().Unix())
	length := 32
	ranStr := make([]byte, length)
	for i := 0; i < length; i++ {
		ranStr[i] = byte(65 + rand.Intn(25))
	}
	return ranStr
}
