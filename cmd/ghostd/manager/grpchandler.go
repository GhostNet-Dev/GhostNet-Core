package manager

import (
	"log"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/bootloader"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/grpc"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/rpc"
)

type GrpcHandler struct {
	loadWallet *bootloader.LoadWallet
	genesis    *bootloader.LoadGenesis
	containers *Containers
	grpcServer *grpc.GrpcServer
	config     *gconfig.GConfig
}

func NewGrpcHandler(loadWallet *bootloader.LoadWallet, genesis *bootloader.LoadGenesis,
	containers *Containers, grpcServer *grpc.GrpcServer, config *gconfig.GConfig) *GrpcHandler {
	gHandler := &GrpcHandler{
		loadWallet: loadWallet,
		genesis:    genesis,
		containers: containers,
		grpcServer: grpcServer,
		config:     config,
	}
	grpcServer.CreateAccountHandler = gHandler.CreateAccountHandler
	grpcServer.CreateGenesisHandler = gHandler.CreateGenesisHandler
	grpcServer.LoginContainerHandler = gHandler.LoginContainerHandler
	grpcServer.ForkContainerHandler = gHandler.ForkContainerHandler
	grpcServer.CreateContainerHandler = gHandler.CreateContainerHandler
	grpcServer.ControlContainerHandler = gHandler.ControlContainerHandler
	grpcServer.ReleaseContainerHandler = gHandler.ReleaseContainerHandler
	grpcServer.GetPrivateKeyHandler = gHandler.GetPrivateKeyHandler
	grpcServer.GetLogHandler = gHandler.GetLogHandler
	grpcServer.CheckStatusHandler = gHandler.CheckStatusHandler
	grpcServer.GetContainerListHandler = gHandler.GetContainerListHandler
	grpcServer.GetInfoHandler = gHandler.GetInfoHandler
	grpcServer.GetAccountHandler = gHandler.GetAccountHandler
	grpcServer.GetBlockInfoHandler = gHandler.GetBlockInfoHandler
	return gHandler
}

func (ghandler *GrpcHandler) GetInfoHandler() *rpc.GetInfoResponse {
	return &rpc.GetInfoResponse{
		TotalContainer: ghandler.containers.Count,
	}
}

func (ghandler *GrpcHandler) GetContainerListHandler(id uint32) *rpc.GetContainerListResponse {
	container, exist := ghandler.containers.GetContainer(id)
	if !exist {
		return nil
	}
	return &rpc.GetContainerListResponse{
		Id:       container.Id,
		PID:      int32(container.PID),
		PubKey:   container.PubKey,
		Port:     container.Port,
		Username: container.Username,
	}
}

func (ghandler *GrpcHandler) CreateAccountHandler(password []byte, username string) bool {
	w, _ := ghandler.loadWallet.OpenWallet(username, password)
	if w == nil {
		w = ghandler.loadWallet.CreateWallet(username, password)
		ghandler.loadWallet.SaveWallet(w, password)
	} else {
		log.Print("already exist account")
		return false
	}
	return true
}

func (ghandler *GrpcHandler) CreateGenesisHandler(id uint32, password []byte) bool {
	if container, exist := ghandler.containers.GetContainer(id); exist {
		return container.Client.CreateGenesis(id, password)
	}
	return false
}

func (ghandler *GrpcHandler) ReleaseContainerHandler(id uint32) bool {
	return ghandler.containers.ReleaseContainer(id)
}

func (ghandler *GrpcHandler) GetPrivateKeyHandler(id uint32, password []byte, username string) ([]byte, bool) {
	w, err := ghandler.loadWallet.OpenWallet(username, password)
	if err != nil {
		log.Print(err)
		return nil, false
	}
	privateKey := w.GetGhostAddress().PrivateKeySerialize()
	cipherKey := gcrypto.Encryption(password, privateKey)
	return cipherKey, true
}

func (ghandler *GrpcHandler) LoginContainerHandler(password []byte, username, ip, port string) bool {
	if w, err := ghandler.loadWallet.OpenWallet(username, password); w == nil {
		log.Print("not exist account = ", username, "or err = ", err)
		return false
	}
	return ghandler.containers.LoginContainer(password, username, ip, port) != nil
}

func (ghandler *GrpcHandler) ForkContainerHandler(password []byte, username, ip, port string) bool {
	log.Print("CreateContainerHandler")
	if w, err := ghandler.loadWallet.OpenWallet(username, password); w == nil {
		log.Print("not exist account = ", username, "or err = ", err)
		return false
	}
	return ghandler.containers.ForkContainer(password, username, ip, port) != nil
}

func (ghandler *GrpcHandler) CreateContainerHandler(password []byte, username, ip, port string) bool {
	log.Print("CreateContainerHandler")
	creators := ghandler.genesis.CreatorList()
	if _, exist := creators[username]; !exist {
		if w, err := ghandler.loadWallet.OpenWallet(username, password); w == nil {
			log.Print("not exist account = ", username, " or err = ", err)
			return false
		}
	}

	return ghandler.containers.CreateContainer(password, username, ip, port) != nil
}

func (ghandler *GrpcHandler) ControlContainerHandler(id uint32, control rpc.ContainerControlType) bool {
	if container, exist := ghandler.containers.GetContainer(id); exist {
		return container.Client.ControlContainer(id, control)
	}
	return false
}

func (ghandler *GrpcHandler) GetLogHandler(id uint32) []byte {
	if container, exist := ghandler.containers.GetContainer(id); exist {
		return container.Client.GetLog(id)
	}
	return nil
}

func (ghandler *GrpcHandler) CheckStatusHandler(id uint32) uint32 {
	if container, exist := ghandler.containers.GetContainer(id); exist {
		return container.Client.CheckStatus(id)
	}
	return 0
}

func (ghandler *GrpcHandler) GetAccountHandler(id uint32) (users []*ptypes.GhostUser) {
	names := ghandler.loadWallet.GetWalletList()
	for _, name := range names {
		users = append(users, &ptypes.GhostUser{Nickname: name})
	}
	return users
}

func (ghandler *GrpcHandler) GetBlockInfoHandler(id, blockId uint32) *ptypes.PairedBlocks {
	if container, exist := ghandler.containers.GetContainer(id); exist {
		return container.Client.GetBlockInfo(id, blockId)
	}
	return nil
}
