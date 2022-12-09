package blocks

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConsensusCalc(t *testing.T) {
	blocks := NewBlocks(blockContainer)

	assert.Equal(t, nil, blocks, "")
}
