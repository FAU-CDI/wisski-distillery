package logging

import (
	"io"
	"sync"
)

type writerIndent struct{}

var indentKey = writerIndent{}

func getIndent(writer io.Writer) int {
	value, ok := getKey(writer, indentKey)
	if !ok {
		value = 0
	}
	return value.(int)
}

func incIndent(writer io.Writer) int {
	value, ok := upsetKey(writer, indentKey, func(value any, fresh bool) any {
		if fresh {
			return 0
		}
		return value.(int) + 1
	})
	if !ok {
		return 0
	}
	return value.(int)
}

func decIndent(writer io.Writer) int {
	value, ok := upsetKey(writer, indentKey, func(value any, fresh bool) any {
		if fresh {
			return 0
		}
		level := value.(int) - 1
		if level < 0 {
			level = 0
		}
		return level
	})
	if !ok {
		return 0
	}
	return value.(int)
}

// KEY-VALUE STORE for writers

var writerDataMutex sync.RWMutex
var writerDataData = make(map[uintptr]map[any]any)

func getKey(writer io.Writer, key any) (value any, ok bool) {
	uid, ok := id(writer)
	if !ok {
		return nil, false
	}

	writerDataMutex.RLock()
	defer writerDataMutex.RUnlock()

	value, ok = writerDataData[uid][key]
	return
}

func setKey(writer io.Writer, key, value any) bool {
	uid, ok := id(writer)
	if !ok {
		return false
	}

	writerDataMutex.Lock()
	defer writerDataMutex.Unlock()

	values, ok := writerDataData[uid]
	if !ok {
		values = make(map[any]any)
		writerDataData[uid] = values
	}
	values[key] = value
	return true
}

func upsetKey(writer io.Writer, key any, update func(value any, fresh bool) any) (any, bool) {
	uid, ok := id(writer)
	if !ok {
		return nil, false
	}

	writerDataMutex.Lock()
	defer writerDataMutex.Unlock()

	values, ok := writerDataData[uid]
	if !ok {
		values = make(map[any]any)
		writerDataData[uid] = values
	}
	values[key] = update(values[key], !ok)
	return values[key], true
}

func id(writer io.Writer) (uintptr, bool) {
	file, ok := writer.(interface{ Fd() uintptr })
	if !ok {
		return 0, false
	}
	return file.Fd(), true
}
