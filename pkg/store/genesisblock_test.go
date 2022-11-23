package store

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestGenesisBlockLoad(t *testing.T) {
	block := GenesisBlock()
	if b, err := json.Marshal(block.Header); err == nil {
		fmt.Printf(string(b))
	}
}
