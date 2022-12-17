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
	scriptBuf := gScript.MakeLockScriptOut(ghostAddr.Get160PubKey())
	inputParam := gScript.MakeInputParam(scriptBuf, ghostAddr)
	result := gvm.PushParam(inputParam)
	assert.Equal(t, true, result, "pushparam 실행에 실패했습니다.")
	result = gvm.ExecuteScript(scriptBuf, scriptBuf)
	assert.Equal(t, true, result, "script 실행에 실패했습니다.")
}
