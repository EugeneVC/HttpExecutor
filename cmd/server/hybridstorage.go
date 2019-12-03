package main

import (
	"errors"
	"models"
	"repository"
	"sort"
	"sync"
)

type HybridStorage struct {
	storage map[int64]*models.Task
	keys []int64 //always sorted
	rw	sync.RWMutex
}

func NewHybridRepository() repository.TaskStorage {
	repository := HybridStorage{ storage:make(map[int64]*models.Task,1000), keys: make([]int64,1000) }
	return &repository
}

func (hs *HybridStorage) Add(task *models.Task){
	hs.rw.Lock()
	defer hs.rw.Unlock()

	hs.keys = append(hs.keys, task.ID)
	hs.storage[task.ID] = task
}

func (hs *HybridStorage) Delete(key int64) error{
	hs.rw.Lock()
	defer hs.rw.Unlock()

	idx:=sort.Search(len(hs.keys),func(i int) bool {return hs.keys[i] >= key})
	if idx==len(hs.keys){
		return errors.New("Not found")
	}

	copy(hs.keys[idx:], hs.keys[idx+1:])
	hs.keys = hs.keys[:len(hs.keys)-1]

	delete(hs.storage, key)

	return nil
}

func (hs *HybridStorage) Get(key int64) (*models.Task,error){
	hs.rw.RLock()
	defer hs.rw.RUnlock()

	task,ok := hs.storage[key]
	if !ok {
		return nil,errors.New("Not found")
	}

	return task,nil
}

func (hs *HybridStorage) Gets(offset,limit int) ([]*models.Task, error){
	hs.rw.RLock()
	defer hs.rw.RUnlock()

	if offset+limit>len(hs.keys){
		return nil,errors.New("Wrong params")
	}

	tasks := []*models.Task{}

	for _,key :=range hs.keys[offset:offset+limit]{
		tasks = append(tasks, hs.storage[key])
	}

	return tasks,nil
}





