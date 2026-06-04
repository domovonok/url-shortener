package inmemory_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	xerrors "github.com/domovonok/url-shortener/internal/errors"
	urlinmemoryrepo "github.com/domovonok/url-shortener/internal/repository/url/inmemory"
)

func TestRepositoryCreateAndGet(t *testing.T) {
	ctx := context.Background()
	repo := urlinmemoryrepo.NewRepository()
	rawURL := "https://example.com/path"

	id, err := repo.Create(ctx, rawURL)
	require.NoError(t, err)
	require.Equal(t, int64(1), id)

	gotURL, err := repo.Get(ctx, id)
	require.NoError(t, err)
	require.Equal(t, rawURL, gotURL)

	gotID, err := repo.GetIDByURL(ctx, rawURL)
	require.NoError(t, err)
	require.Equal(t, id, gotID)
}

func TestRepositoryCreateDuplicateURL(t *testing.T) {
	ctx := context.Background()
	repo := urlinmemoryrepo.NewRepository()
	rawURL := "https://example.com/path"

	id, err := repo.Create(ctx, rawURL)
	require.NoError(t, err)
	require.Equal(t, int64(1), id)

	duplicateID, err := repo.Create(ctx, rawURL)
	require.ErrorIs(t, err, xerrors.ErrUrlExists)
	require.Zero(t, duplicateID)

	gotID, err := repo.GetIDByURL(ctx, rawURL)
	require.NoError(t, err)
	require.Equal(t, id, gotID)
}

func TestRepositoryCreateUsesSequentialIDs(t *testing.T) {
	ctx := context.Background()
	repo := urlinmemoryrepo.NewRepository()

	firstID, err := repo.Create(ctx, "https://example.com/first")
	require.NoError(t, err)
	require.Equal(t, int64(1), firstID)

	secondID, err := repo.Create(ctx, "https://example.com/second")
	require.NoError(t, err)
	require.Equal(t, int64(2), secondID)
}

func TestRepositoryGetMissingURL(t *testing.T) {
	ctx := context.Background()
	repo := urlinmemoryrepo.NewRepository()

	rawURL, err := repo.Get(ctx, 1)
	require.ErrorIs(t, err, xerrors.ErrUrlNotFound)
	require.Empty(t, rawURL)
}

func TestRepositoryGetIDByMissingURL(t *testing.T) {
	ctx := context.Background()
	repo := urlinmemoryrepo.NewRepository()

	id, err := repo.GetIDByURL(ctx, "https://example.com/missing")
	require.ErrorIs(t, err, xerrors.ErrUrlNotFound)
	require.Zero(t, id)
}
