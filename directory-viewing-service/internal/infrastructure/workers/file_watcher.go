package workers

import (
	"context"
	"directory-viewing-service/internal/config"
	"directory-viewing-service/internal/domain/services"
	"directory-viewing-service/internal/infrastructure/broker"
	"directory-viewing-service/internal/infrastructure/dto"
	"directory-viewing-service/pkg"
	"errors"
	"fmt"
	"log/slog"

	"path/filepath"
	"time"
)

var (
	ErrEmptyDir = errors.New("work directory is empty")
)

const packageName = "infrastructure/workers"

// FileWatcher
// Сканирует директорию на наличие необработанных файлов
type FileWatcher struct {
	watchInterval time.Duration
	workPath      string
	service       services.FileTaskService
	sender        *broker.Sender
	logger        *slog.Logger
}

func NewFileWatcher(
	c *config.DirectoryCfg,
	s services.FileTaskService,
	sender *broker.Sender,
	logger *slog.Logger,
) *FileWatcher {
	return &FileWatcher{
		watchInterval: c.WatchInterval,
		workPath:      c.In,
		service:       s,
		sender:        sender,
		logger:        logger,
	}
}

func (f *FileWatcher) scanDirectory() ([]string, error) {
	pattern := fmt.Sprintf("%s/%s", f.workPath, "*.tsv")
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, pkg.PackageError(packageName, "find .tsv matches", err)
	}
	if len(matches) == 0 {
		return nil, ErrEmptyDir
	}
	f.handleMatches(matches)
	return matches, nil
}

func (f *FileWatcher) handleMatches(matches []string) {
	for i, m := range matches {
		matches[i] = filepath.Base(m)
	}
}

// Start
// Запускает worker для отслеживания не обработанных файлов,
// логгирует непредвиденные кейсы в stdout и файл ./logs/.logs, ставит задачи для обработки,
// отправляя id задачи имя файла для обработки
func (f *FileWatcher) Start(ctx context.Context) {
	f.logger.Info(fmt.Sprintf("File Watcher starts! Scans directory every %s", f.watchInterval.String()))
	ticker := time.NewTicker(f.watchInterval)

	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			f.logger.Debug("File Watcher stops")
			return
		case <-ticker.C:
			names, err := f.scanDirectory()
			if err != nil {
				if errors.Is(err, ErrEmptyDir) {
					continue
				}
				f.logger.Error(
					"Scan work directory:",
					"error", err,
				)
				return
			}

			tasks, err := f.service.UploadTasks(ctx, names)
			if err != nil {
				if errors.Is(err, services.ErrZeroUnProcessedFiles) || errors.Is(err, services.ErrCreatedZeroTasks) {
					f.logger.Info("Zero unprocessed files from work directory")
					continue
				}
				f.logger.Warn(
					"failed to upload tasks",
					"error", err,
				)
			}

			for _, task := range tasks {
				msg := dto.TaskToMessage(task)
				if err := f.sender.Send(msg); err != nil {
					f.logger.Warn(
						"failed to send message",
						"error", err,
						"filename", msg.Filename,
					)
				}
			}
		}
	}
}
