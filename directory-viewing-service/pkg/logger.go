package pkg

import (
	"io"
	"log/slog"
	"os"
)

func NewLogger() *slog.Logger {
	var w io.Writer
	logDir := "../logs"

	if _, err := os.Stat(logDir); err != nil {
		_ = os.Mkdir("../logs/", 0755)
	}

	f, err := os.OpenFile(logDir+"/logs.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		slog.Warn("open file logs.log:", "error", err)
		w = os.Stdout
	} else {
		w = io.MultiWriter(os.Stdout, f)
	}

	return slog.New(slog.NewTextHandler(w, nil))
}
