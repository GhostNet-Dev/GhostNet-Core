package blockfilesystem

import (
	"bytes"
	"sync"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blockmanager"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/cloudservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/fileservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/txs"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type ScriptIo struct {
	bc           *store.BlockContainer
	cloud        *cloudservice.CloudService
	wallet       *gcrypto.Wallet
	blockManager *blockmanager.BlockManager
	liteStore    *store.LiteStore
	tXs          *txs.TXs
}

type ScriptIoHandler struct {
	wallet        *gcrypto.Wallet
	scriptIo      *ScriptIo
	toAddr        []byte
	brokerAddr    []byte
	feeBrokerAddr []byte
	scriptTxPtr   []types.PrevOutputParam
}

/*
todo list
1. make litestore uniqkey -> datatxid map
2. when called open, read script tx
3. read key ->  search key map with datatxid
*/

func NewScriptIo(w *gcrypto.Wallet, liteStore *store.LiteStore, cloud *cloudservice.CloudService,
	bc *store.BlockContainer, blkMgr *blockmanager.BlockManager, tXs *txs.TXs) *ScriptIo {
	return &ScriptIo{
		wallet:       w,
		liteStore:    liteStore,
		cloud:        cloud,
		bc:           bc,
		blockManager: blkMgr,
		tXs:          tXs,
	}
}

func (io *ScriptIo) CreateScript(w *gcrypto.Wallet, namespace, script string) ([]byte, bool) {
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
	tx, dataTx := io.tXs.CreateScriptTx(*txInfo, []byte(namespace), []byte(script))
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
	output := tx.Body.Vout[0]
	if output.Type != types.TxTypeScript {
		return nil
	}
	dataTxId := output.ScriptEx

	fileObj := io.cloud.ReadFromCloudSync(fileservice.ByteToFilename(dataTxId),
		io.wallet.GetMasterNode().Ip.GetUdpAddr())

	dataTx := &types.GhostDataTransaction{}
	if !dataTx.Deserialize(bytes.NewBuffer(fileObj.Buffer)).Result() {
		return nil
	}
	return &ScriptIoHandler{
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
	}
}

func (io *ScriptIo) CloseScript(handler *ScriptIoHandler) {}

func (io *ScriptIoHandler) ReadScriptData(key []byte) (data []byte) {
	// to avoid key collision
	uniqKey := append(io.toAddr, key...)
	dataTxId, err := io.scriptIo.liteStore.SelectEntry(store.DefaultDataKeyMappingTable, uniqKey)
	if err != nil {
		return nil
	}

	fileObj := io.scriptIo.cloud.ReadFromCloudSync(fileservice.ByteToFilename(dataTxId),
		io.wallet.GetMasterNode().Ip.GetUdpAddr())

	dataTx := &types.GhostDataTransaction{}
	if !dataTx.Deserialize(bytes.NewBuffer(fileObj.Buffer)).Result() {
		return nil
	}
	return dataTx.Data
}

func (io *ScriptIoHandler) WriteScriptData(uniqKey, data []byte) (key []byte) {
	prevMap := map[types.TxOutputType][]types.PrevOutputParam{}
	prevMap[types.TxTypeScriptStore] = io.scriptTxPtr // for mapping

	txInfo := &txs.TransferTxInfo{
		Prevs:     prevMap,
		FromAddr:  io.wallet.MyPubKey(),
		ToAddr:    io.toAddr,
		Broker:    io.wallet.GetMasterNodeAddr(),
		FeeAddr:   store.AdamsAddress(),
		FeeBroker: io.feeBrokerAddr,
	}
	tx, dataTx := io.scriptIo.tXs.CreateScriptDataTx(*txInfo, uniqKey, data)
	tx = io.scriptIo.tXs.InkTheContract(tx, io.wallet.GetGhostAddress())

	io.scriptIo.blockManager.SendDataTx(tx, dataTx, nil)

	key = dataTx.TxId
	return key
}
