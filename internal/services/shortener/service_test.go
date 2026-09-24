package shortener

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"shortener/internal/common/textutils"
)

type mockLinkRepository struct {
	mock.Mock
}

func (m *mockLinkRepository) CreateOne(ctx context.Context, insert CreateLinkParams) (Link, error) {
	args := m.Called(ctx, insert)

	return args.Get(0).(Link), args.Error(1)
}

func (m *mockLinkRepository) GetByID(ctx context.Context, id uint) (Link, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(Link), args.Error(1)
}

func (m *mockLinkRepository) GetMany(ctx context.Context, options LinkListOptions) ([]Link, error) {
	args := m.Called(ctx, options)

	return args.Get(0).([]Link), args.Error(1)
}

func (m *mockLinkRepository) Count(ctx context.Context) (int, error) {
	args := m.Called(ctx)

	return args.Int(0), args.Error(1)
}

func (m *mockLinkRepository) UpdateByID(ctx context.Context, id uint, update UpdateLinkParams) (Link, error) {
	args := m.Called(ctx, id, update)

	return args.Get(0).(Link), args.Error(1)
}

func (m *mockLinkRepository) DeleteByID(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)

	return args.Error(0)
}

type mockLinkVisitRepository struct {
	mock.Mock
}

func (m *mockLinkVisitRepository) CreateOne(ctx context.Context, insert CreateLinkVisitParams) (LinkVisit, error) {
	args := m.Called(ctx, insert)

	return args.Get(0).(LinkVisit), args.Error(1)
}

func (m *mockLinkVisitRepository) GetMany(ctx context.Context, options LinkVisitListOptions) ([]LinkVisit, error) {
	args := m.Called(ctx, options)

	return args.Get(0).([]LinkVisit), args.Error(1)
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

	service := NewService(linkRepo, linkVisitRepo)

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
		)

		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("CreateOne", t.Context(), CreateLinkParams{
				OriginalURL: originalURL,
				ShortName:   shortName,
			}).
			Return(Link{
				ID:          id,
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
		}, result)
	})

	t.Run("should generate short name when not provided", func(t *testing.T) {
		const (
			id          uint = 42
			originalURL      = "https://example.com/some/page"
		)

		mocks := newServiceMocks(t)
		expectedShortName := textutils.RandomString(mocks.service.opts.shortNameDefaultLength)
		mocks.linkRepo.
			On("CreateOne", t.Context(), mock.MatchedBy(func(insert CreateLinkParams) bool {
				return insert.OriginalURL == originalURL && len(insert.ShortName) == mocks.service.opts.shortNameDefaultLength
			})).
			Return(Link{
				ID:          id,
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
		}, result)
	})

	t.Run("should return conflict when short name already exists", func(t *testing.T) {
		const (
			originalURL = "https://example.com/some/page"
			shortName   = "example"
		)

		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("CreateOne", t.Context(), CreateLinkParams{
				OriginalURL: originalURL,
				ShortName:   shortName,
			}).
			Return(Link{}, ErrRepositoryConflict).
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
			On("CreateOne", t.Context(), mock.MatchedBy(func(insert CreateLinkParams) bool {
				return insert.OriginalURL == originalURL && len(insert.ShortName) == mocks.service.opts.shortNameDefaultLength
			})).
			Return(Link{}, ErrRepositoryConflict).
			Once()

		expectedShortName := textutils.RandomString(mocks.service.opts.shortNameDefaultLength)

		mocks.linkRepo.
			On("CreateOne", t.Context(), mock.MatchedBy(func(insert CreateLinkParams) bool {
				return insert.OriginalURL == originalURL && len(insert.ShortName) == mocks.service.opts.shortNameDefaultLength
			})).
			Return(Link{
				ID:          id,
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
		}, result)

		mocks.linkRepo.AssertNumberOfCalls(t, "CreateOne", 2)
	})

	t.Run("should return error when short name generation attempts are exhausted", func(t *testing.T) {
		const originalURL = "https://example.com/some/page"

		mocks := newServiceMocks(t)

		mocks.linkRepo.
			On("CreateOne", t.Context(), mock.MatchedBy(func(insert CreateLinkParams) bool {
				return insert.OriginalURL == originalURL &&
					len(insert.ShortName) == mocks.service.opts.shortNameDefaultLength
			})).
			Return(Link{}, ErrRepositoryConflict).
			Times(mocks.service.opts.shortNameGenerationMaxAttempts)

		result, err := mocks.service.CreateLink(t.Context(), originalURL, "")

		require.ErrorIs(t, err, ErrShortNameGenerationExhausted)
		assert.Equal(t, Link{}, result)
		mocks.linkRepo.AssertNumberOfCalls(t, "CreateOne", mocks.service.opts.shortNameGenerationMaxAttempts)
	})

	t.Run("should return validation error for invalid short name", func(t *testing.T) {
		type subTest struct {
			name        string
			shortName   string
			originalURL string
		}

		subTests := []subTest{
			{
				name:        "short name too short",
				shortName:   "ab",
				originalURL: "https://example.com/updated",
			},
			{
				name:        "short name too long",
				shortName:   "abcdefghijklmnopqrstuvwxyz1234567",
				originalURL: "https://example.com/updated",
			},
		}
		for _, subTest := range subTests {
			t.Run(subTest.name, func(t *testing.T) {
				var validationErr *ValidationError

				mocks := newServiceMocks(t)

				result, err := mocks.service.CreateLink(t.Context(), subTest.originalURL, subTest.shortName)

				require.Error(t, err)
				require.ErrorAs(t, err, &validationErr)
				assert.Equal(t, Link{}, result)
			})
		}
	})
}

