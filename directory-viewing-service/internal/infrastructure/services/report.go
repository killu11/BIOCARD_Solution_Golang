package services

import (
	"context"
	"directory-viewing-service/internal/domain/models"
	"directory-viewing-service/internal/domain/repository"
)

type ReportService struct {
	repo repository.ReportRepository
}

func NewReportService(repo repository.ReportRepository) *ReportService {
	return &ReportService{repo: repo}
}

func (r *ReportService) ReportError(ctx context.Context, filename, msg string) error {
	return r.repo.Save(ctx, models.NewReport(filename, msg))
}
