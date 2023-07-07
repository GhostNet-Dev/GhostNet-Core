package gvm

import (
	"testing"

	"github.com/GhostNet-Dev/GhostNet-Core/pkg/gcrypto"
	"github.com/stretchr/testify/assert"
)

func TestMakeLockScriptOut(t *testing.T) {
	gvm := NewGVM()
	gCompiler := NewGCompiler()
	ghostAddr := gcrypto.GenerateKeyPair()
	scriptBuf := MakeLockScriptOut(ghostAddr.Get160PubKey())
	inputParam := gCompiler.MakeInputParam(scriptBuf, ghostAddr)
	result := gvm.PushParam(inputParam)
	assert.Equal(t, true, result, "pushparam 실행에 실패했습니다.")
	result = gvm.ExecuteScript(scriptBuf, scriptBuf)
	assert.Equal(t, true, result, "script 실행에 실패했습니다.")
}

func TestRootFsTx(t *testing.T) {
	gvm := NewGVM()
	gCompiler := NewGCompiler()
	ghostAddr := gcrypto.GenerateKeyPair()
	scriptBuf := MakeRootAccount(ghostAddr.Get160PubKey(), "test")
	inputParam := gCompiler.MakeInputParam(scriptBuf, ghostAddr)
	result := gvm.PushParam(inputParam)
	assert.Equal(t, true, result, "pushparam 실행에 실패했습니다.")
	result = gvm.ExecuteScript(scriptBuf, scriptBuf)
	assert.Equal(t, true, result, "MakeRootAccount script 실행에 실패했습니다.")
}
