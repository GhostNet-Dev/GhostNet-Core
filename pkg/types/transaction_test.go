package types

import (
	"fmt"
	"testing"

	ghostBytes "github.com/GhostNet-Dev/GhostNet-Core/libs/bytes"
	"github.com/stretchr/testify/assert"
)

func TestOutPointGetSize(t *testing.T) {
	txOutPoint := TxOutPoint{
		TxId: make([]byte, ghostBytes.HashSize),
	}
	size := txOutPoint.Size()
	fmt.Println(size)
	assert.Equal(t, ghostBytes.HashSize+4, size, "Size is different.")
}

func TestTxInputGetSize(t *testing.T) {
	input := TxInput{}
	size := input.Size()
	fmt.Println(size)
	if size < 44 {
		t.Errorf("size: %d < 44", size)
	}
}
