package workload

import (
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
	glog *glogger.GLogger) *Workload {
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
	if exist, err := w.CheckAccountTx(); !exist && err == nil {
		w.MakeAccountTx()
	} else {
		w.glog.DebugOutput(w, err.Error(), glogger.Default)
	}
}

func (w *Workload) Run() {
	//w.Running = false
	time.Sleep(time.Second)
}

func (w *Workload) CheckAccountTx() (bool, error) {
	if exist, err := w.conn.CheckNickname(w.workerName); exist && err == nil {
		return true, nil
	} else if !exist && err == nil {
		return false, nil
	} else {
		return false, err
	}
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

	w.blockMgr.SendTx(tx)
}

func (w *Workload) MakeDummyTransaction(wallet *gcrypto.Wallet, to []byte, broker []byte) (*types.GhostTransaction, *types.GhostDataTransaction) {
	return nil, nil
}

func (w *Workload) MakeDummyFile() []byte {
	return []byte("test tx")
}
