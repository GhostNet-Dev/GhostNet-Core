package workload

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/cmd/dummy/common"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/bootloader"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
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
	conn           *common.ConnectMaster
	Running        bool
}

var password = "worker"

func NewWorkload(workerName string, loadWallet *bootloader.LoadWallet,
	bc *store.BlockContainer, tXs *txs.TXs, conn *common.ConnectMaster) *Workload {
	return &Workload{
		workerName:     workerName,
		loadWallet:     loadWallet,
		blockContainer: bc,
		tXs:            tXs,
		conn:           conn,
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
	if exist := w.CheckAccountTx(); !exist {
		w.MakeAccountTx()
	}
}

func (w *Workload) Run() {
	w.Running = false
}

func (w *Workload) CheckAccountTx() bool {
	if exist, _ := w.conn.CheckNickname(w.workerName); exist {
		return true
	}
	return false
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

	if _, err := w.conn.SendTx(tx); err != nil {
		log.Fatal(err)
	}
}

func (w *Workload) MakeDummyTransaction(wallet *gcrypto.Wallet, to []byte, broker []byte) (*types.GhostTransaction, *types.GhostDataTransaction) {
	return nil, nil
}

func (w *Workload) MakeDummyFile() []byte {
	return []byte("test tx")
}
