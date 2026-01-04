package link

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"

	"github.com/domovonok/url-shortener/internal/logger"
	"github.com/domovonok/url-shortener/internal/model"
	"github.com/domovonok/url-shortener/internal/transport/http/dto/link"
)

type Controller struct {
	create createUsecase
	get    getUsecase
	log    logger.Logger
}

func New(c createUsecase, g getUsecase, l logger.Logger) *Controller {
	return &Controller{create: c, get: g, log: l}
}

func (c *Controller) Create(w http.ResponseWriter, r *http.Request) {
	var req link.CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		c.responseError(w, model.ErrInvalidInput)
		return
	}

	res, err := c.create.Create(r.Context(), req.Url)
	if err != nil {
		c.responseError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}

func (c *Controller) Get(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")

	res, err := c.get.Get(r.Context(), code)
	if err != nil {
		c.responseError(w, err)
		return
	}

	http.Redirect(w, r, res.Url, http.StatusMovedPermanently)
}

func (c *Controller) responseError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, model.ErrInvalidInput):
		http.Error(w, `{"error": "Invalid Input"}`, http.StatusBadRequest)
	case errors.Is(err, model.ErrCodeNotFound):
		http.Error(w, `{"error": "Code Not Found"}`, http.StatusBadRequest)
	default:
		c.log.Error("Internal error", logger.Error(err))
		http.Error(w, `{"error": "Internal Server Error"}`, http.StatusInternalServerError)
	}
}
