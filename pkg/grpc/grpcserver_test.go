package grpc

import (
	"sync"
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/internal/gconfig"
	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/rpc"
	"github.com/stretchr/testify/assert"
)

var (
	server = NewGrpcServer()
	client = NewGrpcClient("localhost", "50229", 3)
	cfg    = gconfig.NewDefaultConfig()
)

func PreCondition() {
	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		wg.Done()
		server.ServeGRPC(cfg.GrpcPort)
	}()
	wg.Wait()

	client.ConnectServer()
}

func TestGrpcGetInfo(t *testing.T) {
	server.GetInfoHandler = func() *rpc.GetInfoResponse {
		return &rpc.GetInfoResponse{TotalContainer: 1}
	}
	PreCondition()
	defer client.CloseServer()
	response := client.GetInfo()
	assert.Equal(t, uint32(1), response.TotalContainer, "not expected container count")
}

func TestGetLog(t *testing.T) {
	server.GetLogHandler = func(id uint32) []byte {
		return []byte("test")
	}
	PreCondition()
	defer client.CloseServer()
	response := client.GetLog(0)
	assert.Equal(t, []byte("test"), response, "not expected container count")
}

func TestCreateContainer(t *testing.T) {
	server.CreateContainerHandler = func(passwdHash []byte, username, ip, port string) bool {
		return true
	}
	PreCondition()
	defer client.CloseServer()
	response := client.CreateContainer("test", "passwd", "localhost", "50229")
	assert.Equal(t, true, response, "not expected container count")
}
