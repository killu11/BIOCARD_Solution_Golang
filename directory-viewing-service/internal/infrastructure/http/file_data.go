package http

import (
	"context"
	"directory-viewing-service/internal/domain/services"
	is "directory-viewing-service/internal/infrastructure/services"
	"directory-viewing-service/pkg"
	"errors"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

var ErrParmsNotFound = errors.New("url parms not found")
var ErrInvalidLimit = errors.New("limit parameter must be greater than 0 or equal to -1 for return all records")

type FileDataHandler struct {
	service services.FileDataService
}

func NewFileDataHandler(service services.FileDataService) *FileDataHandler {
	return &FileDataHandler{service: service}
}

func (h *FileDataHandler) InitRoutes(mux *mux.Router) {
	mux.HandleFunc("/{uid}", h.getRows).Methods("GET")
}

func (h *FileDataHandler) getRows(w http.ResponseWriter, r *http.Request) {
	parms := mux.Vars(r)

	uid, ok := parms["uid"]
	if !ok {
		h.matchError(w, ErrParmsNotFound)
		return
	}
	var limit int
	stringLim := r.URL.Query().Get("limit")
	if stringLim == "" {
		limit = -1
	}
	limit, err := strconv.Atoi(stringLim)
	if err != nil {
		pkg.JSONError(w, ErrInvalidLimit.Error(), http.StatusBadRequest)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	rows, err := h.service.ProcessedData(ctx, uid, limit)
	if err != nil {
		h.matchError(w, err)
		return
	}

	pkg.JSONResponse(w, rows, http.StatusOK)
}

func (h *FileDataHandler) matchError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrParmsNotFound):
		pkg.JSONError(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, is.ErrInvalidLimit):
		pkg.JSONError(w, ErrInvalidLimit.Error(), http.StatusBadRequest)
	case errors.Is(err, is.ErrZeroParsedRows):
		pkg.JSONError(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, is.ErrInvalidUID):
		pkg.JSONError(w, err.Error(), http.StatusBadRequest)
	case errors.Is(err, is.ErrEmptyData):
		pkg.JSONResponse(w, map[string]string{}, http.StatusOK)
	default:
		slog.Warn("unexcepted result", "error", err)
		pkg.JSONError(w, "internal server error", http.StatusInternalServerError)
	}
}
