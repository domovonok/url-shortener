package url

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"go.uber.org/zap"

	xerrors "github.com/domovonok/url-shortener/internal/errors"
)

//go:generate mockgen -source=handler.go -destination=mocks/handler_mocks.go -package=mocks
type service interface {
	Create(ctx context.Context, url string) (string, error)
	Get(ctx context.Context, url string) (string, error)
}

type handler struct {
	logger  *zap.Logger
	service service
}

func NewHandler(logger *zap.Logger, service service) *handler {
	return &handler{logger, service}
}

func (h *handler) Create(w http.ResponseWriter, r *http.Request) {
	var req createRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.responseError(w, errors.Join(xerrors.ErrInvalidInput, err))
		return
	}

	code, err := h.service.Create(r.Context(), req.Url)
	if err != nil {
		h.responseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(createResponse{Code: code})
}

func (h *handler) Get(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	url, err := h.service.Get(r.Context(), code)
	if err != nil {
		h.responseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(getResponse{Url: url})
}

func (h *handler) responseError(w http.ResponseWriter, err error) {
	statusCode := httpStatusCode(err)
	message := err.Error()

	if statusCode == http.StatusInternalServerError {
		h.logger.Error("internal server error", zap.Error(err))
		message = http.StatusText(statusCode)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(errorResponse{Error: message})
}

func httpStatusCode(err error) int {
	switch {
	case errors.Is(err, xerrors.ErrInvalidInput),
		errors.Is(err, xerrors.ErrInvalidUrl),
		errors.Is(err, xerrors.ErrInvalidCode):
		return http.StatusBadRequest
	case errors.Is(err, xerrors.ErrUrlNotFound):
		return http.StatusNotFound
	default:
		return http.StatusInternalServerError
	}
}
