package repository

import (
	"context"
	"directory-viewing-service/internal/domain/models"
)

type FileTaskRepository interface {
	CreateTasks(ctx context.Context, fn []string) ([]*models.FileTask, error)
	GetUnProcessedFiles(ctx context.Context, fn []string) ([]string, error)
	UpdateStatusByID(ctx context.Context, id int, status string) error
}
