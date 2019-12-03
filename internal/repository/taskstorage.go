package repository

import "models"

type TaskStorage interface {
	Add(task *models.Task)
	Delete(key int64) error
	Get(key int64) (*models.Task,error)
	Gets(offset,limit int) ([]*models.Task, error)
}
