package workload

import (
	"errors"
	"math/rand"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/cmd/dummy/common"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blockmanager"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/bootloader"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/txs"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type Workload struct {
	workerName     string
	loadWallet     *bootloader.LoadWallet
	wallet         *gcrypto.Wallet
	blockContainer *store.BlockContainer
	tXs            *txs.TXs
	blockMgr       *blockmanager.BlockManager
	conn           *common.ConnectMaster
	glog           *glogger.GLogger
	Running        bool
}

var password = "worker"

func NewWorkload(workerName string, loadWallet *bootloader.LoadWallet, blkMgr *blockmanager.BlockManager,
	bc *store.BlockContainer, tXs *txs.TXs, conn *common.ConnectMaster, user *gcrypto.Wallet,
	glog *glogger.GLogger) IWorkload {
	return &Workload{
		workerName:     workerName,
		loadWallet:     loadWallet,
		blockContainer: bc,
		tXs:            tXs,
		blockMgr:       blkMgr,
		conn:           conn,
		wallet:         user,
		glog:           glog,
		Running:        false,
	}
}

func (worker *Workload) CheckRunning() bool {
	return worker.Running
}

func (worker *Workload) LoadWorker(masterNode *ptypes.GhostUser) {
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

func (w *Workload) PrepareRun() {
	w.Running = true
	for {
		if exist, err := w.CheckAccountTx(); !exist && err == nil {
			w.MakeAccountTx()
		} else if exist {
			w.MakeDataTx()
			break
		} else {
			w.glog.DebugOutput(w, err.Error(), glogger.Default)
			break
		}
		time.Sleep(time.Second * 3)
	}
}

func (w *Workload) Run() {
	//w.Running = false
	time.Sleep(time.Second * 3)
	w.MakeDataTx()
}

func (w *Workload) CheckAccountTx() (result bool, err error) {
	eventChannel := make(chan bool, 1)
	w.blockMgr.RequestCheckExistFsRoot([]byte(w.workerName), func(result bool) {
		eventChannel <- result
	})

	select {
	case result = <-eventChannel:
	case <-time.After(time.Second * 8):
		return false, errors.New("timeout")
	}

	return result, nil
}

func (w *Workload) MakeAccountTx() {
	txInfo := &txs.TransferTxInfo{
		MyWallet:  w.wallet,
		ToAddr:    w.wallet.MyPubKey(),
		Broker:    w.wallet.GetMasterNodeAddr(),
		FeeAddr:   store.AdamsAddress(),
		FeeBroker: w.wallet.GetMasterNodeAddr(),
	}
	tx := w.tXs.CreateRootFsTx(*txInfo, w.workerName)
	tx = w.tXs.InkTheContract(tx, w.wallet.GetGhostAddress())

	w.blockMgr.SendTx(tx, nil)
}

func (w *Workload) MakeDataTx() {
	outputParams, ok := w.tXs.GetRootFsTx(w.wallet.MyPubKey())
	if !ok {
		return
	}

	prevMap := map[types.TxOutputType][]types.PrevOutputParam{}
	prevMap[types.TxTypeDataStore] = outputParams // for mapping

	txInfo := &txs.TransferTxInfo{
		Prevs:     prevMap,
		MyWallet:  w.wallet,
		ToAddr:    w.wallet.MyPubKey(),
		Broker:    w.wallet.GetMasterNodeAddr(),
		FeeAddr:   store.AdamsAddress(),
		FeeBroker: w.wallet.GetMasterNodeAddr(),
	}
	tx, dataTx := w.tXs.CreateDataTx(*txInfo, []byte("adam"), w.MakeDummyFile())
	tx = w.tXs.InkTheContract(tx, w.wallet.GetGhostAddress())

	w.blockMgr.SendDataTx(tx, dataTx, nil)
}

func (w *Workload) MakeDummyTransaction(wallet *gcrypto.Wallet, to []byte, broker []byte) (*types.GhostTransaction, *types.GhostDataTransaction) {
	return nil, nil
}

func (w *Workload) MakeDummyFile() []byte {
	rand.Seed(time.Now().Unix())
	length := 32
	ranStr := make([]byte, length)
	for i := 0; i < length; i++ {
		ranStr[i] = byte(65 + rand.Intn(25))
	}
	return ranStr
}
