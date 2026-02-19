package repository

import (
	"context"
	"directory-viewing-service/internal/domain/models"
)

type ReportRepository interface {
	Save(ctx context.Context, task *models.Report) error
}
