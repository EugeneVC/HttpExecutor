package main

import (
	"bytes"
	"common"
	"encoding/json"
	"io/ioutil"
	"log"
	"math"
	"models"
	"net/http"
	"repository"
	"strconv"
	"time"
)

type RequestHandler struct {
	mux         *http.ServeMux
	taskStorage repository.TaskStorage
	counter     common.CounterInt64
	httpClient  http.Client
}

func NewRequestHandler(ts repository.TaskStorage, timeout time.Duration) http.Handler {
	mux := http.NewServeMux()

	s := &RequestHandler{mux: mux, taskStorage: ts, counter: common.NewCounterInt64(), httpClient: http.Client{Timeout: timeout}}

	//REST
	mux.HandleFunc("/task", s.tasksHandler)

	return s
}

func (s *RequestHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}

func (s *RequestHandler) tasksHandler(w http.ResponseWriter, r *http.Request) {
	var task models.Task

	switch r.Method {
	case "POST", "PUT":
		s.taskCreate(w, r, &task)
	case "GET":
		s.tasksList(w, r)
	default:
		http.Error(w, "Only POST,PUT,GET allowed", http.StatusMethodNotAllowed)
	}

}

func (s *RequestHandler) taskCreate(w http.ResponseWriter, r *http.Request, task *models.Task) {
	log.Printf("taskCreate")

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

	s.taskStorage.Add(task)
}

func (s *RequestHandler) tasksList(w http.ResponseWriter, r *http.Request) {
	log.Printf("tasksList")

	params := r.URL.Query()

	var err error
	var offset,limit int = 0,math.MaxInt32

	val := params.Get("offset")
	if val!=""{
		offset,err = getHttpIntParam(val)
		if err!=nil{
			http.Error(w, "Wrong offset params", http.StatusBadRequest)
			return
		}
	}

	val = params.Get("limit")
	if val!=""{
		limit,err = getHttpIntParam(val)
		if err!=nil{
			http.Error(w, "Wrong limit params", http.StatusBadRequest)
			return
		}
	}

	tasks, err := s.taskStorage.Gets(offset, limit)
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

func getHttpIntParam(param string) (int,error){
	if param==""{
		return 0,nil
	}

	val,err := strconv.Atoi(param)
	if err!=nil {
		return 0,err
	}

	return val,nil
}