func TestService_GetLink(t *testing.T) {
	t.Run("should get link successfully", func(t *testing.T) {
		const (
			id          uint = 42
			originalURL      = "https://example.com/some/page"
			shortName        = "example"
		)

		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("GetByID", t.Context(), id).
			Return(Link{
				ID:          id,
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
		}, result)
	})

	t.Run("should return not found when link does not exist", func(t *testing.T) {
		const id uint = 42

		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("GetByID", t.Context(), id).
			Return(Link{}, ErrRepositoryNotFound).
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
			On("GetMany", t.Context(), LinkListOptions{
				ListOptions: ListOptions{
					Limit:     1,
					Offset:    0,
					SortOrder: "DESC",
				},
				SortBy:     LinkSortByID,
				ShortNames: []string{shortName},
			}).
			Return([]Link{
				{
					ID:          id,
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
		}, result)
	})

	t.Run("should return not found when link does not exist", func(t *testing.T) {
		const shortName = "missing"

		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("GetMany", t.Context(), LinkListOptions{
				ListOptions: ListOptions{
					Limit:     1,
					Offset:    0,
					SortOrder: "DESC",
				},
				SortBy:     LinkSortByID,
				ShortNames: []string{shortName},
			}).
			Return([]Link{}, nil).
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
			On("GetMany", t.Context(), LinkListOptions{
				ListOptions: ListOptions{
					Limit:     10,
					Offset:    0,
					SortOrder: "DESC",
				},
				SortBy: LinkSortByID,
			}).
			Return([]Link{
				{
					ID:          1,
					OriginalURL: "https://example.com/first",
					ShortName:   "first",
				},
				{
					ID:          2,
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
			},
			{
				ID:          2,
				OriginalURL: "https://example.com/second",
				ShortName:   "second",
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
			On("GetMany", t.Context(), LinkListOptions{
				ListOptions: ListOptions{
					Limit:     10,
					Offset:    10,
					SortOrder: "ASC",
				},
				SortBy:     LinkSortByShortName,
				ShortNames: []string{"first", "second"},
			}).
			Return([]Link{}, nil).
			Once()

		result, count, err := mocks.service.ListLinksWithCount(t.Context(), builder)

		require.NoError(t, err)
		assert.Empty(t, result)
		assert.Zero(t, count)
	})

	t.Run("should return validation error for invalid options", func(t *testing.T) {
		var validationErr *ValidationError

		mocks := newServiceMocks(t)
		builder := NewLinkListOptionsBuilder()
		builder.WithRange(-1, 10)

		result, count, err := mocks.service.ListLinksWithCount(t.Context(), builder)

		require.Error(t, err)
		require.ErrorAs(t, err, &validationErr)
		assert.Equal(t, "range", validationErr.Field)
		assert.Nil(t, result)
		assert.Zero(t, count)
		assert.Equal(t, "range start must be non-negative: -1", validationErr.Message)
		mocks.linkRepo.AssertNotCalled(t, "GetMany")
		mocks.linkRepo.AssertNotCalled(t, "Count")
	})

	t.Run("should fail when range contains more records than max limit", func(t *testing.T) {
		var validationErr *ValidationError

		mocks := newServiceMocks(t)
		builder := NewLinkListOptionsBuilder()
		builder.WithRange(0, 1001)

		result, count, err := mocks.service.ListLinksWithCount(t.Context(), builder)

		require.Error(t, err)
		require.ErrorAs(t, err, &validationErr)
		assert.Equal(t, "range", validationErr.Field)
		assert.Nil(t, result)
		assert.Zero(t, count)
		assert.Equal(t, "range must contain at most 1001 records, got 1002", validationErr.Message)
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
			On("UpdateByID", t.Context(), id, UpdateLinkParams{
				OriginalURL: originalURL,
				ShortName:   shortName,
			}).
			Return(Link{
				ID:          id,
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
			On("UpdateByID", t.Context(), id, UpdateLinkParams{
				OriginalURL: originalURL,
				ShortName:   shortName,
			}).
			Return(Link{}, ErrRepositoryNotFound).
			Once()

		result, err := mocks.service.UpdateLink(t.Context(), id, originalURL, shortName)

		require.Error(t, err)
		assert.Equal(t, Link{}, result)
	})

	t.Run("should return validation error for invalid short name", func(t *testing.T) {
		type subTest struct {
			name        string
			id          uint
			shortName   string
			originalURL string
		}

		subTests := []subTest{
			{
				name:        "short name too short",
				shortName:   "ab",
				id:          1,
				originalURL: "https://example.com/updated",
			},
			{
				name:        "short name too long",
				shortName:   "abcdefghijklmnopqrstuvwxyz1234567",
				id:          1,
				originalURL: "https://example.com/updated",
			},
		}
		for _, subTest := range subTests {
			t.Run(subTest.name, func(t *testing.T) {
				mocks := newServiceMocks(t)

				var validationErr *ValidationError

				result, err := mocks.service.UpdateLink(t.Context(), subTest.id, subTest.originalURL, subTest.shortName)

				require.Error(t, err)
				require.ErrorAs(t, err, &validationErr)
				assert.Equal(t, Link{}, result)
			})
		}
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
			Return(ErrRepositoryNotFound).
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
			On("CreateOne", t.Context(), CreateLinkVisitParams{
				LinkID:    linkID,
				IP:        ip,
				UserAgent: userAgent,
				Referrer:  referrer,
				Status:    status,
			}).
			Return(LinkVisit{
				ID:        1,
				CreatedAt: createdAt,
				UpdatedAt: updatedAt,
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
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
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
			On("CreateOne", t.Context(), CreateLinkVisitParams{
				LinkID:    linkID,
				IP:        ip,
				UserAgent: userAgent,
				Referrer:  referrer,
				Status:    status,
			}).
			Return(LinkVisit{}, ErrRepositoryConflict).
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
			Return([]LinkVisit{
				{
					ID:        1,
					LinkID:    42,
					IP:        "127.0.0.1",
					UserAgent: "Mozilla/5.0",
					Referrer:  "https://example.com",
					Status:    302,
				},
				{
					ID:        2,
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
			On("GetMany", t.Context(), mock.AnythingOfType("LinkVisitListOptions")).
			Return([]LinkVisit{}, nil).
			Once()

		result, count, err := mocks.service.ListLinkVisitsWithCount(t.Context(), nil)

		require.NoError(t, err)
		assert.Empty(t, result)
		assert.Zero(t, count)
	})

	t.Run("should return validation error for invalid options", func(t *testing.T) {
		var validationErr *ValidationError

		mocks := newServiceMocks(t)
		builder := NewLinkVisitListOptionsBuilder()
		builder.WithRange(-1, 10)

		result, count, err := mocks.service.ListLinkVisitsWithCount(t.Context(), builder)

		require.Error(t, err)
		require.ErrorAs(t, err, &validationErr)
		assert.Equal(t, "range", validationErr.Field)
		assert.Nil(t, result)
		assert.Zero(t, count)
		assert.Equal(t, "range start must be non-negative: -1", validationErr.Message)
		mocks.linkRepo.AssertNotCalled(t, "GetMany")
		mocks.linkRepo.AssertNotCalled(t, "Count")
	})

	t.Run("should fail when range contains more records than max limit", func(t *testing.T) {
		var validationErr *ValidationError

		mocks := newServiceMocks(t)
		builder := NewLinkVisitListOptionsBuilder()
		builder.WithRange(0, 999)

		result, count, err := mocks.service.ListLinkVisitsWithCount(t.Context(), builder)

		require.Error(t, err)
		require.ErrorAs(t, err, &validationErr)
		assert.Equal(t, "range", validationErr.Field)
		assert.Nil(t, result)
		assert.Zero(t, count)
		assert.Equal(t, "range must contain at most 100 records, got 1000", validationErr.Message)
		mocks.linkRepo.AssertNotCalled(t, "GetMany")
		mocks.linkRepo.AssertNotCalled(t, "Count")
	})

	t.Run("should return error when getting link visits fails", func(t *testing.T) {
		mocks := newServiceMocks(t)
		builder := NewLinkVisitListOptionsBuilder()

		mocks.linkVisitRepo.
			On("GetMany", t.Context(), builder.build()).
			Return([]LinkVisit{}, ErrRepositoryNotFound).
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
			Return([]LinkVisit{}, nil).
			Once()
		mocks.linkVisitRepo.
			On("Count", t.Context()).
			Return(0, ErrRepositoryNotFound).
			Once()

		result, count, err := mocks.service.ListLinkVisitsWithCount(t.Context(), builder)

		require.Error(t, err)
		assert.Nil(t, result)
		assert.Zero(t, count)
	})
}
