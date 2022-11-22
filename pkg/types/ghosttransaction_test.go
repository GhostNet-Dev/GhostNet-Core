package types

import (
	"fmt"
	"testing"

	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	"github.com/stretchr/testify/assert"
)

func TestOutPointGetSize(t *testing.T) {
	txOutPoint := TxOutPoint{}
	size := txOutPoint.Size()
	fmt.Println(size)
	assert.Equal(t, size, ghostBytes.HashSize+4, "Size가 다릅니다.")
}

func TestTxInputGetSize(t *testing.T) {
	input := TxInput{}
	size := input.Size()
	fmt.Println(size)
	if size < 44 {
		t.Errorf("size: %d < 44", size)
	}
}
