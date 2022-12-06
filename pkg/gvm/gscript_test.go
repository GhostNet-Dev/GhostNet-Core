package gvm

import (
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/stretchr/testify/assert"
)

func TestMakeLockScriptOut(t *testing.T) {
	gvm := NewGVM()
	gScript := NewGScript()
	ghostAddr := gcrypto.GenerateKeyPair()
	toAddr := ghostAddr.PubKey
	scriptBuf := gScript.MakeLockScriptOut(toAddr)
	inputParam := gScript.MakeInputParam(scriptBuf, ghostAddr)
	result := gvm.ExecuteScript(scriptBuf, inputParam)
	assert.Equal(t, true, result, "script 실행에 실패했습니다.")
}
