package consensus

import "github.com/GhostNet-Dev/GhostNet-Core/pkg/types"

type EventHandler interface {
	MiningStart()
	MiningStop()
	BroadcastBlockChainNotification()
	CheckBlockHeight()
	SetHeighestCandidatePool()
	LoadHashFromTempDb(blockIdx uint32) string
	RequestGetBlock(blockIdx uint32)
	RequestGetBlockHash(currBlockHeight uint32)
	CheckAndSave(pairedBlock *types.PairedBlock) bool
	MergeErrorNotification()
	LocalBlockDbValidation() bool
}
