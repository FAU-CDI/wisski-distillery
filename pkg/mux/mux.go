// Package mux provides mux
package mux

import (
	"context"
	"net/http"
)

// Mux represents a mux that can handle different requests
type Mux[C any] struct {
	prefixes map[string][]handler
	exacts   map[string][]handler

	Context func(r *http.Request) C // called to set context on the given request

	Panic    func(panic any, w http.ResponseWriter, r *http.Request) // called on panic
	NotFound http.Handler                                            // optional handler to be called in case of a not found
}

type contextKey struct{}

var theContextKey = contextKey{}

type handler struct {
	Predicate Predicate
	http.Handler
}

func (mux *Mux[T]) Prepare(r *http.Request) *http.Request {
	if mux == nil || mux.Context == nil {
		return r
	}

	ctx := context.WithValue(r.Context(), theContextKey, mux.Context(r))
	return r.WithContext(ctx)
}

func (mux *Mux[T]) ContextOf(r *http.Request) (t T) {
	value, ok := r.Context().Value(theContextKey).(T)
	if !ok {
		return t
	}
	return value
}

// Add adds a handler for the given path
func (mux *Mux[T]) Add(path string, predicate Predicate, exact bool, h http.Handler) {
	if mux.exacts == nil {
		mux.exacts = make(map[string][]handler)
	}
	if mux.prefixes == nil {
		mux.prefixes = make(map[string][]handler)
	}

	mPath := NormalizePath(path)
	mHandler := handler{Predicate: predicate, Handler: h}
	if exact {
		mux.exacts[mPath] = append(mux.exacts[mPath], mHandler)
	} else {
		mux.prefixes[mPath] = append(mux.prefixes[mPath], mHandler)
	}
}

// Match returns the handler to be applied for the given request.
func (mux *Mux[T]) Match(r *http.Request, prepare bool) (http.Handler, bool) {
	if mux == nil {
		return nil, false
	}

	if prepare {
		r = mux.Prepare(r)
	}

	candidate := NormalizePath(r.URL.Path)

	// match the exact path first
	for _, h := range mux.exacts[candidate] {
		if h.Predicate.Call(r) {
			return h.Handler, true
		}
	}

	// iterate over path segment candidates
	for {
		// check the current candidate
		for _, h := range mux.prefixes[candidate] {
			if h.Predicate.Call(r) {
				return h.Handler, true
			}
		}

		// if the candidate is the root url, we can bail out now
		if len(candidate) == 0 || candidate == "/" {
			return nil, false
		}

		// move to the parent segment
		candidate = parentSegment(candidate)
	}

}

func (mux *Mux[T]) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// handle panics with the panic handler
	defer func() {
		caught := recover()
		if caught == nil {
			return
		}

		if mux == nil || mux.Panic == nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// silently ignore any panic()s in the panic handler
		defer func() {
			recover()
		}()

		// call the panic handler
		mux.Panic(caught, w, r)
	}()

	// prepare the request
	r = mux.Prepare(r)

	// find the right handler
	// or go into 404 mode
	handler, ok := mux.Match(r, false)
	if !ok {
		if mux == nil || mux.NotFound == nil {
			http.NotFound(w, r)
			return
		}
		mux.NotFound.ServeHTTP(w, r)
		return
	}

	// call the actual handling
	handler.ServeHTTP(w, r)
}

// Predicate represents a matching predicate for a given request.
// The nil predicate always matches
type Predicate func(r *http.Request) bool

// Call checks if this predicate matches the given request.
func (p Predicate) Call(r *http.Request) bool {
	if p == nil {
		return true
	}
	return p(r)
}
