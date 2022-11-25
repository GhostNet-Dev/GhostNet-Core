package types

import (
	"fmt"
	"testing"

	//ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	//"github.com/stretchr/testify/assert"
	pb "github.com/GhostNet-Dev/GhostNet-Core/pkg/proto/types"
)

func TestGetHashKey(t *testing.T) {
	ghostNetBlock := pb.GhostNetBlock{
		Id: 1,
	}
	fmt.Printf("%02x", ghostNetBlock.Id)
	/*
		hash := ghostNetBlock.GetHashKey()
		//fmt.Printf("%02x", hash)
		size := uint32(len(hash))
		assert.Equal(t, size, ghostBytes.HashSize, "Size가 다릅니다.")
	*/
}
