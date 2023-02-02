package deamongrpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/rpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient struct {
	GrpcPort string
	c        rpc.GApiClient
	conn     *grpc.ClientConn
}

func NewGrpcClient(grpcPort string) *GrpcClient {
	return &GrpcClient{
		GrpcPort: grpcPort,
	}
}

func (client *GrpcClient) ConnectServer() {
	conn, err := grpc.Dial("localhost:"+client.GrpcPort, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	client.c = rpc.NewGApiClient(conn)
	client.conn = conn
}

func (client *GrpcClient) CloseServer() {
	client.conn.Close()
}

func (client *GrpcClient) GetStateServer() {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.c.CheckStatus(ctx, &rpc.CoreStateRequest{})
	if err != nil {
		log.Fatalf("could not connect: %v", err)
	}
	fmt.Printf("get State = %d", r.State)
}
