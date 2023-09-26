package txs

import (
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

func (txs *TXs) CandidateUTXO(withDrawCoin uint64, account []byte) ([]types.PrevOutputParam, bool) {
	outputParams := []types.PrevOutputParam{}
	getherCoin := uint64(0)
	checkBalance := false

	for _, outputParam := range txs.blockContainer.TxContainer.GetUnusedOutputList(types.TxTypeCoinTransfer, account) {
		getherCoin += outputParam.Vout.Value
		outputParams = append(outputParams, outputParam)
		if getherCoin >= withDrawCoin {
			checkBalance = true
			break
		}
	}

	return outputParams, checkBalance
}

func (txs *TXs) TotalUTXO(account []byte) uint64 {
	getherCoin := uint64(0)
	for _, outputParam := range txs.blockContainer.TxContainer.GetUnusedOutputList(types.TxTypeCoinTransfer, account) {
		getherCoin += outputParam.Vout.Value
	}
	return getherCoin
}

func (txs *TXs) GetRootFsTx(account []byte) ([]types.PrevOutputParam, bool) {
	outputParams := txs.blockContainer.TxContainer.SearchOutputList(types.TxTypeFSRoot, account)
	if len(outputParams) != 0 {
		return outputParams, true
	}
	return nil, false
}

func (txs *TXs) GetDataTx(key, toAddr []byte) ([]types.PrevOutputParam, bool) {
	return nil, false
}