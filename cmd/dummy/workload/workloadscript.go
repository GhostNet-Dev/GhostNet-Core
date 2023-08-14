package workload

import (
	"errors"
	"math/rand"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/common"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blockfilesystem"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blockmanager"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/bootloader"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type WorkloadScript struct {
	workerName      string
	loadWallet      *bootloader.LoadWallet
	wallet          *gcrypto.Wallet
	conn            *common.ConnectMaster
	blockMgr        *blockmanager.BlockManager
	scriptIo        *blockfilesystem.ScriptIo
	scriptIoHandler *blockfilesystem.ScriptIoHandler
	glog            *glogger.GLogger
	Running         bool
}

func NewWorkloadScript(workerName string, loadWallet *bootloader.LoadWallet, blockMgr *blockmanager.BlockManager,
	scriptIo *blockfilesystem.ScriptIo, conn *common.ConnectMaster, user *gcrypto.Wallet,
	glog *glogger.GLogger) IWorkload {
	return &WorkloadScript{
		workerName: workerName,
		loadWallet: loadWallet,
		conn:       conn,
		blockMgr:   blockMgr,
		scriptIo:   scriptIo,
		wallet:     user,
		glog:       glog,
		Running:    false,
	}
}

func (worker *WorkloadScript) CheckRunning() bool {
	return worker.Running
}

func (w *WorkloadScript) PrepareRun() {
	w.Running = true
	var txId []byte
	ok := false
	sampleCode := `
	let key = saveKeyValue("sampleKey", "testValue");
	let result = loadKeyValue(key);
	result;
	`

	for {
		if txId == nil {
			if txId, ok = w.scriptIo.CreateScript(ptypes.ScriptType_Default, w.wallet, "workload", sampleCode); !ok {
				time.Sleep(time.Second * 3)
				continue
			}
		}
		if w.scriptIoHandler = w.scriptIo.OpenScript(txId); w.scriptIoHandler == nil {
			time.Sleep(time.Second * 3)
			continue
		} else {
			w.scriptIoHandler.ExecuteScript()
			break
		}
	}
}

func (w *WorkloadScript) Run() {
	//w.Running = false
	time.Sleep(time.Second * 3)
	//w.scriptIoHandler.WriteScriptData([]byte(w.wallet.GetNickname()), w.MakeDummyFile())
	result := w.scriptIoHandler.ExecuteScript()
	w.glog.DebugOutput(w, result, glogger.Default)
}

func (w *WorkloadScript) CheckAccountTx() (result bool, err error) {
	eventChannel := make(chan bool, 1)
	w.blockMgr.RequestCheckExistFsRoot([]byte(w.wallet.GetNickname()), func(result bool) {
		eventChannel <- result
	})

	select {
	case result = <-eventChannel:
	case <-time.After(time.Second * 8):
		return false, errors.New("timeout")
	}

	return result, nil
}

func (w *WorkloadScript) MakeDummyTransaction(wallet *gcrypto.Wallet, to []byte, broker []byte) (*types.GhostTransaction, *types.GhostDataTransaction) {
	return nil, nil
}

func (w *WorkloadScript) MakeDummyFile() []byte {
	rand.Seed(time.Now().Unix())
	length := 32
	ranStr := make([]byte, length)
	for i := 0; i < length; i++ {
		ranStr[i] = byte(65 + rand.Intn(25))
	}
	return ranStr
}
