package shortener

import (
	"context"
	"shortener/internal/common/textutils"
	"shortener/internal/config"
	"shortener/internal/db/models"
	"shortener/internal/db/models/link"
	"shortener/internal/db/models/linkvisit"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockLinkRepository struct {
	mock.Mock
}

func (m *mockLinkRepository) CreateOne(ctx context.Context, insert link.Insert) (link.Record, error) {
	args := m.Called(ctx, insert)

	return args.Get(0).(link.Record), args.Error(1)
}

func (m *mockLinkRepository) GetByID(ctx context.Context, id uint) (link.Record, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(link.Record), args.Error(1)
}

func (m *mockLinkRepository) GetMany(ctx context.Context, options link.ListOptions) ([]link.Record, error) {
	args := m.Called(ctx, options)

	return args.Get(0).([]link.Record), args.Error(1)
}

func (m *mockLinkRepository) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)

	return args.Int(0), args.Error(1)
}

func (m *mockLinkRepository) UpdateByID(ctx context.Context, id uint, update link.Update) (link.Record, error) {
	args := m.Called(ctx, id, update)

	return args.Get(0).(link.Record), args.Error(1)
}

func (m *mockLinkRepository) DeleteByID(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

type mockLinkVisitRepository struct {
	mock.Mock
}

func (m *mockLinkVisitRepository) CreateOne(ctx context.Context, insert linkvisit.Insert) (linkvisit.Record, error) {
	args := m.Called(ctx, insert)

	return args.Get(0).(linkvisit.Record), args.Error(1)
}

func (m *mockLinkVisitRepository) GetMany(ctx context.Context, options linkvisit.ListOptions) ([]linkvisit.Record, error) {
	args := m.Called(ctx, options)

	return args.Get(0).([]linkvisit.Record), args.Error(1)
}

func (m *mockLinkVisitRepository) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)

	return args.Int(0), args.Error(1)
}

type serviceMocks struct {
	service       *Service
	linkRepo      *mockLinkRepository
	linkVisitRepo *mockLinkVisitRepository
}

func newServiceMocks(t *testing.T) serviceMocks {
	t.Helper()

	linkRepo := new(mockLinkRepository)
	linkVisitRepo := new(mockLinkVisitRepository)

	t.Cleanup(func() {
		linkRepo.AssertExpectations(t)
		linkVisitRepo.AssertExpectations(t)
	})

	appConfig := &config.App{BaseURL: "https://short.example.com"}
	service := NewService(linkRepo, linkVisitRepo, appConfig)

	return serviceMocks{
		service:       service,
		linkRepo:      linkRepo,
		linkVisitRepo: linkVisitRepo,
	}
}

