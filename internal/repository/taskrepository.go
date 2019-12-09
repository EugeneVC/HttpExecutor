package repository

import "models"

type TaskRepository interface {
	Add(task *models.Task)
	Delete(key int64) error
	Get(key int64) (*models.Task,error)
	GetPage(pageNumber,pageSize int) ([]*models.Task, error)
}
