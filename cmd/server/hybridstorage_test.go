package main

import (
	"errors"
	"models"
	"testing"
)

func TestInsert(t *testing.T) {
	storage := NewHybridRepository()

	var i int64
	var cnt int = 10

	for i = 0; i < int64(cnt); i++ {
		task := models.Task{ID: i}
		storage.Add(&task)
	}

	tasks, err := storage.Gets(0, cnt)
	if err != nil {
		t.Error(err)
	}

	if len(tasks) != cnt {
		t.Error("Wrong count elements")
	}

	//t.Logf("%#v",tasks)
}

func TestDelete(t *testing.T) {
	storage := NewHybridRepository()

	var i int64
	var cnt int = 10

	for i = 0; i < int64(cnt); i++ {
		task := models.Task{ID: i}
		storage.Add(&task)
	}

	err := storage.Delete(1)
	if err != nil {
		t.Error(err)
	}
	err = storage.Delete(5)
	if err != nil {
		t.Error(err)
	}

	tasks, err := storage.Gets(0, cnt-2)
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != cnt-2 {
		t.Fatal("Wrong count elements")
	}

	//t.Logf("%#v",tasks)
}

func TestFind(t *testing.T) {
	storage := NewHybridRepository()

	var i int64
	var cnt int = 10

	for i = 0; i < int64(cnt); i++ {
		task := models.Task{ID: i}
		storage.Add(&task)
	}

	task, err := storage.Get(5)
	if err != nil {
		t.Error(err)
	}
	if task.ID != 5 {
		t.Error(
			"expected", 5,
			"got", task.ID,
		)
	}

	task, err = storage.Get(20)
	if err == nil {
		t.Error(
			"expected error",
			"got", task.ID,
		)
	}
}

func TestBoundary(t *testing.T) {
	storage := NewHybridRepository()

	var i int64
	var cnt int = 10

	for i = 0; i < int64(cnt); i++ {
		task := models.Task{ID: i}
		storage.Add(&task)
	}

	tasks, err := storage.Gets(0, 0)
	if err != nil {
		t.Error(err)
	}

	if len(tasks) != 0 {
		t.Error("Wrong count elements")
	}

	tasks, err = storage.Gets(10, 0)
	if err == nil {
		t.Error(errors.New("Error element"))
	}

	tasks, err = storage.Gets(9, 20)
	if err != nil {
		t.Error(errors.New("Wrong count elements"))
	}

}
