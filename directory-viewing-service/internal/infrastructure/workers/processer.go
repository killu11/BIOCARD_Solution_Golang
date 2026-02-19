package workers

import (
	"context"
	"directory-viewing-service/internal/domain/services"
	"directory-viewing-service/internal/domain/workers/parser"
	"directory-viewing-service/internal/infrastructure/broker"
	"directory-viewing-service/internal/infrastructure/dto"
	"directory-viewing-service/pkg"
	"log/slog"
	"os"
	"sync"
	"time"
)

type Processor struct {
	parser          parser.TSVParser
	fileDataService services.FileDataService
	fileTaskService services.FileTaskService
	reportService   services.ReportService
	receiver        *broker.Receiver
	logger          *slog.Logger
	mutex           *sync.Mutex
}

func NewProcessor(
	parser parser.TSVParser,
	fds services.FileDataService,
	fts services.FileTaskService,
	rs services.ReportService,
	receiver *broker.Receiver,
	logger *slog.Logger,
	mutex *sync.Mutex,
) *Processor {
	return &Processor{
		parser:          parser,
		fileDataService: fds,
		fileTaskService: fts,
		reportService:   rs,
		receiver:        receiver,
		logger:          logger,
		mutex:           mutex,
	}
}

func (p *Processor) Start(ctx context.Context) {
	const numWorkers = 5
	out := make(chan *dto.FileDataMessage, numWorkers)
	if err := p.receiver.Receive(ctx, out, p.fileTaskService); err != nil {
		p.logger.Error("start receive", "error", err)
	}

	for i := 0; i < 5; i++ {
		go func() {
			for fm := range out {
				p.processMessage(ctx, fm)
			}
		}()
	}
}

func (p *Processor) processMessage(ctx context.Context, msg *dto.FileDataMessage) {
	path := p.parser.TSVPath(msg.Filename)
	f, err := os.Open(path)
	if err != nil {
		p.logger.Warn(
			"open file",
			"filename", msg.Filename,
			"error", err,
		)
		return
	}

	defer f.Close()

	rows, err := p.parser.Parse(f, msg.Filename)
	if err != nil {
		p.handleFailure(ctx, msg.ID, msg.Filename, err)
		return
	}

	uploadCtx, cancel := context.WithTimeout(ctx, time.Second*10)
	defer cancel()
	if err = p.fileDataService.UploadParsed(uploadCtx, rows); err != nil {
		p.handleFailure(ctx, msg.ID, msg.Filename, err)
		return
	}
	p.mutex.Lock()
	p.parser.PublishPDF(rows)
	p.mutex.Unlock()
}

func (p *Processor) handleFailure(ctx context.Context, id int, filename string, err error) {
	handleCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	_ = p.fileTaskService.ChangeStatus(ctx, id, services.StatusFailed)
	msg := pkg.ErrDescription(err)

	if err = p.reportService.ReportError(handleCtx, filename, msg); err != nil {
		p.logger.Warn(
			"generate report",
			"error", err.Error(),
			"filename", filename,
			"msg", msg,
		)
	}
}
