package manager

import (
	"os/exec"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/grpc"
)

type Container struct {
	Id       uint32
	PID      int
	Ip       string
	Port     string
	PubKey   string
	Username string
	Cmd      *exec.Cmd
	Client   *grpc.GrpcClient
}

func NewContainer(id uint32, ip, port string) *Container {
	return &Container{
		Id: id, Ip: ip, Port: port,
		Client: grpc.NewGrpcClient(ip, port),
	}
}

func (c *Container) ConnectServer() {
	c.Client.ConnectServer()
}
