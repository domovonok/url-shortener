package url_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"
	"go.uber.org/zap"

	xerrors "github.com/domovonok/url-shortener/internal/errors"
	urlhandler "github.com/domovonok/url-shortener/internal/handler/url"
	"github.com/domovonok/url-shortener/internal/handler/url/mocks"
)

type testHandler interface {
	Create(http.ResponseWriter, *http.Request)
	Get(http.ResponseWriter, *http.Request)
}

func TestHandlerCreateSuccess(t *testing.T) {
	service, handler := newTestHandler(t)
	rawURL := "https://example.com/path"
	code := "abcdef"

	service.EXPECT().Create(gomock.Any(), rawURL).Return(code, nil)

	response := performCreate(handler, newCreateRequest(t, rawURL))

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "application/json", response.Header().Get("Content-Type"))
	require.Equal(t, map[string]string{"code": code}, readJSONResponse(t, response))
}

func TestHandlerCreateInvalidJSON(t *testing.T) {
	_, handler := newTestHandler(t)

	response := performCreate(handler, httptest.NewRequest(http.MethodPost, "/", strings.NewReader("{")))

	require.Equal(t, http.StatusBadRequest, response.Code)
	require.Contains(t, readJSONResponse(t, response)["error"], xerrors.ErrInvalidInput.Error())
}

func TestHandlerCreateServiceError(t *testing.T) {
	service, handler := newTestHandler(t)
	rawURL := "example.com"

	service.EXPECT().Create(gomock.Any(), rawURL).Return("", xerrors.ErrInvalidUrl)

	response := performCreate(handler, newCreateRequest(t, rawURL))

	require.Equal(t, http.StatusBadRequest, response.Code)
	require.Equal(t, map[string]string{"error": xerrors.ErrInvalidUrl.Error()}, readJSONResponse(t, response))
}

func TestHandlerGetSuccess(t *testing.T) {
	service, handler := newTestHandler(t)
	code := "abcdef"
	rawURL := "https://example.com/path"

	service.EXPECT().Get(gomock.Any(), code).Return(rawURL, nil)

	response := performGet(handler, code)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "application/json", response.Header().Get("Content-Type"))
	require.Equal(t, map[string]string{"url": rawURL}, readJSONResponse(t, response))
}

func TestHandlerGetNotFound(t *testing.T) {
	service, handler := newTestHandler(t)
	code := "abcdef"

	service.EXPECT().Get(gomock.Any(), code).Return("", xerrors.ErrUrlNotFound)

	response := performGet(handler, code)

	require.Equal(t, http.StatusNotFound, response.Code)
	require.Equal(t, map[string]string{"error": xerrors.ErrUrlNotFound.Error()}, readJSONResponse(t, response))
}

func TestHandlerGetInternalError(t *testing.T) {
	service, handler := newTestHandler(t)
	code := "abcdef"

	service.EXPECT().Get(gomock.Any(), code).Return("", errors.New("service error"))

	response := performGet(handler, code)

	require.Equal(t, http.StatusInternalServerError, response.Code)
	require.Equal(t, map[string]string{"error": http.StatusText(http.StatusInternalServerError)}, readJSONResponse(t, response))
}

func newTestHandler(t *testing.T) (*mocks.Mockservice, testHandler) {
	t.Helper()

	ctrl := gomock.NewController(t)
	service := mocks.NewMockservice(ctrl)

	return service, urlhandler.NewHandler(zap.NewNop(), service)
}

func performCreate(handler testHandler, request *http.Request) *httptest.ResponseRecorder {
	response := httptest.NewRecorder()
	handler.Create(response, request)

	return response
}

func performGet(handler testHandler, code string) *httptest.ResponseRecorder {
	request := httptest.NewRequest(http.MethodGet, "/"+code, nil)
	routeCtx := chi.NewRouteContext()
	routeCtx.URLParams.Add("code", code)
	request = request.WithContext(context.WithValue(request.Context(), chi.RouteCtxKey, routeCtx))

	response := httptest.NewRecorder()
	handler.Get(response, request)

	return response
}

func newCreateRequest(t *testing.T, rawURL string) *http.Request {
	t.Helper()

	body, err := json.Marshal(map[string]string{"url": rawURL})
	require.NoError(t, err)

	return httptest.NewRequest(http.MethodPost, "/", strings.NewReader(string(body)))
}

func readJSONResponse(t *testing.T, response *httptest.ResponseRecorder) map[string]string {
	t.Helper()

	var body map[string]string
	require.NoError(t, json.NewDecoder(response.Body).Decode(&body))

	return body
}
