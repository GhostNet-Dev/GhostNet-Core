package gvm

import (
	"fmt"

	"github.com/GhostNet-Dev/gscript/evaluator"
	"github.com/GhostNet-Dev/gscript/lexer"
	"github.com/GhostNet-Dev/gscript/object"
	"github.com/GhostNet-Dev/gscript/parser"
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
	env := object.NewEnvironment(nil)
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
	case *object.String:
		result = obj.Value
	case string:
		result = obj
	case int64:
		result = fmt.Sprint(obj)
	}

	return result
}
