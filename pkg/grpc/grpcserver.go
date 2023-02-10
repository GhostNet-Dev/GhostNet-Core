package grpc

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
	CreateAccountHandler    func(passwdHash []byte, username string) bool
	CreateGenesisHandler    func(id uint32, passwdHash []byte) bool
	CreateContainerHandler  func(passwdHash []byte, username, ip, port string) bool
	ControlContainerHandler func(id uint32, control rpc.ContainerControlType) bool
	ReleaseContainerHandler func(id uint32) bool
	GetPrivateKeyHandler    func(id uint32, passwdHash []byte, username string) ([]byte, bool)
	GetLogHandler           func(id uint32) []byte
	GetContainerListHandler func(id uint32) *rpc.GetContainerListResponse
	CheckStatusHandler      func(id uint32) uint32
	GetInfoHandler          func() *rpc.GetInfoResponse
}

func NewGrpcServer() *GrpcServer {
	return &GrpcServer{}
}

func (s *GrpcServer) GetInfo(ctx context.Context, in *rpc.GetInfoRequest) (*rpc.GetInfoResponse, error) {
	var response *rpc.GetInfoResponse = nil
	if s.GetInfoHandler != nil {
		response = s.GetInfoHandler()
	} else {
		response = &rpc.GetInfoResponse{}
	}
	return response, nil
}

func (s *GrpcServer) GetContainerList(ctx context.Context, in *rpc.GetContainerListRequest) (*rpc.GetContainerListResponse, error) {
	var response *rpc.GetContainerListResponse = nil
	if s.GetContainerListHandler != nil {
		response = s.GetContainerListHandler(in.Id)
	} else {
		response = &rpc.GetContainerListResponse{Id: in.Id}
	}
	return response, nil
}

func (s *GrpcServer) GetLog(ctx context.Context, in *rpc.GetLogRequest) (*rpc.GetLogResponse, error) {
	var result []byte
	if s.GetLogHandler != nil {
		result = s.GetLogHandler(in.Id)
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
		result = s.CreateContainerHandler(in.Password, in.Username, in.Ip, in.Port)
	}
	return &rpc.CreateContainerResponse{Result: result}, nil
}

func (s *GrpcServer) GetPrivateKey(ctx context.Context, in *rpc.PrivateKeyRequest) (*rpc.PrivateKeyResponse, error) {
	var result bool
	var privateKey []byte
	if s.GetPrivateKeyHandler != nil {
		privateKey, result = s.GetPrivateKeyHandler(in.Id, in.Password, in.Username)
	}
	return &rpc.PrivateKeyResponse{Result: result, PrivateKey: privateKey}, nil
}

func (s *GrpcServer) CreateGenesis(ctx context.Context, in *rpc.CreateGenesisRequest) (*rpc.CreateGenesisResponse, error) {
	result := false
	if s.CreateGenesisHandler != nil {
		result = s.CreateGenesisHandler(in.Id, in.Password)
	}

	return &rpc.CreateGenesisResponse{Result: result}, nil
}

func (s *GrpcServer) CreateAccount(ctx context.Context, in *rpc.CreateAccountRequest) (*rpc.CreateAccountResponse, error) {
	result := false
	if s.CreateAccountHandler != nil {
		result = s.CreateAccountHandler(in.Password, in.Username)
	}
	return &rpc.CreateAccountResponse{Result: result}, nil
}

func (s *GrpcServer) CheckStatus(ctx context.Context, in *rpc.CoreStateRequest) (*rpc.CoreStateResponse, error) {
	status := uint32(0)
	if s.CheckStatusHandler != nil {
		status = s.CheckStatusHandler(in.Id)
	}
	return &rpc.CoreStateResponse{State: status}, nil
}

func (grpcServer *GrpcServer) ServeGRPC(cfg *gconfig.GConfig) error {
	lis, err := net.Listen("tcp", ":"+cfg.GrpcPort)
	if err != nil {
		return err
	}

	s := grpc.NewServer()
	rpc.RegisterGApiServer(s, grpcServer)
	glogger.DebugOutput(grpcServer, fmt.Sprint("start gRPC Server on ", cfg.GrpcPort, "\n"), 0)
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return nil
}
