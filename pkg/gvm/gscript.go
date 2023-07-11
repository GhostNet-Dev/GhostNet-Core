package gvm

import (
	"github.com/GhostNet-Dev/glambda/evaluator"
	"github.com/GhostNet-Dev/glambda/lexer"
	"github.com/GhostNet-Dev/glambda/object"
	"github.com/GhostNet-Dev/glambda/parser"
)

type GScript struct {
}

func NewGScript() *GScript {
	gScript := &GScript{}

	return gScript
}

func Eval(code string) interface{} {
	l := lexer.NewLexer(code)
	p := parser.NewParser(l)
	program := p.ParseProgram()
	env := object.NewEnvironment()
	resultObj := evaluator.Eval(program, env)
	switch obj := resultObj.(type) {
	case *object.String:
		return obj.Value
	case *object.Integer:
		return obj.Value
	case *object.Identifier:
		return obj.Value
	}
	return ""
}

func ExecuteScript(code string) (result string) {
	ret := Eval(code)
	responseParam := make(map[string]string)

	switch obj := ret.(type) {
	case *object.Hash:
		for _, hashPair := range obj.Pairs {
			responseParam[hashPair.Key.Inspect()] = hashPair.Value.Inspect()
		}
	}

	return result
}
