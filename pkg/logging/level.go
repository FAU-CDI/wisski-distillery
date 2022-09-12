package logging

import (
	"sync"

	"github.com/tkw1536/goprogram/stream"
)

var logLevelMutex sync.Mutex
var logLevelMap = make(map[uintptr]int)

func getIndent(io stream.IOStream) int {
	logLevelMutex.Lock()
	defer logLevelMutex.Unlock()

	id, ok := logID(io)
	if !ok {
		return 0
	}

	return logLevelMap[id]
}

func incIndent(io stream.IOStream) int {
	logLevelMutex.Lock()
	defer logLevelMutex.Unlock()

	id, ok := logID(io)
	if !ok { // if we don't have an id, then inc statically returns 1
		return 1
	}

	logLevelMap[id]++
	return logLevelMap[id]
}

func decIndent(io stream.IOStream) int {
	logLevelMutex.Lock()
	defer logLevelMutex.Unlock()
	id, ok := logID(io)

	if !ok { // if we don't have an id, then dec statically returns 0
		return 0
	}

	logLevelMap[id]--
	if logLevelMap[id] < 0 {
		panic("DecLogIdent: decrease below 0")
	}
	return logLevelMap[id]
}

func logID(io stream.IOStream) (uintptr, bool) {
	file, ok := io.Stdin.(interface{ Fd() uintptr })
	if !ok {
		return 0, false
	}
	return file.Fd(), true
}
