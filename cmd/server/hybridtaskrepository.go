package main

import (
	"errors"
	"models"
	"repository"
	"sort"
	"sync"
)

type HybridTaskRepository struct {
	storage map[int64]*models.Task
	keys    []int64 //always sorted
	rw      sync.RWMutex
}

func NewHybridTaskRepository() repository.TaskRepository {
	repository := HybridTaskRepository{storage: map[int64]*models.Task{}, keys: []int64{}}
	return &repository
}

func (hs *HybridTaskRepository) Add(task *models.Task) {
	hs.rw.Lock()
	defer hs.rw.Unlock()

	hs.keys = append(hs.keys, task.ID)
	hs.storage[task.ID] = task
}

func (hs *HybridTaskRepository) Delete(key int64) error {
	hs.rw.Lock()
	defer hs.rw.Unlock()

	idx := sort.Search(len(hs.keys), func(i int) bool { return hs.keys[i] >= key })
	if idx == len(hs.keys) {
		return errors.New("Not found")
	}

	copy(hs.keys[idx:], hs.keys[idx+1:])
	hs.keys = hs.keys[:len(hs.keys)-1]

	delete(hs.storage, key)

	return nil
}

func (hs *HybridTaskRepository) Get(key int64) (*models.Task, error) {
	hs.rw.RLock()
	defer hs.rw.RUnlock()

	task, ok := hs.storage[key]
	if !ok {
		return nil, errors.New("Not found")
	}

	return task, nil
}

func (hs *HybridTaskRepository) Gets(offset, limit int) ([]*models.Task, error) {
	hs.rw.RLock()
	defer hs.rw.RUnlock()

	if offset < 0 || limit < 0 || offset > len(hs.keys) {
		return nil, errors.New("Wrong params")
	}

	if limit == 0 {
		return nil, nil
	}

	if offset+limit >= len(hs.keys) {
		limit = len(hs.keys) - offset
	}

	tasks := []*models.Task{}

	for _, key := range hs.keys[offset : offset+limit] {
		tasks = append(tasks, hs.storage[key])
	}

	return tasks, nil
}
