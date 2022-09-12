package env

import (
	"io"
	"net/http"
)

// TODO: Move this into dis!

// Server represents a server for this distillery
type Server struct {
	dis *Distillery
}

func (dis *Distillery) Server() *Server {
	return &Server{
		dis: dis,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	instances, err := s.dis.AllInstances()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		io.WriteString(w, "Something went wrong")
		return
	}

	w.WriteHeader(http.StatusOK)
	for _, instance := range instances {
		io.WriteString(w, instance.Slug+"\n")
	}
}