func TestService_CreateLink(t *testing.T) {
	t.Run("should create link successfully", func(t *testing.T) {
		const (
			id          uint = 42
			originalURL      = "https://example.com/some/page"
			shortName        = "example"
			baseURL          = "https://short.example.com"
		)

		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("CreateOne", t.Context(), link.Insert{
				OriginalURL: originalURL,
				ShortName:   shortName,
			}).
			Return(link.Record{
				Model:       gorm.Model{ID: id},
				OriginalURL: originalURL,
				ShortName:   shortName,
			}, nil).
			Once()

		result, err := mocks.service.CreateLink(t.Context(), originalURL, shortName)

		require.NoError(t, err)
		assert.Equal(t, Link{
			ID:          id,
			OriginalURL: originalURL,
			ShortName:   shortName,
			ShortURL:    "https://short.example.com/r/example",
		}, result)
	})

	t.Run("should generate short name when not provided", func(t *testing.T) {
		const (
			id          uint = 42
			originalURL      = "https://example.com/some/page"
		)

		mocks := newServiceMocks(t)
		expectedShortName := textutils.RandomString(6)
		mocks.linkRepo.
			On("CreateOne", t.Context(), mock.MatchedBy(func(insert link.Insert) bool {
				return insert.OriginalURL == originalURL && len(insert.ShortName) == 6
			})).
			Return(link.Record{
				Model:       gorm.Model{ID: id},
				OriginalURL: originalURL,
				ShortName:   expectedShortName,
			}, nil).
			Once()

		result, err := mocks.service.CreateLink(t.Context(), originalURL, "")

		require.NoError(t, err)
		assert.Equal(t, Link{
			ID:          id,
			OriginalURL: originalURL,
			ShortName:   expectedShortName,
			ShortURL:    "https://short.example.com/r/" + expectedShortName,
		}, result)
	})

	t.Run("should return conflict when short name already exists", func(t *testing.T) {
		const (
			originalURL = "https://example.com/some/page"
			shortName   = "example"
		)

		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("CreateOne", t.Context(), link.Insert{
				OriginalURL: originalURL,
				ShortName:   shortName,
			}).
			Return(link.Record{}, models.ErrObjectAlreadyExists).
			Once()

		result, err := mocks.service.CreateLink(t.Context(), originalURL, shortName)

		require.Error(t, err)
		assert.Equal(t, Link{}, result)
	})

	t.Run("should retry when generated short name already exists", func(t *testing.T) {
		const (
			id          uint = 42
			originalURL      = "https://example.com/some/page"
		)

		mocks := newServiceMocks(t)

		mocks.linkRepo.
			On("CreateOne", t.Context(), mock.MatchedBy(func(insert link.Insert) bool {
				return insert.OriginalURL == originalURL && len(insert.ShortName) == 6
			})).
			Return(link.Record{}, models.ErrObjectAlreadyExists).
			Once()

		expectedShortName := textutils.RandomString(6)

		mocks.linkRepo.
			On("CreateOne", t.Context(), mock.MatchedBy(func(insert link.Insert) bool {
				return insert.OriginalURL == originalURL && len(insert.ShortName) == 6
			})).
			Return(link.Record{
				Model:       gorm.Model{ID: id},
				OriginalURL: originalURL,
				ShortName:   expectedShortName,
			}, nil).
			Once()

		result, err := mocks.service.CreateLink(t.Context(), originalURL, "")

		require.NoError(t, err)
		assert.Equal(t, Link{
			ID:          id,
			OriginalURL: originalURL,
			ShortName:   expectedShortName,
			ShortURL:    "https://short.example.com/r/" + expectedShortName,
		}, result)

		mocks.linkRepo.AssertNumberOfCalls(t, "CreateOne", 2)
	})
}

func TestService_GetLink(t *testing.T) {
	t.Run("should get link successfully", func(t *testing.T) {
		const (
			id          uint = 42
			originalURL      = "https://example.com/some/page"
			shortName        = "example"
			baseURL          = "https://short.example.com"
		)

		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("GetByID", t.Context(), id).
			Return(link.Record{
				Model:       gorm.Model{ID: id},
				OriginalURL: originalURL,
				ShortName:   shortName,
			}, nil).
			Once()

		result, err := mocks.service.GetLink(t.Context(), id)

		require.NoError(t, err)
		assert.Equal(t, Link{
			ID:          id,
			OriginalURL: originalURL,
			ShortName:   shortName,
			ShortURL:    "https://short.example.com/r/example",
		}, result)
	})

	t.Run("should return not found when link does not exist", func(t *testing.T) {
		const id uint = 42

		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("GetByID", t.Context(), id).
			Return(link.Record{}, models.ErrObjectDoesNotExist).
			Once()

		result, err := mocks.service.GetLink(t.Context(), id)

		require.Error(t, err)
		assert.Equal(t, Link{}, result)
	})
}

