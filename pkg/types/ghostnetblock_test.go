package types

import (
	//"fmt"
	"testing"

	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	"github.com/stretchr/testify/assert"
)

func TestGetHashKey(t *testing.T) {
	ghostNetBlock := GhostNetBlock{}
	hash := ghostNetBlock.GetHashKey()
	//fmt.Printf("%02x", hash)
	size := uint32(len(hash))
	assert.Equal(t, size, ghostBytes.HashSize, "Size가 다릅니다.")
}
