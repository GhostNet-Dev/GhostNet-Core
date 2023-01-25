package glogger

import "fmt"

type LogLevel int32

const (
	Default         LogLevel = 0
	BlockConsensus  LogLevel = 1
	MasterListSync  LogLevel = 2
	ForwardingTrace LogLevel = 3
)

func DebugOutput(obj interface{}, msg string, level LogLevel) {
	fmt.Printf(msg)
}
