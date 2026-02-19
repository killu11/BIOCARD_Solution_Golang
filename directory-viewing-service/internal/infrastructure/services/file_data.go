package services

import (
	"context"
	"directory-viewing-service/internal/domain/models"
	"directory-viewing-service/internal/domain/repository"
	"errors"
)

var (
	ErrInvalidLimit   = errors.New("invalid limit parm")
	ErrInvalidUID     = errors.New("invalid unit_guid")
	ErrEmptyData      = errors.New("got zero records")
	ErrZeroParsedRows = errors.New("got zero parse rows")
)

type FileDataService struct {
	repo repository.FileDataRepository
}

func (f *FileDataService) UploadParsed(ctx context.Context, data []*models.FileData) error {
	if len(data) == 0 {
		return ErrZeroParsedRows
	}

	return f.repo.Save(ctx, data)
}

func (f *FileDataService) ProcessedData(ctx context.Context, uid string, limit int) ([]*models.FileData, error) {
	if limit <= -2 {
		return nil, ErrInvalidLimit
	}
	if uid == "" {
		return nil, ErrInvalidUID
	}

	records, err := f.repo.RecordsByUID(ctx, uid, limit)
	if err != nil {
		return nil, err
	}

	if len(records) == 0 {
		return nil, ErrEmptyData
	}

	return records, nil
}

func NewFileDataService(repo repository.FileDataRepository) *FileDataService {
	return &FileDataService{repo: repo}
}
