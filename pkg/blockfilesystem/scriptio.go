package blockfilesystem

import (
	"bytes"
	"log"
	"sync"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blockmanager"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/cloudservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/fileservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gvm"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/txs"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
	"google.golang.org/protobuf/proto"
)

type ScriptIo struct {
	bc              *store.BlockContainer
	cloud           *cloudservice.CloudService
	wallet          *gcrypto.Wallet
	blockManager    *blockmanager.BlockManager
	liteStore       *store.LiteStore
	tXs             *txs.TXs
	scriptIoHandler *ScriptIoHandler
}

type ScriptIoHandler struct {
	wallet        *gcrypto.Wallet
	scriptIo      *ScriptIo
	toAddr        []byte
	brokerAddr    []byte
	feeBrokerAddr []byte
	scriptTxPtr   []types.PrevOutputParam
	code          *ptypes.ScriptHeader
}

/*
todo list
1. make litestore uniqkey -> datatxid map
2. when called open, read script tx
3. read key ->  search key map with datatxid
*/

func NewScriptIo(blkMgr *blockmanager.BlockManager,
	bc *store.BlockContainer, tXs *txs.TXs, cloud *cloudservice.CloudService,
	w *gcrypto.Wallet, liteStore *store.LiteStore) *ScriptIo {
	gScriptIo := &ScriptIo{
		wallet:       w,
		liteStore:    liteStore,
		cloud:        cloud,
		bc:           bc,
		blockManager: blkMgr,
		tXs:          tXs,
	}
	return gScriptIo
}

func (io *ScriptIo) CreateScript(scriptType ptypes.ScriptType, w *gcrypto.Wallet, namespace, script string) ([]byte, bool) {
	outputParams, ok := io.tXs.GetRootFsTx(w.MyPubKey())
	if !ok {
		return nil, false
	}
	prevMap := map[types.TxOutputType][]types.PrevOutputParam{}
	prevMap[types.TxTypeScript] = outputParams // for mapping
	txInfo := &txs.TransferTxInfo{
		Prevs:     prevMap,
		FromAddr:  w.MyPubKey(),
		ToAddr:    w.MyPubKey(),
		Broker:    w.GetMasterNodeAddr(),
		FeeAddr:   store.AdamsAddress(),
		FeeBroker: w.GetMasterNodeAddr(),
	}
	scriptHeader := &ptypes.ScriptHeader{Type: scriptType, Version: 0, Script: script}
	scriptData, err := proto.Marshal(scriptHeader)
	if err != nil {
		log.Fatal(err)
	}
	tx, dataTx := io.tXs.CreateScriptTx(*txInfo, []byte(namespace), scriptData)
	tx = io.tXs.InkTheContract(tx, w.GetGhostAddress())

	wg := &sync.WaitGroup{}
	result := false
	io.blockManager.SendDataTx(tx, dataTx, func(b bool) {
		defer wg.Done()
		result = b
	})
	wg.Wait()
	return tx.TxId, result
}

func (io *ScriptIo) OpenScript(txId []byte) *ScriptIoHandler {
	tx, _ := io.bc.TxContainer.GetTx(txId)
	if tx == nil {
		return nil
	}
	output := tx.Body.Vout[0]
	if output.Type != types.TxTypeScript {
		return nil
	}
	dataTxId := output.ScriptEx
	dataTx := io.loadDataTx(dataTxId)
	if dataTx == nil {
		return nil
	}

	scriptHeader := &ptypes.ScriptHeader{}
	if err := proto.Unmarshal(dataTx.Data, scriptHeader); err != nil {
		log.Fatal(err)
	}
	io.scriptIoHandler = &ScriptIoHandler{
		scriptIo:      io,
		wallet:        io.wallet,
		toAddr:        output.Addr,
		brokerAddr:    io.wallet.MyPubKey(),
		feeBrokerAddr: output.Addr,
		scriptTxPtr: []types.PrevOutputParam{{
			TxType: types.TxTypeScript,
			VOutPoint: types.TxOutPoint{
				TxId:       txId,
				TxOutIndex: 0,
			},
			Vout: output,
		}},
		code: scriptHeader,
	}
	return io.scriptIoHandler
}

func (io *ScriptIo) CloseScript(handler *ScriptIoHandler) {}

func (io *ScriptIoHandler) ExecuteScript() string {
	return gvm.ExecuteScript(io.code.Script)
}

func (io *ScriptIoHandler) ReadScriptData(key []byte) (data []byte) {
	// to avoid key collision
	dataTxId := key
	dataTx := io.scriptIo.loadDataTx(dataTxId)
	if dataTx == nil {
		return nil
	}
	return dataTx.Data
}

