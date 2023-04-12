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

// cli(GrpcClient) -> cli(rpc.GApiClient) ->
// ghostd(GrpcServer) -> ghostd(GrpcDeamonHandler) -> container(GrpcClient) ->
// container(GrpcServer) -> container(GhostContainerApi)
type GrpcDeamonHandler struct {
	loadWallet *bootloader.LoadWallet
	genesis    *bootloader.LoadGenesis
	containers *Containers
	grpcServer *grpc.GrpcServer
	config     *gconfig.GConfig
}

func NewGrpcDeamonHandler(loadWallet *bootloader.LoadWallet, genesis *bootloader.LoadGenesis,
	containers *Containers, grpcServer *grpc.GrpcServer, config *gconfig.GConfig) *GrpcDeamonHandler {
	gHandler := &GrpcDeamonHandler{
		loadWallet: loadWallet,
		genesis:    genesis,
		containers: containers,
		grpcServer: grpcServer,
		config:     config,
	}
	grpcServer.CreateAccountHandler = gHandler.CreateAccountDeamon
	grpcServer.CreateGenesisHandler = gHandler.CreateGenesisDeamon
	grpcServer.LoginContainerHandler = gHandler.LoginContainerDeamon
	grpcServer.ForkContainerHandler = gHandler.ForkContainerDeamon
	grpcServer.CreateContainerHandler = gHandler.CreateContainerDeamon
	grpcServer.ControlContainerHandler = gHandler.ControlContainerDeamon
	grpcServer.ReleaseContainerHandler = gHandler.ReleaseContainerDeamon
	grpcServer.GetPrivateKeyHandler = gHandler.GetPrivateKeyDeamon
	grpcServer.GetLogHandler = gHandler.GetLogDeamon
	grpcServer.CheckStatusHandler = gHandler.CheckStatusDeamon
	grpcServer.GetContainerListHandler = gHandler.GetContainerListDeamon
	grpcServer.GetInfoHandler = gHandler.GetInfoDeamon
	grpcServer.GetAccountHandler = gHandler.GetAccountDeamon
	grpcServer.GetBlockInfoHandler = gHandler.GetBlockInfoDeamon
	return gHandler
}

func (ghandler *GrpcDeamonHandler) GetInfoDeamon() *rpc.GetInfoResponse {
	return &rpc.GetInfoResponse{
		TotalContainer: ghandler.containers.Count,
	}
}

func (ghandler *GrpcDeamonHandler) GetContainerListDeamon(id uint32) *rpc.GetContainerListResponse {
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

func (ghandler *GrpcDeamonHandler) CreateAccountDeamon(password []byte, username string) bool {
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

func (ghandler *GrpcDeamonHandler) CreateGenesisDeamon(id uint32, password []byte) bool {
	if container, exist := ghandler.containers.GetContainer(id); exist {
		return container.Client.CreateGenesis(id, password)
	}
	return false
}

func (ghandler *GrpcDeamonHandler) ReleaseContainerDeamon(id uint32) bool {
	return ghandler.containers.ReleaseContainer(id)
}

func (ghandler *GrpcDeamonHandler) GetPrivateKeyDeamon(id uint32, password []byte, username string) ([]byte, bool) {
	w, err := ghandler.loadWallet.OpenWallet(username, password)
	if err != nil {
		log.Print(err)
		return nil, false
	}
	privateKey := w.GetGhostAddress().PrivateKeySerialize()
	cipherKey := gcrypto.Encryption(password, privateKey)
	return cipherKey, true
}

func (ghandler *GrpcDeamonHandler) LoginContainerDeamon(id uint32, password []byte, username, ip, port string) bool {
	creators := ghandler.genesis.CreatorList()
	if creator, exist := creators[username]; exist {
		if _, err := ghandler.genesis.LoadCreatorKeyFile(creator.Nickname,
			creator.PubKey, password); err != nil {
			log.Println("Load Creator Key File Fail..")
			return false
		}

	} else {
		if w, err := ghandler.loadWallet.OpenWallet(username, password); w == nil {
			log.Print("not exist account = ", username, " or err = ", err)
			return false
		}
	}

	return ghandler.containers.LoginContainer(id, password, username, ip, port) != nil
}

func (ghandler *GrpcDeamonHandler) ForkContainerDeamon(password []byte, username, ip, port string) bool {
	log.Print("CreateContainerHandler")
	if w, err := ghandler.loadWallet.OpenWallet(username, password); w == nil {
		log.Print("not exist account = ", username, "or err = ", err)
		return false
	}
	return ghandler.containers.ForkContainer(password, username, ip, port) != nil
}

func (ghandler *GrpcDeamonHandler) CreateContainerDeamon(password []byte, username, ip, port string) bool {
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

func (ghandler *GrpcDeamonHandler) ControlContainerDeamon(id uint32, control rpc.ContainerControlType) bool {
	if container, exist := ghandler.containers.GetContainer(id); exist {
		return container.Client.ControlContainer(id, control)
	}
	return false
}

func (ghandler *GrpcDeamonHandler) GetLogDeamon(id uint32) []byte {
	if container, exist := ghandler.containers.GetContainer(id); exist {
		return container.Client.GetLog(id)
	}
	return nil
}

func (ghandler *GrpcDeamonHandler) CheckStatusDeamon(id uint32) uint32 {
	if container, exist := ghandler.containers.GetContainer(id); exist {
		return container.Client.CheckStatus(id)
	}
	return 0
}

func (ghandler *GrpcDeamonHandler) GetAccountDeamon(id uint32) (users []*ptypes.GhostUser) {
	names := ghandler.loadWallet.GetWalletList()
	for _, name := range names {
		users = append(users, &ptypes.GhostUser{Nickname: name})
	}
	return users
}

func (ghandler *GrpcDeamonHandler) GetBlockInfoDeamon(id, blockId uint32) *ptypes.PairedBlocks {
	if container, exist := ghandler.containers.GetContainer(id); exist {
		return container.Client.GetBlockInfo(id, blockId)
	}
	return nil
}
