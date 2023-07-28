package grpc

import (
	"context"
	"log"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/ptypes"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// cli(GrpcClient) -> cli(proto.rpc.GApiClient) ->
// ghostd(grpc.GrpcServer) -> ghostd(manager.GrpcDeamonHandler) ->
// container(grpc.GrpcClient) ->
// container(grpc.GrpcServer) -> container(gapi.GhostContainerApi)
type GrpcClient struct {
	GrpcIp   string
	GrpcPort string
	c        rpc.GApiClient
	conn     *grpc.ClientConn
	Timeout  uint32
}

func NewGrpcClient(grpcIp, grpcPort string, timeout uint32) *GrpcClient {
	return &GrpcClient{
		GrpcIp:   grpcIp,
		GrpcPort: grpcPort,
		Timeout:  timeout,
	}
}

func (client *GrpcClient) ConnectServer() {
	conn, err := grpc.Dial(client.GrpcIp+":"+client.GrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("did not connect: %v", err)
	}
	client.c = rpc.NewGApiClient(conn)
	client.conn = conn
}

func (client *GrpcClient) CloseServer() {
	client.conn.Close()
}

func (client *GrpcClient) GetInfo() *rpc.GetInfoResponse {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(client.Timeout))
	defer cancel()
	r, err := client.c.GetInfo(ctx, &rpc.GetInfoRequest{})
	if err != nil {
		log.Printf("could not connect: %v", err)
	}
	return r
}

func (client *GrpcClient) GetContainerList(id uint32) *rpc.GetContainerListResponse {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(client.Timeout))
	defer cancel()
	r, err := client.c.GetContainerList(ctx, &rpc.GetContainerListRequest{Id: id})
	if err != nil {
		log.Printf("could not connect: %v", err)
	}
	return r
}

func (client *GrpcClient) CreateAccount(username, password string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(client.Timeout))
	defer cancel()
	r, err := client.c.CreateAccount(ctx, &rpc.CreateAccountRequest{Username: username, Password: gcrypto.PasswordToSha256(password)})
	if err != nil {
		log.Printf("could not connect: %v", err)
		return false
	}
	return r.Result
}

func (client *GrpcClient) CreateGenesisEncrypt(id uint32, password string) bool {
	return client.CreateGenesis(id, gcrypto.PasswordToSha256(password))
}

func (client *GrpcClient) CreateGenesis(id uint32, password []byte) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(client.Timeout))
	defer cancel()
	r, err := client.c.CreateGenesis(ctx, &rpc.CreateGenesisRequest{Id: id, Password: password})
	if err != nil {
		log.Printf("could not connect: %v", err)
		return false
	}
	return r.Result
}

func (client *GrpcClient) GetPrivateKey(id uint32, username, password string) ([]byte, bool) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(client.Timeout))
	defer cancel()
	r, err := client.c.GetPrivateKey(ctx, &rpc.PrivateKeyRequest{Id: id, Username: username, Password: gcrypto.PasswordToSha256(password)})
	if err != nil {
		log.Printf("could not connect: %v", err)
		return nil, false
	}
	return r.PrivateKey, r.Result
}

func (client *GrpcClient) ForkContainer(username, password, ip, port string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(client.Timeout))
	defer cancel()
	r, err := client.c.ForkContainer(ctx, &rpc.ForkContainerRequest{
		Username: username, Password: gcrypto.PasswordToSha256(password), Ip: ip, Port: port})
	if err != nil {
		log.Printf("could not connect: %v", err)
		return false
	}
	return r.Result
}

func (client *GrpcClient) CreateContainer(username, password, ip, port string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(client.Timeout))
	defer cancel()
	r, err := client.c.CreateContainer(ctx, &rpc.CreateContainerRequest{
		Username: username, Password: gcrypto.PasswordToSha256(password), Ip: ip, Port: port})
	if err != nil {
		log.Printf("could not connect: %v", err)
		return false
	}
	return r.Result
}

func (client *GrpcClient) LoginContainer(id uint32, password []byte, username, ip, port string) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(client.Timeout))
	defer cancel()
	r, err := client.c.LoginContainer(ctx, &rpc.LoginRequest{
		Id: id, Username: username, Password: password, Ip: ip, Port: port})
	if err != nil {
		log.Printf("could not connect: %v", err)
		return false
	}
	return r.Result

}

func (client *GrpcClient) ControlContainer(id uint32, control rpc.ContainerControlType) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(client.Timeout))
	defer cancel()
	r, err := client.c.ControlContainer(ctx, &rpc.ControlContainerRequest{Id: id, Control: control})
	if err != nil {
		log.Printf("could not connect: %v", err)
		return false
	}
	return r.Result
}

func (client *GrpcClient) ReleaseContainer(id uint32) bool {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(client.Timeout))
	defer cancel()
	r, err := client.c.ReleaseContainer(ctx, &rpc.ReleaseRequest{Id: id})
	if err != nil {
		log.Printf("could not connect: %v", err)
		return false
	}
	return r.Result
}

func (client *GrpcClient) GetLog(id uint32) []byte {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(client.Timeout))
	defer cancel()
	r, err := client.c.GetLog(ctx, &rpc.GetLogRequest{Id: id})
	if err != nil {
		log.Printf("could not connect: %v", err)
	}
	return r.Logbuf
}

func (client *GrpcClient) CheckStatus(id uint32) uint32 {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(client.Timeout))
	defer cancel()
	r, err := client.c.CheckStatus(ctx, &rpc.CoreStateRequest{Id: id})
	if err != nil {
		log.Printf("could not connect: %v", err)
	}
	return r.State
}

func (client *GrpcClient) GetAccount(id uint32) []*ptypes.GhostUser {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(client.Timeout))
	defer cancel()
	r, err := client.c.GetAccount(ctx, &rpc.GetAccountRequest{Id: id})
	if err != nil {
		log.Printf("could not connect: %v", err)
		return nil
	}
	return r.User
}

func (client *GrpcClient) GetBlockInfo(id, blockId uint32) *ptypes.PairedBlocks {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(client.Timeout))
	defer cancel()
	r, err := client.c.GetBlockInfo(ctx, &rpc.GetBlockInfoRequest{Id: id, BlockId: blockId})
	if err != nil {
		log.Printf("could not connect: %v", err)
		return nil
	}
	return r.Pair
}
