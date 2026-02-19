package services

import (
	"context"
	"directory-viewing-service/internal/domain/models"
)

type FileDataService interface {
	UploadParsed(ctx context.Context, fd []*models.FileData) error
	ProcessedData(ctx context.Context, uid string, limit int) ([]*models.FileData, error)
}
