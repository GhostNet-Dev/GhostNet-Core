package gapi

import (
	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blocks"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/bootloader"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/grpc"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/rpc"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

type GhostApi struct {
	grpcServer     *grpc.GrpcServer
	block          *blocks.Blocks
	blockContainer *store.BlockContainer
	loadWallet     *bootloader.LoadWallet
	config         *gconfig.GConfig
}

func NewGhostApi(grpcServer *grpc.GrpcServer, block *blocks.Blocks,
	blockContainer *store.BlockContainer,
	loadWallet *bootloader.LoadWallet, config *gconfig.GConfig) *GhostApi {
	ghostApi := &GhostApi{
		grpcServer:     grpcServer,
		block:          block,
		blockContainer: blockContainer,
		loadWallet:     loadWallet,
		config:         config,
	}

	grpcServer.CreateGenesisHandler = ghostApi.CreateGenesisHandler
	grpcServer.ControlContainerHandler = ghostApi.ControlContainerHandler
	grpcServer.GetLogHandler = ghostApi.GetLogHandler
	grpcServer.GetBlockInfoHandler = ghostApi.GetBlockInfoHandler
	return ghostApi
}

func (gApi *GhostApi) CreateGenesisHandler(id uint32, password []byte) bool {
	gApi.block.MakeGenesisBlock(func(name string, address *gcrypto.GhostAddress) {
		gApi.loadWallet.SaveWallet(gcrypto.NewWallet(name, address, nil, nil), password)
	})
	return false
}

func (gApi *GhostApi) ControlContainerHandler(id uint32, control rpc.ContainerControlType) bool {
	return false
}

func (gApi *GhostApi) GetLogHandler(id uint32) []byte {
	return nil
}

func (gApi *GhostApi) CheckStatusHandler(id uint32) uint32 {
	return 0
}

func (gApi *GhostApi) GetBlockInfoHandler(id, blockId uint32) *ptypes.PairedBlocks {
	pair := gApi.blockContainer.GetBlock(blockId)
	protoPairedBlock := types.GhostBlockToProtoType(pair)
	return protoPairedBlock
}
