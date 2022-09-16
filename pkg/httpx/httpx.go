package httpx

import "errors"

// ErrNotFound should be returned from any httpx error to indicate that the item was not found
var ErrNotFound = errors.New("httpx: Error 404")
