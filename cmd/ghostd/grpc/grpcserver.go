package deamongrpc

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/glogger"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/rpc"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	rpc.UnimplementedGApiServer
	CreateAccountHandler    func(username, password string) bool
	CreateGenesisHandler    func(password string) bool
	CreateContainerHandler  func(username, password string) bool
	ControlContainerHandler func(id uint32, control rpc.ContainerControlType) bool
	ReleaseContainerHandler func(id uint32) bool
	GetPrivateKeyHandler    func(username, password string) ([]byte, bool)
	GetLogHandler           func() []byte
	CheckStatusHandler      func() uint32
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}

func (s *GrpcServer) GetLog(ctx context.Context, in *rpc.GetLogRequest) (*rpc.GetLogResponse, error) {
	var result []byte
	if s.GetLogHandler != nil {
		result = s.GetLogHandler()
	}
	return &rpc.GetLogResponse{Logbuf: result}, nil
}

func (s *GrpcServer) ReleaseContainer(ctx context.Context, in *rpc.ReleaseRequest) (*rpc.ReleaseResponse, error) {
	result := false
	if s.ReleaseContainerHandler != nil {
		result = s.ReleaseContainerHandler(in.Id)
	}
	return &rpc.ReleaseResponse{Result: result}, nil
}

func (s *GrpcServer) ControlContainer(ctx context.Context, in *rpc.ControlContainerRequest) (*rpc.ControlContainerResponse, error) {
	result := false
	if s.ControlContainerHandler != nil {
		result = s.ControlContainerHandler(in.Id, in.Control)
	}
	return &rpc.ControlContainerResponse{Result: result}, nil
}

func (s *GrpcServer) CreateContainer(ctx context.Context, in *rpc.CreateContainerRequest) (*rpc.CreateContainerResponse, error) {
	result := false
	if s.CreateContainerHandler != nil {
		result = s.CreateContainerHandler(in.Username, in.Password)
	}
	return &rpc.CreateContainerResponse{Result: result}, nil
}

func (s *GrpcServer) GetPrivateKey(ctx context.Context, in *rpc.PrivateKeyRequest) (*rpc.PrivateKeyResponse, error) {
	var result bool
	var privateKey []byte
	if s.GetPrivateKeyHandler != nil {
		privateKey, result = s.GetPrivateKeyHandler(in.Username, in.Password)
	}
	return &rpc.PrivateKeyResponse{Result: result, PrivateKey: privateKey}, nil
}

func (s *GrpcServer) CreateGenesis(ctx context.Context, in *rpc.CreateGenesisRequest) (*rpc.CreateGenesisResponse, error) {
	result := false
	if s.CreateGenesisHandler != nil {
		result = s.CreateGenesisHandler(in.Password)
	}

	return &rpc.CreateGenesisResponse{Result: result}, nil
}

func (s *GrpcServer) CreateAccount(ctx context.Context, in *rpc.CreateAccountRequest) (*rpc.CreateAccountResponse, error) {
	result := false
	if s.CreateAccountHandler != nil {
		result = s.CreateAccountHandler(in.Username, in.Password)
	}
	return &rpc.CreateAccountResponse{Result: result}, nil
}

func (s *GrpcServer) CheckStatus(ctx context.Context, in *rpc.CoreStateRequest) (*rpc.CoreStateResponse, error) {
	status := uint32(0)
	if s.CheckStatusHandler != nil {
		status = s.CheckStatusHandler()
	}
	return &rpc.CoreStateResponse{State: status}, nil
}

func (grpcServer *GrpcServer) ServeGRPC(cfg gconfig.GConfig) error {
	lis, err := net.Listen("tcp", ":"+cfg.GrpcPort)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	rpc.RegisterGApiServer(s, grpcServer)
	glogger.DebugOutput(grpcServer, fmt.Sprint("start gRPC Server on ", cfg.Port), 0)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return nil
}