func TestService_GetRedirectLink(t *testing.T) {
	t.Run("should get redirect link successfully", func(t *testing.T) {
		const (
			id          uint = 42
			originalURL      = "https://example.com/some/page"
			shortName        = "example"
		)

		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("GetMany", t.Context(), link.ListOptions{
				ListOptions: models.ListOptions{
					Limit:     1,
					Offset:    0,
					SortBy:    "id",
					SortOrder: "DESC",
				},
				Filters: link.Filters{ShortNames: []string{shortName}},
			}).
			Return([]link.Record{
				{
					Model:       gorm.Model{ID: id},
					OriginalURL: originalURL,
					ShortName:   shortName,
				},
			}, nil).
			Once()

		result, err := mocks.service.GetRedirectLink(t.Context(), shortName)

		require.NoError(t, err)
		assert.Equal(t, Link{
			ID:          id,
			OriginalURL: originalURL,
			ShortName:   shortName,
			ShortURL:    "https://short.example.com/r/example",
		}, result)
	})

	t.Run("should return not found when link does not exist", func(t *testing.T) {
		const shortName = "missing"

		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("GetMany", t.Context(), link.ListOptions{
				ListOptions: models.ListOptions{
					Limit:     1,
					Offset:    0,
					SortBy:    "id",
					SortOrder: "DESC",
				},
				Filters: link.Filters{ShortNames: []string{shortName}},
			}).
			Return([]link.Record{}, nil).
			Once()

		result, err := mocks.service.GetRedirectLink(t.Context(), shortName)

		require.Error(t, err)
		assert.Equal(t, Link{}, result)
	})
}

func TestService_ListLinksWithCount(t *testing.T) {
	t.Run("should list links with count successfully", func(t *testing.T) {
		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("Count", t.Context()).
			Return(42, nil).
			Once()
		mocks.linkRepo.
			On("GetMany", t.Context(), link.ListOptions{
				ListOptions: models.ListOptions{
					Limit:     10,
					Offset:    0,
					SortBy:    "id",
					SortOrder: "DESC",
				},
			}).
			Return([]link.Record{
				{
					Model:       gorm.Model{ID: 1},
					OriginalURL: "https://example.com/first",
					ShortName:   "first",
				},
				{
					Model:       gorm.Model{ID: 2},
					OriginalURL: "https://example.com/second",
					ShortName:   "second",
				},
			}, nil).
			Once()

		result, count, err := mocks.service.ListLinksWithCount(t.Context(), nil)

		require.NoError(t, err)
		assert.Equal(t, 42, count)
		assert.Equal(t, []Link{
			{
				ID:          1,
				OriginalURL: "https://example.com/first",
				ShortName:   "first",
				ShortURL:    "https://short.example.com/r/first",
			},
			{
				ID:          2,
				OriginalURL: "https://example.com/second",
				ShortName:   "second",
				ShortURL:    "https://short.example.com/r/second",
			},
		}, result)
	})

	t.Run("should list links with provided options", func(t *testing.T) {
		mocks := newServiceMocks(t)
		builder := NewLinkListOptionsBuilder()
		builder.WithShortNames("first", "second")
		builder.WithRange(10, 19)
		builder.WithSort("short_name", "asc")
		mocks.linkRepo.
			On("Count", t.Context()).
			Return(0, nil).
			Once()
		mocks.linkRepo.
			On("GetMany", t.Context(), link.ListOptions{
				ListOptions: models.ListOptions{
					Limit:     10,
					Offset:    10,
					SortBy:    "short_name",
					SortOrder: "ASC",
				},
				Filters: link.Filters{
					ShortNames: []string{"first", "second"},
				},
			}).
			Return([]link.Record{}, nil).
			Once()

		result, count, err := mocks.service.ListLinksWithCount(t.Context(), builder)

		require.NoError(t, err)
		assert.Empty(t, result)
		assert.Zero(t, count)
	})

	t.Run("should return validation error for invalid options", func(t *testing.T) {
		mocks := newServiceMocks(t)
		builder := NewLinkListOptionsBuilder()
		builder.WithRange(-1, 10)

		result, count, err := mocks.service.ListLinksWithCount(t.Context(), builder)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, count)
		mocks.linkRepo.AssertNotCalled(t, "GetMany")
		mocks.linkRepo.AssertNotCalled(t, "Count")
	})
}

