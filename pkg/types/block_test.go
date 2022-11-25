package types

import (
	"fmt"
	"testing"

	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	pb "github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/types"

	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

func TestGetHashKey(t *testing.T) {
	ghostNetBlock := GhostNetBlock{}
	hash := ghostNetBlock.GetHashKey()
	//fmt.Printf("%02x", hash)
	size := uint32(len(hash))
	assert.Equal(t, size, ghostBytes.HashSize, "Size가 다릅니다.")
}

func TestProtoBlock(t *testing.T) {
	ghostNetBlock := &pb.GhostNetBlock{
		Header: &pb.GhostNetBlockHeader{
			Id:               1,
			AliceCount:       1212,
			TransactionCount: 3232,
		},
	}

	fmt.Printf("%02x\n", ghostNetBlock.Header.Id)

	out, err := proto.Marshal(ghostNetBlock)
	if err != nil {
		assert.Fail(t, "Failed to encode ", err)
	}
	fmt.Println(len(out))

	newGhostNetBlock := &pb.GhostNetBlock{}
	if err := proto.Unmarshal(out, newGhostNetBlock); err != nil {
		assert.Fail(t, "Failed to decode ", err)
	}
	fmt.Printf("%02x\n", newGhostNetBlock.Header.Id)
	/*
		hash := ghostNetBlock.GetHashKey()
		//fmt.Printf("%02x", hash)
		size := uint32(len(hash))
		assert.Equal(t, size, ghostBytes.HashSize, "Size가 다릅니다.")
	*/
}
