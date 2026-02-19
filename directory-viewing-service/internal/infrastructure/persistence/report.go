package persistence

import (
	"context"
	"directory-viewing-service/internal/domain/models"
	"directory-viewing-service/pkg"

	"github.com/jackc/pgx/v5/pgxpool"
)

type ReportRepository struct {
	db *pgxpool.Pool
}

func NewReportRepository(db *pgxpool.Pool) *ReportRepository {
	return &ReportRepository{db: db}
}

func (r *ReportRepository) Save(ctx context.Context, task *models.Report) error {
	query := `INSERT INTO reports(filename, msg) 
	VALUES ($1, $2)`
	_, err := r.db.Exec(ctx, query, task.Filename, task.Msg)
	if err != nil {
		return pkg.PackageError(packageName, "save report", err)
	}
	return nil
}
