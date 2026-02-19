package persistence

import (
	"context"
	"directory-viewing-service/internal/domain/models"
	"directory-viewing-service/pkg"
	"errors"

	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	StatusNotProc    = "not processed"
	StatusProcessing = "processing"
	StatusFailed     = "failed"
)

type FileTaskRepository struct {
	db *pgxpool.Pool
}

func NewFileTaskRepository(db *pgxpool.Pool) *FileTaskRepository {
	return &FileTaskRepository{db: db}
}

func (d *FileTaskRepository) GetUnProcessedFiles(ctx context.Context, fn []string) ([]string, error) {
	query := `SELECT files_names.name 
	FROM unnest($1::TEXT[]) AS files_names(name)
	LEFT JOIN files_tasks ON files_names.name = files_tasks.filename
	WHERE files_tasks.filename IS NULL 
   	OR files_tasks.status = 'not processed'`

	rows, err := d.db.Query(
		ctx,
		query,
		fn,
	)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		return nil, pkg.PackageError(packageName, "query unprocessed files", err)
	}

	defer rows.Close()

	var result []string
	for rows.Next() {
		var name string
		if err = rows.Scan(&name); err != nil {
			return nil, pkg.PackageError(packageName, "scan unprocessed filenames", err)
		}
		result = append(result, name)
	}

	if err = rows.Err(); err != nil {
		return nil, pkg.PackageError(packageName, "scan unprocessed filenames", err)
	}

	return result, nil
}

func (d *FileTaskRepository) CreateTasks(ctx context.Context, fn []string) ([]*models.FileTask, error) {
	query := `INSERT INTO files_tasks (filename)
	SELECT val.filename
	FROM unnest($1::TEXT[]) AS val(filename)
	ON CONFLICT (filename) 
	DO UPDATE SET filename = EXCLUDED.filename
	RETURNING id, filename`

	rows, err := d.db.Query(ctx, query, fn)

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		return nil, pkg.PackageError(packageName, "create files tasks", err)
	}

	defer rows.Close()

	var result []*models.FileTask
	for rows.Next() {
		var ft models.FileTask
		if err = rows.Scan(&ft.ID, &ft.Filename); err != nil {
			return nil, pkg.PackageError(packageName, "scan inserted files tasks", err)
		}
		result = append(result, &ft)
	}

	if err = rows.Err(); err != nil {
		return nil, pkg.PackageError(packageName, "scan inserted files tasks", err)
	}

	return result, nil
}

func (d *FileTaskRepository) UpdateStatusByID(ctx context.Context, id int, status string) error {
	query := `UPDATE files_tasks SET status = $1 WHERE id = $2`
	_, err := d.db.Exec(ctx, query, status, id)
	if err != nil {
		return pkg.PackageError(packageName, "update status by id", err)
	}
	return nil
}
