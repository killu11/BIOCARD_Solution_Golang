package app

import (
	"context"
	"directory-viewing-service/internal/config"
	"directory-viewing-service/internal/domain/repository"
	is "directory-viewing-service/internal/domain/services"
	"directory-viewing-service/internal/infrastructure/broker"
	http2 "directory-viewing-service/internal/infrastructure/http"
	"directory-viewing-service/internal/infrastructure/persistence"
	"directory-viewing-service/internal/infrastructure/services"
	"directory-viewing-service/internal/infrastructure/workers"
	"directory-viewing-service/internal/infrastructure/workers/parsers"
	"directory-viewing-service/pkg"
	"log"
	"log/slog"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
)

type handlers struct {
	fd *http2.FileDataHandler
}
type repos struct {
	ft     repository.FileTaskRepository
	fd     repository.FileDataRepository
	report repository.ReportRepository
}

type service struct {
	ft     is.FileTaskService
	fd     is.FileDataService
	report is.ReportService
}

type App struct {
	handlers handlers
	services service
	repos    repos
	w        *workers.FileWatcher
	p        *workers.Processor
}

func NewApp() *App {
	cfg, err := config.NewConfig()
	if err != nil {
		log.Fatalln(err)
	}

	db, err := persistence.NewPostgresConnections(cfg.Postgres)
	if err != nil {
		log.Fatalln(err)
	}

	rabbit, err := broker.NewRabbitMQ(cfg.RabbitMQ)
	if err != nil {
		log.Fatalln(err)
	}

	l := pkg.NewLogger()

	//Repositories, services, handlers
	r := repos{
		ft:     persistence.NewFileTaskRepository(db),
		fd:     persistence.NewFileDataRepository(db),
		report: persistence.NewReportRepository(db),
	}

	s := service{
		ft:     services.NewFileTaskService(r.ft),
		fd:     services.NewFileDataService(r.fd),
		report: services.NewReportService(r.report),
	}

	h := handlers{
		fd: http2.NewFileDataHandler(s.fd),
	}
	sdr := broker.NewSender(rabbit.Sch)
	rvr := broker.NewReceiver(rabbit.Rch)
	parser := parsers.NewTSVParser(cfg.Directory.In, cfg.Directory.Out)

	fileWatcher := workers.NewFileWatcher(
		cfg.Directory,
		s.ft,
		sdr,
		l,
	)

	processer := workers.NewProcessor(
		parser,
		s.fd,
		s.ft,
		s.report,
		rvr,
		l,
		&sync.Mutex{},
	)
	return &App{
		w:        fileWatcher,
		p:        processer,
		services: s,
		repos:    r,
		handlers: h,
	}
}
func (a *App) Start() {
	ctx, cancel := context.WithCancel(context.Background())

	go a.w.Start(ctx)
	a.p.Start(ctx)
	a.runServer(cancel)
}

func (a *App) runServer(cancel context.CancelFunc) {
	r := mux.NewRouter()
	port := 8080

	a.handlers.fd.InitRoutes(r)
	defer cancel()

	log.Printf("Service start on port: %d\n", port)
	if err := http.ListenAndServe(":8080", r); err != nil {
		slog.Error("Service fatal: ", slog.StringValue(err.Error()))
	}
}
