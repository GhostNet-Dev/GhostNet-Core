package gapi

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/blocks"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/bootloader"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/grpc"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/rpc"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/store"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/types"
)

// cli(GrpcClient) -> cli(rpc.GApiClient) ->
// ghostd(GrpcServer) -> ghostd(GrpcDeamonHandler) -> container(GrpcClient) ->
// container(GrpcServer) -> container(GhostContainerApi)
type GhostContainerApi struct {
	grpcServer     *grpc.GrpcServer
	block          *blocks.Blocks
	blockContainer *store.BlockContainer
	loadWallet     *bootloader.LoadWallet
	eventListener  map[rpc.ContainerControlType][]func(rpc.ContainerControlType, interface{})
}

func NewGhostContainerApi(grpcServer *grpc.GrpcServer, block *blocks.Blocks,
	blockContainer *store.BlockContainer,
	loadWallet *bootloader.LoadWallet) *GhostContainerApi {
	ghostApi := &GhostContainerApi{
		grpcServer:     grpcServer,
		block:          block,
		blockContainer: blockContainer,
		loadWallet:     loadWallet,
		eventListener:  make(map[rpc.ContainerControlType][]func(rpc.ContainerControlType, interface{})),
	}

	grpcServer.CreateGenesisHandler = ghostApi.CreateGenesisHandler
	grpcServer.ControlContainerHandler = ghostApi.ControlContainerHandler
	grpcServer.LoginContainerHandler = ghostApi.LoginContainerHandler
	grpcServer.GetLogHandler = ghostApi.GetLogHandler
	grpcServer.GetBlockInfoHandler = ghostApi.GetBlockInfoHandler
	return ghostApi
}

func (gApi *GhostContainerApi) RegisterEventListener(control rpc.ContainerControlType, handler func(rpc.ContainerControlType, interface{})) {
	gApi.eventListener[control] = append(gApi.eventListener[control], handler)
}

func (gApi *GhostContainerApi) CreateGenesisHandler(id uint32, password []byte) bool {
	gApi.block.MakeGenesisBlock(func(name string, address *gcrypto.GhostAddress) {
		gApi.loadWallet.SaveWallet(gcrypto.NewWallet(name, address, nil, nil), password)
	})
	return false
}

func (gApi *GhostContainerApi) LoginContainerHandler(id uint32, password []byte, username, ip, port string) bool {
	if w, _ := gApi.loadWallet.OpenWallet(username, password); w == nil {
		log.Println("Login fail user = ", username)
		return false
	}
	control := rpc.ContainerControlType_StartResume
	for _, handler := range gApi.eventListener[control] {
		handler(control, struct {
			Nickname string
			PassHash []byte
		}{Nickname: username, PassHash: password})
	}
	return true
}

func (gApi *GhostContainerApi) ControlContainerHandler(id uint32, control rpc.ContainerControlType) bool {
	for _, handler := range gApi.eventListener[control] {
		handler(control, nil)
	}
	return true
}

func (gApi *GhostContainerApi) GetLogHandler(id uint32) []byte {
	return nil
}

func (gApi *GhostContainerApi) CheckStatusHandler(id uint32) uint32 {
	return 0
}

func (gApi *GhostContainerApi) GetBlockInfoHandler(id, blockId uint32) *ptypes.PairedBlocks {
	pair := gApi.blockContainer.GetBlock(blockId)
	if pair == nil {
		return nil
	}
	protoPairedBlock := types.GhostBlockToProtoType(pair)
	return protoPairedBlock
}
