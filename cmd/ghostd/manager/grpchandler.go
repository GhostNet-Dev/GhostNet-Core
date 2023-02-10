package manager

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/bootloader"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/grpc"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/rpc"
)

type GrpcHandler struct {
	containers *Containers
	grpcServer *grpc.GrpcServer
	config     *gconfig.GConfig
}

func NewGrpcHandler(containers *Containers, grpcServer *grpc.GrpcServer,
	config *gconfig.GConfig) *GrpcHandler {
	gHandler := &GrpcHandler{
		containers: containers,
		grpcServer: grpcServer,
		config:     config,
	}
	grpcServer.CreateAccountHandler = gHandler.CreateAccountHandler
	grpcServer.CreateGenesisHandler = gHandler.CreateGenesisHandler
	grpcServer.CreateContainerHandler = gHandler.CreateContainerHandler
	grpcServer.ControlContainerHandler = gHandler.ControlContainerHandler
	grpcServer.ReleaseContainerHandler = gHandler.ReleaseContainerHandler
	grpcServer.GetPrivateKeyHandler = gHandler.GetPrivateKeyHandler
	grpcServer.GetLogHandler = gHandler.GetLogHandler
	grpcServer.CheckStatusHandler = gHandler.CheckStatusHandler
	grpcServer.GetContainerListHandler = gHandler.GetContainerListHandler
	grpcServer.GetInfoHandler = gHandler.GetInfoHandler
	return gHandler
}

func (ghandler *GrpcHandler) GetInfoHandler() *rpc.GetInfoResponse {
	return &rpc.GetInfoResponse{
		TotalContainer: ghandler.containers.Count,
	}
}

func (ghandler *GrpcHandler) GetContainerListHandler(id uint32) *rpc.GetContainerListResponse {
	container := ghandler.containers.GetContainer(id)
	return &rpc.GetContainerListResponse{
		Id:       container.Id,
		PID:      int32(container.PID),
		PubKey:   container.PubKey,
		Port:     container.Port,
		Username: container.Username,
	}
}

func (ghandler *GrpcHandler) CreateAccountHandler(password []byte, username string) bool {
	db := bootloader.NewLiteStore(ghandler.config.SqlPath, bootloader.Tables)
	wallet := bootloader.NewLoadWallet(bootloader.Tables[1], db, nil)
	_, err := wallet.OpenWallet(username, password)
	if err != nil {
		w := wallet.CreateWallet(username, password)
		wallet.SaveWallet(w, password)
	} else {
		log.Print("already exist account")
		return false
	}
	return true
}

func (ghandler *GrpcHandler) CreateGenesisHandler(id uint32, password []byte) bool {
	container := ghandler.containers.GetContainer(id)
	return container.Client.CreateGenesis(id, password)
}

func (ghandler *GrpcHandler) ReleaseContainerHandler(id uint32) bool {
	return ghandler.containers.ReleaseContainer(id)
}

func (ghandler *GrpcHandler) GetPrivateKeyHandler(id uint32, password []byte, username string) ([]byte, bool) {
	db := bootloader.NewLiteStore(ghandler.config.SqlPath, bootloader.Tables)
	wallet := bootloader.NewLoadWallet(bootloader.Tables[1], db, nil)
	w, err := wallet.OpenWallet(username, password)
	if err != nil {
		log.Print(err)
		return nil, false
	}
	privateKey := w.GetGhostAddress().PrivateKeySerialize()
	cipherKey := wallet.Encryption(password, privateKey)
	return cipherKey, true
}

func (ghandler *GrpcHandler) CreateContainerHandler(password []byte, username, ip, port string) bool {
	log.Print("CreateContainerHandler")
	return ghandler.containers.CreateContainer(password, username, ip, port) != nil
}

func (ghandler *GrpcHandler) ControlContainerHandler(id uint32, control rpc.ContainerControlType) bool {
	container := ghandler.containers.GetContainer(id)
	return container.Client.ControlContainer(id, control)
}

func (ghandler *GrpcHandler) GetLogHandler(id uint32) []byte {
	container := ghandler.containers.GetContainer(id)
	return container.Client.GetLog(id)
}

func (ghandler *GrpcHandler) CheckStatusHandler(id uint32) uint32 {
	container := ghandler.containers.GetContainer(id)
	return container.Client.CheckStatus(id)
}