func (io *ScriptIoHandler) WriteScriptData(uniqKey, data []byte) (key []byte,
	tx *types.GhostTransaction, dataTx *types.GhostDataTransaction) {
	prevs := io.scriptTxPtr // for mapping

	if ref := io.scriptIo.loadRefTx(uniqKey, io.toAddr); len(ref) != 0 {
		prevs = append(prevs, ref...)
	}
	prevMap := map[types.TxOutputType][]types.PrevOutputParam{}
	prevMap[types.TxTypeScriptStore] = prevs

	txInfo := &txs.TransferTxInfo{
		Prevs:     prevMap,
		FromAddr:  io.wallet.MyPubKey(),
		ToAddr:    io.toAddr,
		Broker:    io.brokerAddr,
		FeeAddr:   store.AdamsAddress(),
		FeeBroker: io.feeBrokerAddr,
	}
	tx, dataTx = io.scriptIo.tXs.CreateScriptDataTx(*txInfo, uniqKey, data)
	tx = io.scriptIo.tXs.InkTheContract(tx, io.wallet.GetGhostAddress())

	io.scriptIo.blockManager.SendDataTx(tx, dataTx, nil)

	key = dataTx.TxId
	return key, tx, dataTx
}

func (io *ScriptIoHandler) MakeScriptData(uniqKey, data []byte) (txInfo *txs.TransferTxInfo) {
	prevs := io.scriptTxPtr // for mapping

	if ref := io.scriptIo.loadRefTx(uniqKey, io.toAddr); len(ref) != 0 {
		prevs = append(prevs, ref...)
	}
	prevMap := map[types.TxOutputType][]types.PrevOutputParam{}
	prevMap[types.TxTypeScriptStore] = prevs

	txInfo = &txs.TransferTxInfo{
		Prevs:     prevMap,
		FromAddr:  io.wallet.MyPubKey(),
		ToAddr:    io.toAddr,
		Broker:    io.brokerAddr,
		FeeAddr:   store.AdamsAddress(),
		FeeBroker: io.feeBrokerAddr,
	}
	return txInfo
}

func (io *ScriptIoHandler) CommitScriptData(uniqKey, data [][]byte, txInfos []*txs.TransferTxInfo) (
	tx *types.GhostTransaction, dataTxs []*types.GhostDataTransaction) {

	prevMap := map[types.TxOutputType][]types.PrevOutputParam{}
	for _, info := range txInfos {
		prevMap[types.TxTypeScriptStore] = append(prevMap[types.TxTypeScriptStore], info.Prevs[types.TxTypeScriptStore]...)
	}
	txInfo := &txs.TransferTxInfo{
		Prevs:     prevMap,
		FromAddr:  io.wallet.MyPubKey(),
		ToAddr:    io.toAddr,
		Broker:    io.brokerAddr,
		FeeAddr:   store.AdamsAddress(),
		FeeBroker: io.feeBrokerAddr,
	}
	tx, dataTxs = io.scriptIo.tXs.CreateScriptMultiDataTx(*txInfo, uniqKey, data)
	tx = io.scriptIo.tXs.InkTheContract(tx, io.wallet.GetGhostAddress())

	for _, dataTx := range dataTxs {
		io.scriptIo.blockManager.SendDataTx(tx, dataTx, nil)
	}

	return tx, dataTxs
}

func (io *ScriptIoHandler) UpdateScriptData(uniqKey, data []byte) (key []byte) {
	prevMap := map[types.TxOutputType][]types.PrevOutputParam{}
	prevMap[types.TxTypeScriptStore] = io.scriptTxPtr // for mapping

	txInfo := &txs.TransferTxInfo{
		Prevs:     prevMap,
		FromAddr:  io.wallet.MyPubKey(),
		ToAddr:    io.toAddr,
		Broker:    io.brokerAddr, //io.wallet.GetMasterNodeAddr(),
		FeeAddr:   store.AdamsAddress(),
		FeeBroker: io.feeBrokerAddr,
	}
	tx, dataTx := io.scriptIo.tXs.CreateScriptDataTx(*txInfo, uniqKey, data)
	tx = io.scriptIo.tXs.InkTheContract(tx, io.wallet.GetGhostAddress())

	io.scriptIo.blockManager.SendDataTx(tx, dataTx, nil)

	key = dataTx.TxId
	return key
}

func (io *ScriptIo) loadRefTx(key, toAddr []byte) (ret []types.PrevOutputParam) {
	ret = io.bc.TxContainer.SelectOutputLatests(types.TxTypeScriptStore, toAddr, key, 0, 1)
	return ret
}

func (io *ScriptIo) loadDataTx(dataTxId []byte) *types.GhostDataTransaction {
	dataTx := io.bc.TxContainer.GetDataTx(dataTxId)
	if dataTx == nil {
		fileObj := io.cloud.ReadFromCloudSync(fileservice.ByteToFilename(dataTxId),
			io.wallet.GetMasterNode().Ip.GetUdpAddr())
		if fileObj == nil {
			log.Print("could not found file: ", fileObj.Filename)
			return nil
		}

		dataTx = &types.GhostDataTransaction{}
		if !dataTx.Deserialize(bytes.NewBuffer(fileObj.Buffer)).Result() {
			return nil
		}
	}
	return dataTx
}
