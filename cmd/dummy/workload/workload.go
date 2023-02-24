package workload

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/bootloader"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type Workload struct {
	workerName     string
	loadWallet     *bootloader.LoadWallet
	wallet         *gcrypto.Wallet
	blockContainer *store.BlockContainer
}

var password = "worker"

func NewWorkload(workerName string, loadWallet *bootloader.LoadWallet, bc *store.BlockContainer) *Workload {
	return &Workload{
		workerName:     workerName,
		loadWallet:     loadWallet,
		blockContainer: bc,
	}
}

func (worker *Workload) LoadWorker() {
	username := worker.workerName
	cipherText := gcrypto.PasswordToSha256(password)

	w, err := worker.loadWallet.OpenWallet(username, cipherText)
	if err != nil {
		w = worker.loadWallet.CreateWallet(username, cipherText)
		worker.loadWallet.SaveWallet(w, cipherText)
	}
	worker.wallet = w
}

func (w *Workload) PrepareRun() {
}

func (w *Workload) Run() {

}

func (w *Workload) MakeDummyTransaction() *types.GhostTransaction {
	return nil
}

func (w *Workload) MakeDummyFile() []byte {
	return []byte("test tx")
}
