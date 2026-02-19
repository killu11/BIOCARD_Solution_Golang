package services

import (
	"context"
	"directory-viewing-service/internal/domain/models"
	"errors"
)

const (
	StatusNotProc    = "not processed"
	StatusProcessing = "processing"
	StatusFailed     = "failed"
)

var (
	ErrZeroUnProcessedFiles = errors.New("zero unprocessed files")
	ErrCreatedZeroTasks     = errors.New("created zero tasks")
)

type FileTaskService interface {
	UploadTasks(ctx context.Context, tsvNames []string) ([]*models.FileTask, error)

	ChangeStatus(ctx context.Context, id int, status string) error
}
