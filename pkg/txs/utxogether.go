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
