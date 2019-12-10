package main

import (
	"errors"
	"math"
	"models"
	"testing"
)

func TestInsert(t *testing.T) {
	hybridTaskRepository := NewHybridTaskRepository()

	var i int64
	var cnt int = 10

	for i = 0; i < int64(cnt); i++ {
		task := models.Task{ID: i}
		hybridTaskRepository.Add(&task)
	}

	tasks, err := hybridTaskRepository.GetPage(0, math.MaxInt32)
	if err != nil {
		t.Error(err)
	}

	if len(tasks) != cnt {
		t.Error("Wrong count elements")
	}

	tasks, err = hybridTaskRepository.GetPage(0, cnt)
	if err != nil {
		t.Error(err)
	}

	if len(tasks) != cnt {
		t.Error("Wrong count elements")
	}

	//t.Logf("%#v",tasks)
}

func TestDelete(t *testing.T) {
	hybridTaskRepository := NewHybridTaskRepository()

	var i int64
	var cnt int = 10

	for i = 0; i < int64(cnt); i++ {
		task := models.Task{ID: i}
		hybridTaskRepository.Add(&task)
	}

	err := hybridTaskRepository.Delete(1)
	if err != nil {
		t.Error(err)
	}

	err = hybridTaskRepository.Delete(1)
	if err == nil {
		t.Error("Re-delete error")
	}

	err = hybridTaskRepository.Delete(5)
	if err != nil {
		t.Error(err)
	}

	tasks, err := hybridTaskRepository.GetPage(0, math.MaxInt32)
	if err != nil {
		t.Fatal(err)
	}

	if len(tasks) != cnt-2 {
		t.Fatal("Wrong count elements")
	}

	//t.Logf("%#v",tasks)
}

func TestFind(t *testing.T) {
	hybridTaskRepository := NewHybridTaskRepository()

	var i int64
	var cnt int = 10

	for i = 0; i < int64(cnt); i++ {
		task := models.Task{ID: i}
		hybridTaskRepository.Add(&task)
	}

	task, err := hybridTaskRepository.Get(5)
	if err != nil {
		t.Error(err)
	}
	if task.ID != 5 {
		t.Error(
			"expected", 5,
			"got", task.ID,
		)
	}

	task, err = hybridTaskRepository.Get(20)
	if err == nil {
		t.Error(
			"expected error",
			"got", task.ID,
		)
	}
}

func TestBoundary(t *testing.T) {
	hybridTaskRepository := NewHybridTaskRepository()

	var i int64
	var cnt int = 10

	for i = 0; i < int64(cnt); i++ {
		task := models.Task{ID: i}
		hybridTaskRepository.Add(&task)
	}

	tasks, err := hybridTaskRepository.GetPage(0, 0)
	if err == nil {
		t.Error(err)
	}

	if len(tasks) != 0 {
		t.Error("Wrong count elements")
	}

	tasks, err = hybridTaskRepository.GetPage(9, 1)
	if err != nil || len(tasks) != 1 || tasks[0].ID!=9 {
		t.Error(errors.New("Error element"))
	}

	tasks, err = hybridTaskRepository.GetPage(9, 20)
	if err != nil {
		t.Error(errors.New("Wrong count elements"))
	}

}
