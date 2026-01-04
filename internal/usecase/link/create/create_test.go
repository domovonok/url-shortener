package create_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	create_test "github.com/domovonok/url-shortener/internal/usecase/link/create/mocks"

	"github.com/domovonok/url-shortener/internal/model"
	"github.com/domovonok/url-shortener/internal/usecase/link/create"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		repo := create_test.NewMocklinkRepo(ctrl)
		uc := create.New(repo)

		url := "https://test.com/some/path/1"
		want := model.Link{
			Url:       url,
			Code:      "Code123",
			CreatedAt: time.Unix(123, 0).UTC(),
		}

		repo.EXPECT().
			Create(gomock.Any(), url).
			Return(want, nil)

		got, err := uc.Create(ctx, url)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		repo := create_test.NewMocklinkRepo(ctrl)
		uc := create.New(repo)

		url := "https://test.com/some/path/1"
		wantErr := errors.New("repo failure")

		repo.EXPECT().
			Create(gomock.Any(), url).
			Return(model.Link{}, wantErr)

		got, err := uc.Create(ctx, url)
		require.Error(t, err)
		require.ErrorIs(t, wantErr, err)
		require.Empty(t, got)
	})
}
