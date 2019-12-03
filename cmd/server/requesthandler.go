package main

import (
	"bytes"
	"common"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"models"
	"net/http"
	"repository"
	"time"
)

type RequestHandler struct {
	mux         *http.ServeMux
	taskStorage repository.TaskStorage
	counter     common.CounterInt64
	httpClient http.Client
}

func NewRequestHandler(ts repository.TaskStorage, timeout time.Duration) http.Handler {
	mux := http.NewServeMux()

	s := &RequestHandler{mux: mux, taskStorage: ts, counter: common.NewCounterInt64(),httpClient:http.Client{Timeout:timeout}}

	mux.HandleFunc("/task", s.tasksHandler)

	return s
}

func (s *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *RequestHandler) tasksHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	if r.Method != "POST"{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w,errors.New("Only POST allowed"))
		return
	}

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

	switch task.Method {
		//case "LIST":
		//	s.tasksList(w, r)
		case "POST","GET":
			s.taskCreate(w, r,&task)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func (s *RequestHandler) taskCreate(w http.ResponseWriter, r *http.Request, task *models.Task) {
	log.Printf("taskCreate")

	w.Header().Set("Content-Type", "application/json")

	request,err := http.NewRequest(task.Method,task.URL,bytes.NewBuffer([]byte(task.RequestBody)))
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
	}

	resp, err := s.httpClient.Do(request)
	if err!=nil{
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
	}

	defer resp.Body.Close()

	body,err := ioutil.ReadAll(resp.Body)

	task.StatusCode = resp.StatusCode
	log.Println(string(body))
	task.Length = len(body)
	task.ID = s.counter.NextValue()

	s.taskStorage.Add(task)

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(w, err)
		return
	}
}

func (s *RequestHandler) tasksList(w http.ResponseWriter, r *http.Request) {
	log.Printf("tasksList")
}
