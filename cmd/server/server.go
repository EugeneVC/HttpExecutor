package main

import (
	"net/http"
)

type Server struct {
	mux *http.ServeMux
}

func NewServer() http.Handler {
	s := &Server{mux: http.NewServeMux()}
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