func TestService_UpdateLink(t *testing.T) {
	t.Run("should update link successfully", func(t *testing.T) {
		const (
			id          uint = 42
			originalURL      = "https://example.com/updated"
			shortName        = "updated"
		)

		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("UpdateByID", t.Context(), id, link.Update{
				OriginalURL: originalURL,
				ShortName:   shortName,
			}).
			Return(link.Record{
				Model:       gorm.Model{ID: id},
				OriginalURL: originalURL,
				ShortName:   shortName,
			}, nil).
			Once()

		result, err := mocks.service.UpdateLink(t.Context(), id, originalURL, shortName)

		require.NoError(t, err)
		assert.Equal(t, Link{
			ID:          id,
			OriginalURL: originalURL,
			ShortName:   shortName,
			ShortURL:    "https://short.example.com/r/updated",
		}, result)
	})

	t.Run("should return not found when link does not exist", func(t *testing.T) {
		const (
			id          uint = 42
			originalURL      = "https://example.com/updated"
			shortName        = "updated"
		)

		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("UpdateByID", t.Context(), id, link.Update{
				OriginalURL: originalURL,
				ShortName:   shortName,
			}).
			Return(link.Record{}, models.ErrObjectDoesNotExist).
			Once()

		result, err := mocks.service.UpdateLink(t.Context(), id, originalURL, shortName)

		require.Error(t, err)
		assert.Equal(t, Link{}, result)
	})
}

func TestService_DeleteLink(t *testing.T) {
	t.Run("should delete link successfully", func(t *testing.T) {
		const id uint = 42

		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("DeleteByID", t.Context(), id).
			Return(nil).
			Once()

		err := mocks.service.DeleteLink(t.Context(), id)

		require.NoError(t, err)
	})

	t.Run("should return not found when link does not exist", func(t *testing.T) {
		const id uint = 42

		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("DeleteByID", t.Context(), id).
			Return(models.ErrObjectDoesNotExist).
			Once()

		err := mocks.service.DeleteLink(t.Context(), id)

		require.Error(t, err)
	})
}

func TestService_SaveLinkVisit(t *testing.T) {
	t.Run("should save link visit successfully", func(t *testing.T) {
		const (
			linkID    uint = 42
			ip             = "127.0.0.1"
			userAgent      = "Mozilla/5.0"
			referrer       = "https://example.com"
			status         = 302
		)

		createdAt := time.Date(2026, 9, 4, 12, 30, 0, 0, time.UTC)
		updatedAt := time.Date(2026, 9, 4, 12, 35, 0, 0, time.UTC)
		mocks := newServiceMocks(t)
		mocks.linkVisitRepo.
			On("CreateOne", t.Context(), linkvisit.Insert{
				LinkID:    linkID,
				IP:        ip,
				UserAgent: userAgent,
				Referrer:  referrer,
				Status:    status,
			}).
			Return(linkvisit.Record{
				Model: gorm.Model{
					ID:        1,
					CreatedAt: createdAt,
					UpdatedAt: updatedAt,
				},
				LinkID:    linkID,
				IP:        ip,
				UserAgent: userAgent,
				Referrer:  referrer,
				Status:    status,
			}, nil).
			Once()

		result, err := mocks.service.SaveLinkVisit(t.Context(), linkID, ip, userAgent, referrer, status)

		require.NoError(t, err)
		assert.Equal(t, LinkVisit{
			ID:        1,
			LinkID:    linkID,
			CreatedAt: createdAt.Format(time.RFC3339),
			UpdatedAt: updatedAt.Format(time.RFC3339),
			IP:        ip,
			UserAgent: userAgent,
			Referrer:  referrer,
			Status:    status,
		}, result)
	})

	t.Run("should return error when saving link visit fails", func(t *testing.T) {
		const (
			linkID    uint = 42
			ip             = "192.168.0.1"
			userAgent      = "Mozilla/5.0"
			referrer       = "https://example.com"
			status         = 404
		)

		mocks := newServiceMocks(t)
		mocks.linkVisitRepo.
			On("CreateOne", t.Context(), linkvisit.Insert{
				LinkID:    linkID,
				IP:        ip,
				UserAgent: userAgent,
				Referrer:  referrer,
				Status:    status,
			}).
			Return(linkvisit.Record{}, models.ErrObjectAlreadyExists).
			Once()

		result, err := mocks.service.SaveLinkVisit(t.Context(), linkID, ip, userAgent, referrer, status)

		require.Error(t, err)
		assert.Equal(t, LinkVisit{}, result)
	})
}

