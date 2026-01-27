package integration

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/domovonok/url-shortener/internal/logger"
	"github.com/domovonok/url-shortener/internal/model"
	linkRepo "github.com/domovonok/url-shortener/internal/repo/link"
	"github.com/domovonok/url-shortener/internal/transport/http/dto/link"
	linkHandler "github.com/domovonok/url-shortener/internal/transport/http/link"
	linkCreateUsecase "github.com/domovonok/url-shortener/internal/usecase/link/create"
	linkGetUsecase "github.com/domovonok/url-shortener/internal/usecase/link/get"
)

func TestLinkController_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	ctx := context.Background()

	pool, err := pgxpool.New(ctx, ConnectionString)
	require.NoError(t, err)
	t.Cleanup(func() {
		pool.Close()
	})

	l := logger.MustInit(true)
	repo := linkRepo.New(pool)
	createUC := linkCreateUsecase.New(repo)
	getUC := linkGetUsecase.New(repo)
	controller := linkHandler.New(createUC, getUC, l)

	t.Run("Successfully create and get link", func(t *testing.T) {
		originalURL := "https://test.com/qwerty123_-"
		reqBody := link.CreateRequest{
			Url: originalURL,
		}

		jsonData, err := json.Marshal(reqBody)
		require.NoError(t, err)

		req := httptest.NewRequest("POST", "/", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		controller.Create(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var createdLink model.Link
		err = json.Unmarshal(w.Body.Bytes(), &createdLink)
		require.NoError(t, err)

		assert.Equal(t, originalURL, createdLink.Url)
		assert.NotEmpty(t, createdLink.Code)
		assert.False(t, createdLink.CreatedAt.IsZero())

		r := chi.NewRouter()
		r.Get("/{code}", controller.Get)

		reqGet := httptest.NewRequest("GET", "/"+createdLink.Code, nil)
		wGet := httptest.NewRecorder()

		r.ServeHTTP(wGet, reqGet)

		assert.Equal(t, http.StatusMovedPermanently, wGet.Code)
		assert.Equal(t, originalURL, wGet.Header().Get("Location"))
	})

	t.Run("Get non-existent link returns error", func(t *testing.T) {
		r := chi.NewRouter()
		r.Get("/{code}", controller.Get)

		reqGet := httptest.NewRequest("GET", "/nonexistent", nil)
		wGet := httptest.NewRecorder()

		r.ServeHTTP(wGet, reqGet)

		assert.Equal(t, http.StatusBadRequest, wGet.Code)
		assert.Contains(t, wGet.Body.String(), "Code Not Found")
	})

	t.Run("Create with invalid input returns error", func(t *testing.T) {
		req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte("invalid json")))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		controller.Create(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
		assert.Contains(t, w.Body.String(), "Invalid Input")
	})
}
