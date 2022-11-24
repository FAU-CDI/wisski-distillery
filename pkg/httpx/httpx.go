package httpx

import (
	"net/http"
)

// Response represents a response to an http request.
type Response struct {
	ContentType string // defaults to text/plain
	Body        []byte
	StatusCode  int // defaults to [http.StatusOK]
}

func (response Response) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	if response.ContentType == "" {
		response.ContentType = "text/plain"
	}
	w.Header().Set("Content-Type", response.ContentType)

	if response.StatusCode <= 0 {
		response.StatusCode = http.StatusOK
	}
	w.WriteHeader(response.StatusCode)
	w.Write(response.Body)
}