func TestService_ListLinkVisitsWithCount(t *testing.T) {
	t.Run("should list link visits with count successfully", func(t *testing.T) {
		mocks := newServiceMocks(t)
		builder := NewLinkVisitListOptionsBuilder()

		mocks.linkVisitRepo.
			On("Count", t.Context()).
			Return(42, nil).
			Once()
		mocks.linkVisitRepo.
			On("GetMany", t.Context(), builder.build()).
			Return([]linkvisit.Record{
				{
					Model:     gorm.Model{ID: 1},
					LinkID:    42,
					IP:        "127.0.0.1",
					UserAgent: "Mozilla/5.0",
					Referrer:  "https://example.com",
					Status:    302,
				},
				{
					Model:     gorm.Model{ID: 2},
					LinkID:    42,
					IP:        "192.168.0.1",
					UserAgent: "Mozilla/5.0",
					Referrer:  "https://example.com",
					Status:    404,
				},
			}, nil).
			Once()

		result, count, err := mocks.service.ListLinkVisitsWithCount(t.Context(), builder)

		require.NoError(t, err)
		assert.Equal(t, 42, count)
		assert.Len(t, result, 2)
	})

	t.Run("should list link visits without filters successfully", func(t *testing.T) {
		mocks := newServiceMocks(t)
		mocks.linkVisitRepo.
			On("Count", t.Context()).
			Return(0, nil).
			Once()
		mocks.linkVisitRepo.
			On("GetMany", t.Context(), mock.AnythingOfType("linkvisit.ListOptions")).
			Return([]linkvisit.Record{}, nil).
			Once()

		result, count, err := mocks.service.ListLinkVisitsWithCount(t.Context(), nil)

		require.NoError(t, err)
		assert.Empty(t, result)
		assert.Zero(t, count)
	})

	t.Run("should return validation error for invalid options", func(t *testing.T) {
		mocks := newServiceMocks(t)
		builder := NewLinkVisitListOptionsBuilder()
		builder.WithRange(-1, 10)

		result, count, err := mocks.service.ListLinkVisitsWithCount(t.Context(), builder)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, count)
		mocks.linkVisitRepo.AssertNotCalled(t, "GetMany")
		mocks.linkVisitRepo.AssertNotCalled(t, "Count")
	})

	t.Run("should return error when getting link visits fails", func(t *testing.T) {
		mocks := newServiceMocks(t)
		builder := NewLinkVisitListOptionsBuilder()

		mocks.linkVisitRepo.
			On("GetMany", t.Context(), builder.build()).
			Return([]linkvisit.Record{}, models.ErrObjectDoesNotExist).
			Once()

		result, count, err := mocks.service.ListLinkVisitsWithCount(t.Context(), builder)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, count)
	})

	t.Run("should return error when counting link visits fails", func(t *testing.T) {
		mocks := newServiceMocks(t)
		builder := NewLinkVisitListOptionsBuilder()
		mocks.linkVisitRepo.
			On("GetMany", t.Context(), builder.build()).
			Return([]linkvisit.Record{}, nil).
			Once()
		mocks.linkVisitRepo.
			On("Count", t.Context()).
			Return(0, models.ErrObjectDoesNotExist).
			Once()

		result, count, err := mocks.service.ListLinkVisitsWithCount(t.Context(), builder)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, count)
	})
}
