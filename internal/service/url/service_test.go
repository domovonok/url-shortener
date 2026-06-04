package url_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	xerrors "github.com/domovonok/url-shortener/internal/errors"
	urlservice "github.com/domovonok/url-shortener/internal/service/url"
	"github.com/domovonok/url-shortener/internal/service/url/mocks"
)

func TestServiceCreateSuccess(t *testing.T) {
	ctx, repo, service := newTestService(t)
	rawURL := "https://example.com/path"
	id := int64(1)

	repo.EXPECT().Create(ctx, rawURL).Return(id, nil)

	code, err := service.Create(ctx, rawURL)

	require.NoError(t, err)
	require.NotEmpty(t, code)

	repo.EXPECT().Get(ctx, id).Return(rawURL, nil)

	got, err := service.Get(ctx, code)

	require.NoError(t, err)
	require.Equal(t, rawURL, got)
}

func TestServiceCreateInvalidURL(t *testing.T) {
	ctx, _, service := newTestService(t)

	code, err := service.Create(ctx, "example.com")

	require.Empty(t, code)
	require.ErrorIs(t, err, xerrors.ErrInvalidUrl)
}

func TestServiceCreateDuplicateURL(t *testing.T) {
	ctx, repo, service := newTestService(t)
	rawURL := "https://example.com/path"
	id := int64(42)

	repo.EXPECT().Create(ctx, rawURL).Return(int64(0), xerrors.ErrUrlExists)
	repo.EXPECT().GetIDByURL(ctx, rawURL).Return(id, nil)

	code, err := service.Create(ctx, rawURL)

	require.NoError(t, err)
	require.NotEmpty(t, code)

	repo.EXPECT().Get(ctx, id).Return(rawURL, nil)

	got, err := service.Get(ctx, code)

	require.NoError(t, err)
	require.Equal(t, rawURL, got)
}

func TestServiceCreateRepoError(t *testing.T) {
	ctx, repo, service := newTestService(t)
	rawURL := "https://example.com/path"
	repoErr := errors.New("repo error")

	repo.EXPECT().Create(ctx, rawURL).Return(int64(0), repoErr)

	code, err := service.Create(ctx, rawURL)

	require.Empty(t, code)
	require.ErrorIs(t, err, repoErr)
}

func TestServiceCreateDuplicateLookupError(t *testing.T) {
	ctx, repo, service := newTestService(t)
	rawURL := "https://example.com/path"
	lookupErr := errors.New("lookup error")

	repo.EXPECT().Create(ctx, rawURL).Return(int64(0), xerrors.ErrUrlExists)
	repo.EXPECT().GetIDByURL(ctx, rawURL).Return(int64(0), lookupErr)

	code, err := service.Create(ctx, rawURL)

	require.Empty(t, code)
	require.ErrorIs(t, err, lookupErr)
}

func TestServiceGetSuccess(t *testing.T) {
	ctx, repo, service := newTestService(t)
	id := int64(42)
	rawURL := "https://example.com/path"
	code := createCode(t, ctx, repo, service, rawURL, id)

	repo.EXPECT().Get(ctx, id).Return(rawURL, nil)

	got, err := service.Get(ctx, code)

	require.NoError(t, err)
	require.Equal(t, rawURL, got)
}

func TestServiceGetInvalidCode(t *testing.T) {
	ctx, _, service := newTestService(t)

	rawURL, err := service.Get(ctx, "short")

	require.Empty(t, rawURL)
	require.ErrorIs(t, err, xerrors.ErrInvalidCode)
}

func TestServiceGetRepoError(t *testing.T) {
	ctx, repo, service := newTestService(t)
	id := int64(7)
	rawURL := "https://example.com/path"
	code := createCode(t, ctx, repo, service, rawURL, id)

	repo.EXPECT().Get(ctx, id).Return("", xerrors.ErrUrlNotFound)

	got, err := service.Get(ctx, code)

	require.Empty(t, got)
	require.ErrorIs(t, err, xerrors.ErrUrlNotFound)
}

type urlService interface {
	Create(context.Context, string) (string, error)
	Get(context.Context, string) (string, error)
}

func newTestService(t *testing.T) (context.Context, *mocks.MockRepo, urlService) {
	t.Helper()

	ctrl := gomock.NewController(t)
	repo := mocks.NewMockRepo(ctrl)

	return context.Background(), repo, urlservice.NewService(repo)
}

func createCode(
	t *testing.T,
	ctx context.Context,
	repo *mocks.MockRepo,
	service urlService,
	rawURL string,
	id int64,
) string {
	t.Helper()

	repo.EXPECT().Create(ctx, rawURL).Return(id, nil)

	code, err := service.Create(ctx, rawURL)
	require.NoError(t, err)
	require.NotEmpty(t, code)

	return code
}
