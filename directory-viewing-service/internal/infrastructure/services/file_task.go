package services

import (
	"context"
	"directory-viewing-service/internal/domain/models"
	"directory-viewing-service/internal/domain/repository"
	"directory-viewing-service/internal/domain/services"
)

type FileTaskService struct {
	repo repository.FileTaskRepository
}

func NewFileTaskService(r repository.FileTaskRepository) *FileTaskService {
	return &FileTaskService{repo: r}
}

func (f *FileTaskService) UploadTasks(ctx context.Context, tsvNames []string) ([]*models.FileTask, error) {
	unProcFiles, err := f.repo.GetUnProcessedFiles(ctx, tsvNames)
	if err != nil {
		return nil, err
	}

	if len(unProcFiles) < 1 {
		return nil, services.ErrZeroUnProcessedFiles
	}
	tasks, err := f.repo.CreateTasks(ctx, unProcFiles)
	if err != nil {
		return nil, err
	}

	if len(tasks) == 0 {
		return nil, services.ErrCreatedZeroTasks
	}
	return tasks, nil
}

func (f *FileTaskService) ChangeStatus(ctx context.Context, id int, status string) error {
	return f.repo.UpdateStatusByID(ctx, id, status)
}
