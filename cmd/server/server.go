package main

import (
	"context"
	"log"
	"net/http"
)

type ServerHTTP struct {
	server http.Server
}

func NewHttpServer(addr string,handler http.Handler) *ServerHTTP{
	serverHTTP := ServerHTTP{server:http.Server{Addr: addr, Handler: handler}}

	return &serverHTTP
}

func (s *ServerHTTP) Start() {
	go func() {
		log.Printf("Listening on %s\n", s.server.Addr)

		if err := s.server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()
}

func (s *ServerHTTP) Stop(){
	s.server.Shutdown(context.Background())
}

