package blockfilesystem

import (
	"bytes"
	"fmt"
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
	"github.com/GhostNet-Dev/glambda/evaluator"
	"github.com/GhostNet-Dev/glambda/object"
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
	evaluator.AddBuiltIn("loadKeyValue", &object.Builtin{
		Fn: func(env interface{}, args ...object.Object) object.Object {
			if len(args) != 1 {
				return object.NewError("wrong number of arguments. got=%d, want=1", len(args))
			}
			var key string
			switch arg := args[0].(type) {
			case *object.String:
				key = arg.Value
			case *object.Identifier:
				key = arg.Value.(*object.String).Value
			case *object.Integer:
				key = fmt.Sprint(arg.Value)
			}
			data := gScriptIo.scriptIoHandler.ReadScriptData([]byte(key))
			return &object.String{Value: string(data)}
		},
	})

	evaluator.AddBuiltIn("saveKeyValue", &object.Builtin{
		Fn: func(env interface{}, args ...object.Object) object.Object {
			if len(args) != 2 {
				return object.NewError("wrong number of arguments. got=%d, want=2", len(args))
			}
			if args[0].Type() != object.STRING_OBJ || args[1].Type() != object.STRING_OBJ {
				return object.NewError("argument to 'saveKeyValue' must be string type. got=%s", args[0].Type())
			}
			key := []byte(args[0].(*object.String).Value)
			value := []byte(args[1].(*object.String).Value)
			// TODO: Script 전용 tx storage가 만들어져야함
			txId := gScriptIo.scriptIoHandler.WriteScriptData(key, value)
			return &object.String{Value: string(txId)}
		},
	})
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

	fileObj := io.cloud.ReadFromCloudSync(fileservice.ByteToFilename(dataTxId),
		io.wallet.GetMasterNode().Ip.GetUdpAddr())
	if fileObj == nil {
		log.Print("could not found file: ", fileObj.Filename)
		return nil
	}

	dataTx := &types.GhostDataTransaction{}
	if !dataTx.Deserialize(bytes.NewBuffer(fileObj.Buffer)).Result() {
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
	fileObj := io.scriptIo.cloud.ReadFromCloudSync(fileservice.ByteToFilename(dataTxId),
		io.wallet.GetMasterNode().Ip.GetUdpAddr())
	if fileObj == nil {
		log.Print("could not found file: ", fileObj.Filename)
		return nil
	}
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
