package blockfilesystem

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blockmanager"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/cloudservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/fileservice"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/txs"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type BlockIo struct {
	wallet       *gcrypto.Wallet
	blockManager *blockmanager.BlockManager
	bc           *store.BlockContainer
	tXs          *txs.TXs
	cloud        *cloudservice.CloudService
	rootTxPtr    []types.PrevOutputParam
}

func NewBlockIo(w *gcrypto.Wallet, blkMgr *blockmanager.BlockManager,
	bc *store.BlockContainer, tXs *txs.TXs, cloud *cloudservice.CloudService) *BlockIo {
	return &BlockIo{
		wallet:       w,
		blockManager: blkMgr,
		bc:           bc,
		tXs:          tXs,
		cloud:        cloud,
	}
}

func (io *BlockIo) CreateFilesystem() {
	txInfo := &txs.TransferTxInfo{
		MyWallet:  io.wallet,
		ToAddr:    io.wallet.MyPubKey(),
		Broker:    io.wallet.GetMasterNodeAddr(),
		FeeAddr:   store.AdamsAddress(),
		FeeBroker: io.wallet.GetMasterNodeAddr(),
	}
	tx := io.tXs.CreateRootFsTx(*txInfo, io.wallet.GetNickname())
	tx = io.tXs.InkTheContract(tx, io.wallet.GetGhostAddress())
	io.blockManager.SendTx(tx)
}

func (io *BlockIo) OpenFilesystem() bool {
	if outputParams, ok := io.tXs.GetRootFsTx(io.wallet.MyPubKey()); ok {
		io.rootTxPtr = outputParams
		return true
	}
	return false
}

func (io *BlockIo) CloseFilesystem() {
}

func (io *BlockIo) ReadData(key []byte) []byte {
	fileObj := io.cloud.ReadFromCloudSync(fileservice.ByteToFilename(key),
		io.wallet.GetMasterNode().Ip.GetUdpAddr())
	if fileObj == nil {
		log.Print("could not found file: ", fileObj.Filename)
		return nil
	}
	return fileObj.Buffer
}

func (io *BlockIo) WriteData(uniqKey []byte, data []byte) (key []byte) {
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
	tx, dataTx := io.tXs.CreateDataTx(*txInfo, uniqKey, data)
	tx = io.tXs.InkTheContract(tx, io.wallet.GetGhostAddress())

	io.blockManager.SendDataTx(tx, dataTx)

	key = dataTx.TxId
	return key
}
