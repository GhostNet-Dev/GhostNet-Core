package blocks

import (
	"math/rand"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestConsensusCalc(t *testing.T) {
	now := uint64(time.Now().Unix())
	prev := now - uint64(rand.Intn(60))
	nextTx := CoreCalculator(now, prev, 10)
	assert.Equal(t, true, nextTx >= 10, "Calculated below expected value:", nextTx)
}
