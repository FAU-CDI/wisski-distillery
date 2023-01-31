package httpx

import (
	"net/http"
	"sync"
)

// SyncedResponseWriter wraps a http ResponseWriter to syncronize all actions
type SyncedResponseWriter struct {
	m sync.Mutex
	http.ResponseWriter
}

func (rw *SyncedResponseWriter) Header() http.Header {
	rw.m.Lock()
	defer rw.m.Unlock()

	return rw.ResponseWriter.Header()
}

func (rw *SyncedResponseWriter) Write(data []byte) (int, error) {
	rw.m.Lock()
	defer rw.m.Unlock()

	return rw.ResponseWriter.Write(data)
}

func (rw *SyncedResponseWriter) WriteHeader(statusCode int) {
	rw.m.Lock()
	defer rw.m.Unlock()

	rw.ResponseWriter.WriteHeader(statusCode)
}

// Flush flushes any partial output to the underlying ResponseWriter.
// If the wrapped ResponseWriter does not implement flush, the function performs no operation.
func (rw *SyncedResponseWriter) Flush() {
	f, ok := rw.ResponseWriter.(http.Flusher)
	if !ok {
		return
	}

	rw.m.Lock()
	defer rw.m.Unlock()

	f.Flush()
}
