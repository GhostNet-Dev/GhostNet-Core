package blockfilesystem

import (
	"log"
	"sync"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blockmanager"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/cloudservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/fileservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/txs"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type BlockIo struct {
	blockManager *blockmanager.BlockManager
	bc           *store.BlockContainer
	tXs          *txs.TXs
	cloud        *cloudservice.CloudService
}

type BlockIoHandler struct {
	wallet    *gcrypto.Wallet
	rootTxPtr []types.PrevOutputParam
	blockIo   *BlockIo
}

func NewBlockIo(blkMgr *blockmanager.BlockManager,
	bc *store.BlockContainer, tXs *txs.TXs, cloud *cloudservice.CloudService) *BlockIo {
	return &BlockIo{
		blockManager: blkMgr,
		bc:           bc,
		tXs:          tXs,
		cloud:        cloud,
	}
}

func (io *BlockIo) CreateFilesystem(w *gcrypto.Wallet) *BlockIoHandler {
	txInfo := &txs.TransferTxInfo{
		MyWallet:  w,
		ToAddr:    w.MyPubKey(),
		Broker:    w.GetMasterNodeAddr(),
		FeeAddr:   store.AdamsAddress(),
		FeeBroker: w.GetMasterNodeAddr(),
	}
	tx := io.tXs.CreateRootFsTx(*txInfo, w.GetNickname())
	tx = io.tXs.InkTheContract(tx, w.GetGhostAddress())
	wg := &sync.WaitGroup{}
	result := false
	io.blockManager.SendTx(tx, func(b bool) {
		defer wg.Done()
		result = b
	})
	wg.Wait()

	if !result {
		return nil
	}
	return &BlockIoHandler{
		wallet:  w,
		blockIo: io,
	}
}

func (io *BlockIo) OpenFilesystem(w *gcrypto.Wallet) *BlockIoHandler {
	if outputParams, ok := io.tXs.GetRootFsTx(w.MyPubKey()); ok {
		return &BlockIoHandler{
			wallet:    w,
			rootTxPtr: outputParams,
			blockIo:   io,
		}
	}
	return nil
}

func (io *BlockIo) CloseFilesystem() {
}

func (io *BlockIoHandler) ReadData(key []byte) []byte {
	fileObj := io.blockIo.cloud.ReadFromCloudSync(fileservice.ByteToFilename(key),
		io.wallet.GetMasterNode().Ip.GetUdpAddr())
	if fileObj == nil {
		log.Print("could not found file: ", fileObj.Filename)
		return nil
	}
	return fileObj.Buffer
}

func (io *BlockIoHandler) WriteData(uniqKey, data []byte) (key []byte) {
	prevMap := map[types.TxOutputType][]types.PrevOutputParam{}
	prevMap[types.TxTypeDataStore] = io.rootTxPtr // for mapping

	txInfo := &txs.TransferTxInfo{
		Prevs:     prevMap,
		MyWallet:  io.wallet,
		ToAddr:    io.wallet.MyPubKey(),
		Broker:    io.wallet.GetMasterNodeAddr(),
		FeeAddr:   store.AdamsAddress(),
		FeeBroker: io.wallet.GetMasterNodeAddr(),
	}
	tx, dataTx := io.blockIo.tXs.CreateDataTx(*txInfo, uniqKey, data)
	tx = io.blockIo.tXs.InkTheContract(tx, io.wallet.GetGhostAddress())

	io.blockIo.blockManager.SendDataTx(tx, dataTx, nil)

	key = dataTx.TxId
	return key
}

func (io *BlockIoHandler) WriteDataBrokerService(uniqKey, data, feeBroker []byte) (key []byte) {
	prevMap := map[types.TxOutputType][]types.PrevOutputParam{}
	prevMap[types.TxTypeDataStore] = io.rootTxPtr // for mapping

	txInfo := &txs.TransferTxInfo{
		Prevs:     prevMap,
		MyWallet:  io.wallet,
		ToAddr:    io.wallet.MyPubKey(),
		Broker:    io.wallet.GetMasterNodeAddr(),
		FeeAddr:   store.AdamsAddress(),
		FeeBroker: feeBroker,
	}
	tx, dataTx := io.blockIo.tXs.CreateDataTx(*txInfo, uniqKey, data)
	tx = io.blockIo.tXs.InkTheContract(tx, io.wallet.GetGhostAddress())

	io.blockIo.blockManager.SendDataTx(tx, dataTx, nil)

	key = dataTx.TxId
	return key
}
