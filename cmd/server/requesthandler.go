package main

import (
	"bytes"
	"common"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"math"
	"models"
	"net/http"
	"repository"
	"time"
)

type RequestHandler struct {
	mux         *mux.Router
	taskStorage repository.TaskRepository
	counter     common.CounterInt64
	httpClient  http.Client
}

func NewRequestHandler(ts repository.TaskRepository, timeout time.Duration) http.Handler {
	mux := mux.NewRouter()

	s := &RequestHandler{mux: mux, taskStorage: ts, counter: common.NewCounterInt64(),
		httpClient: http.Client{Timeout: timeout}}

	//REST API
	mux.HandleFunc("/task", s.taskCreate).Methods("POST", "PUT")
	mux.HandleFunc("/task", s.tasksList).Methods("GET")
	mux.HandleFunc("/task/{id}", s.taskGet).Methods("GET")
	mux.HandleFunc("/task/{id}", s.taskDelete).Methods("DELETE")

	return s
}

func (s *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *RequestHandler) taskCreate(w http.ResponseWriter, r *http.Request) {
	log.Printf("taskCreate")

	var task models.Task

	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err = task.ValidateCreateRequest(); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request, err := http.NewRequest(task.Method, task.URL, bytes.NewBuffer([]byte(task.RequestBody)))
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	request.Header=task.Header

	resp, err := s.httpClient.Do(request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	task.StatusCode = resp.StatusCode
	//log.Println(string(body))
	task.Length = len(body)
	task.ID = s.counter.NextValue()

	w.Header().Set("Content-Type", "application/json")

	err = json.NewEncoder(w).Encode(task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	s.taskStorage.Add(&task)
}

func (s *RequestHandler) taskGet(w http.ResponseWriter, r *http.Request) {
	log.Printf("taskGet")
}

func (s *RequestHandler) tasksList(w http.ResponseWriter, r *http.Request) {
	log.Printf("tasksList")

	params := r.URL.Query()

	var err error
	var pageNumber, pageSize int = 0, math.MaxUint32

	val := params.Get("pagenumber")
	if val != "" {
		pageNumber, err = common.ConvertStringToInt(val)
		if err != nil {
			http.Error(w, "Wrong pagenumber params", http.StatusBadRequest)
			return
		}
	}

	val = params.Get("pagesize")
	if val != "" {
		pageSize, err = common.ConvertStringToInt(val)
		if err != nil {
			http.Error(w, "Wrong pagesize params", http.StatusBadRequest)
			return
		}
	}

	tasks, err := s.taskStorage.GetPage(pageNumber, pageSize)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(tasks)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}

func (s *RequestHandler) taskDelete(w http.ResponseWriter, r *http.Request) {
	log.Printf("taskDelete")

	vars := mux.Vars(r)

	taskID, err := common.ConvertStringToInt64(vars["id"])
	if err != nil {
		http.Error(w, "Wrong ID params", http.StatusBadRequest)
		return
	}

	err = s.taskStorage.Delete(taskID)
	if err!=nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
