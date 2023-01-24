package pools

import (
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
