package main

import (
	"common"
	"encoding/json"
	"fmt"
	"log"
	"models"
	"net/http"
	"repository"
)

type Server struct {
	mux         *http.ServeMux
	taskStorage *repository.TaskStorage
	counter     common.CounterInt64
}

func NewServer(ts *repository.TaskStorage) http.Handler {
	mux := http.NewServeMux()

	s := &Server{mux: mux, taskStorage: ts, counter: common.NewCounter()}

	//REST
	mux.HandleFunc("/tasks", s.TasksHandler)

	return s
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s Server) TasksHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		s.tasksList(w, r)
	case "POST":
		s.taskCreate(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *Server) taskCreate(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	log.Printf("taskCreate")

	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	if err = task.ValidateRequest(); err!=nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}

	task.ID = s.counter.NextValue()



	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
}

func (s *Server) tasksList(w http.ResponseWriter, r *http.Request) {
	log.Printf("tasksList")
}
