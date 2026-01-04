package get_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/domovonok/url-shortener/internal/model"
	"github.com/domovonok/url-shortener/internal/usecase/link/get"
	get_test "github.com/domovonok/url-shortener/internal/usecase/link/get/mocks"
)

func TestGet(t *testing.T) {
	t.Parallel()

	t.Run("success", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		repo := get_test.NewMocklinkRepo(ctrl)
		uc := get.New(repo)

		code := "Code123"
		want := model.Link{
			Url:       "https://test.com/some/path/1",
			Code:      code,
			CreatedAt: time.Unix(123, 0).UTC(),
		}

		repo.EXPECT().
			Get(gomock.Any(), code).
			Return(want, nil)

		got, err := uc.Get(ctx, code)
		require.NoError(t, err)
		require.Equal(t, want, got)
	})

	t.Run("error", func(t *testing.T) {
		t.Parallel()

		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		ctx := context.Background()
		repo := get_test.NewMocklinkRepo(ctrl)
		uc := get.New(repo)

		code := "Code123"
		wantErr := errors.New("repo failure")

		repo.EXPECT().
			Get(gomock.Any(), code).
			Return(model.Link{}, wantErr)

		got, err := uc.Get(ctx, code)
		require.ErrorIs(t, wantErr, err)
		require.Empty(t, got)
	})
}
