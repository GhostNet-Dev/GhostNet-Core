package manager

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"

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
	var grpcPort string
	if portInt, err := strconv.Atoi(port); err != nil {
		log.Fatal(err)
	} else {
		grpcPort = fmt.Sprint(portInt + 1)
	}

	return &Container{
		Id: id, Ip: ip, Port: port,
		Client: grpc.NewGrpcClient(ip, grpcPort),
	}
}

func (c *Container) ConnectServer() {
	c.Client.ConnectServer()
}
