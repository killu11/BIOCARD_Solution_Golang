package parser

import (
	"directory-viewing-service/internal/domain/models"
	"io"
)

type TSVParser interface {
	Parse(rd io.Reader, filename string) ([]*models.FileData, error)
	PublishPDF([]*models.FileData)
	TSVPath(name string) string
}
