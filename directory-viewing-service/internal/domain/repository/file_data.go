package repository

import (
	"context"
	"directory-viewing-service/internal/domain/models"
)

type FileDataRepository interface {
	Save(ctx context.Context, fd []*models.FileData) error
	RecordsByUID(ctx context.Context, uid string, limit int) ([]*models.FileData, error)
}
