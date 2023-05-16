package glogger

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"runtime"
	"time"
)

type LogLevel int32

type GLogger struct {
	id  uint32
	out io.Writer // destination for output
	buf []byte    // for accumulating text to write
}

const (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

const (
	Default         LogLevel = 0
	BlockConsensus  LogLevel = 1
	MasterListSync  LogLevel = 2
	ForwardingTrace LogLevel = 3
	PacketLog       LogLevel = 4
)

var globalGlog = NewGLogger(0)

func NewGLogger(id uint32) *GLogger {
	return &GLogger{
		id:  id,
		out: os.Stderr,
	}
}

func GetType(myvar interface{}) string {
	if myvar == nil {
		return "{nil}"
	}
	if t := reflect.TypeOf(myvar); t.Kind() == reflect.Ptr {
		return "*" + t.Elem().Name()
	} else {
		return t.Name()
	}
}

func GlobalDebugOutput(obj interface{}, s string, level LogLevel) error {
	return globalGlog.Output(obj, s, level, 2)
}

func (log *GLogger) ShowMetaData(i interface{}) {
	elements := reflect.ValueOf(i).Elem()

	for index := 0; index < elements.NumField(); index++ {
		typeField := elements.Type().Field(index)
		fmt.Println(typeField.Name, typeField.Type, typeField.Tag, elements.Field(index))
	}
}

func (l *GLogger) Output(obj interface{}, s string, level LogLevel, depth int) error {
	now := time.Now() // get this early.
	_, file, line, ok := runtime.Caller(depth)
	if !ok {
		file = "???"
		line = 0
	}
	l.buf = l.buf[:0]
	formatTimeHeader(&l.buf, now, file, line)
	switch level {
	case PacketLog:
		l.buf = append(l.buf, Purple...)
	case BlockConsensus:
		l.buf = append(l.buf, Blue...)
	}
	l.buf = append(l.buf, s...)
	l.buf = append(l.buf, Reset...)

	if len(s) == 0 || s[len(s)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}
	_, err := l.out.Write(l.buf)
	return err
}

func (l *GLogger) DebugOutput(obj interface{}, s string, level LogLevel) error {
	return l.Output(obj, s, level, 2)
}

func formatTimeHeader(buf *[]byte, t time.Time, file string, line int) {
	//date
	year, month, day := t.Date()
	itoa(buf, year, 4)
	*buf = append(*buf, '/')
	itoa(buf, int(month), 2)
	*buf = append(*buf, '/')
	itoa(buf, day, 2)
	*buf = append(*buf, ' ')
	//time
	hour, min, _ /*sec*/ := t.Clock()
	itoa(buf, hour, 2)
	*buf = append(*buf, ':')
	itoa(buf, min, 2)
	/*
		*buf = append(*buf, ':')
		itoa(buf, sec, 2)
	*/
	*buf = append(*buf, ' ')
	//file line
	*buf = append(*buf, Green...)
	short := file
	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' {
			short = file[i+1:]
			break
		}
	}
	file = short
	*buf = append(*buf, file...)
	*buf = append(*buf, ':')
	*buf = append(*buf, Yellow...)
	itoa(buf, line, -1)
	*buf = append(*buf, ": "...)
	*buf = append(*buf, Reset...)
}

func itoa(buf *[]byte, i int, wid int) {
	// Assemble decimal in reverse order.
	var b [20]byte
	bp := len(b) - 1
	for i >= 10 || wid > 1 {
		wid--
		q := i / 10
		b[bp] = byte('0' + i - q*10)
		bp--
		i = q
	}
	// i < 10
	b[bp] = byte('0' + i)
	*buf = append(*buf, b[bp:]...)
}
