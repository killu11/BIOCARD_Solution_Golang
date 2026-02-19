package services

import "context"

type ReportService interface {
	ReportError(ctx context.Context, filename, msg string) error
}
