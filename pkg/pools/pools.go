// Package pools holds various pools for reuse
package pools

import (
	"bytes"
	"strings"
	"sync"
)

var builders = sync.Pool{
	New: func() any { return new(strings.Builder) },
}

func GetBuilder() *strings.Builder {
	return builders.Get().(*strings.Builder)
}

func ReleaseBuilder(builder *strings.Builder) {
	builder.Reset()
	builders.Put(builder)
}

var buffers = sync.Pool{
	New: func() any { return new(bytes.Buffer) },
}

func GetBuffer() *bytes.Buffer {
	return buffers.Get().(*bytes.Buffer)
}

func ReleaseBuffer(buffer *bytes.Buffer) {
	buffer.Reset()
	buffers.Put(buffer)
}
