package main

import (
	"models"
	"repository"
)

type MapStorage struct {
	storage map[int]models.Task
}

func NewMapStorageRepository() repository.TaskStorage {
	repository := MapStorage{}
	return &repository
}
