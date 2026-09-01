package shortener

import (
	"shortener/internal/config"
	"shortener/internal/database/storage"
	"shortener/internal/database/storage/link"
	"shortener/internal/database/storage/linkvisit"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

type mockLinkRepository struct {
	mock.Mock
}

func (m *mockLinkRepository) CreateOne(insert link.Insert) (link.Record, error) {
	args := m.Called(insert)
	return args.Get(0).(link.Record), args.Error(1)
}

func (m *mockLinkRepository) GetByID(id uint) (link.Record, error) {
	args := m.Called(id)
	return args.Get(0).(link.Record), args.Error(1)
}

func (m *mockLinkRepository) GetMany(options link.ListOptions) ([]link.Record, error) {
	args := m.Called(options)
	return args.Get(0).([]link.Record), args.Error(1)
}

func (m *mockLinkRepository) Count() (int, error) {
	args := m.Called()
	return args.Int(0), args.Error(1)
}

func (m *mockLinkRepository) UpdateByID(id uint, update link.Update) (link.Record, error) {
	args := m.Called(id, update)
	return args.Get(0).(link.Record), args.Error(1)
}

func (m *mockLinkRepository) DeleteByID(id uint) error {
	args := m.Called(id)
	return args.Error(0)
}

type mockLinkVisitRepository struct {
	mock.Mock
}

func (m *mockLinkVisitRepository) CreateOne(insert linkvisit.Insert) (linkvisit.Record, error) {
	args := m.Called(insert)
	return args.Get(0).(linkvisit.Record), args.Error(1)
}

func (m *mockLinkVisitRepository) GetMany(options linkvisit.ListOptions) ([]linkvisit.Record, error) {
	args := m.Called(options)
	return args.Get(0).([]linkvisit.Record), args.Error(1)
}

func (m *mockLinkVisitRepository) Count() (int, error) {
	args := m.Called()
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
			On("CreateOne", link.Insert{
				OriginalURL: originalURL,
				ShortName:   shortName,
			}).
			Return(link.Record{
				Model:       gorm.Model{ID: id},
				OriginalURL: originalURL,
				ShortName:   shortName,
			}, nil).
			Once()

		result, err := mocks.service.CreateLink(originalURL, shortName)

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
		expectedShortName := toHashString(originalURL, 6)
		mocks.linkRepo.
			On("CreateOne", link.Insert{
				OriginalURL: originalURL,
				ShortName:   expectedShortName,
			}).
			Return(link.Record{
				Model:       gorm.Model{ID: id},
				OriginalURL: originalURL,
				ShortName:   expectedShortName,
			}, nil).
			Once()

		result, err := mocks.service.CreateLink(originalURL, "")

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
			On("CreateOne", link.Insert{
				OriginalURL: originalURL,
				ShortName:   shortName,
			}).
			Return(link.Record{}, storage.ErrObjectAlreadyExists).
			Once()

		result, err := mocks.service.CreateLink(originalURL, shortName)

		require.Error(t, err)
		assert.Equal(t, Link{}, result)
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
			On("GetByID", id).
			Return(link.Record{
				Model:       gorm.Model{ID: id},
				OriginalURL: originalURL,
				ShortName:   shortName,
			}, nil).
			Once()

		result, err := mocks.service.GetLink(id)

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
			On("GetByID", id).
			Return(link.Record{}, storage.ErrObjectDoesNotExist).
			Once()

		result, err := mocks.service.GetLink(id)

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
			On("GetMany", link.ListOptions{
				ListOptions: storage.ListOptions{
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

		result, err := mocks.service.GetRedirectLink(shortName)

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
			On("GetMany", link.ListOptions{
				ListOptions: storage.ListOptions{
					Limit:     1,
					Offset:    0,
					SortBy:    "id",
					SortOrder: "DESC",
				},
				Filters: link.Filters{ShortNames: []string{shortName}},
			}).
			Return([]link.Record{}, nil).
			Once()

		result, err := mocks.service.GetRedirectLink(shortName)

		require.Error(t, err)
		assert.Equal(t, Link{}, result)
	})
}

func TestService_ListLinksWithCount(t *testing.T) {
	t.Run("should list links with count successfully", func(t *testing.T) {
		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("Count").
			Return(42, nil).
			Once()
		mocks.linkRepo.
			On("GetMany", link.ListOptions{
				ListOptions: storage.ListOptions{
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

		result, count, err := mocks.service.ListLinksWithCount(nil)

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
			On("Count").
			Return(0, nil).
			Once()
		mocks.linkRepo.
			On("GetMany", link.ListOptions{
				ListOptions: storage.ListOptions{
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

		result, count, err := mocks.service.ListLinksWithCount(builder)

		require.NoError(t, err)
		assert.Empty(t, result)
		assert.Zero(t, count)
	})

	t.Run("should return validation error for invalid options", func(t *testing.T) {
		mocks := newServiceMocks(t)
		builder := NewLinkListOptionsBuilder()
		builder.WithRange(-1, 10)

		result, count, err := mocks.service.ListLinksWithCount(builder)

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
			On("UpdateByID", id, link.Update{
				OriginalURL: originalURL,
				ShortName:   shortName,
			}).
			Return(link.Record{
				Model:       gorm.Model{ID: id},
				OriginalURL: originalURL,
				ShortName:   shortName,
			}, nil).
			Once()

		result, err := mocks.service.UpdateLink(id, originalURL, shortName)

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
			On("UpdateByID", id, link.Update{
				OriginalURL: originalURL,
				ShortName:   shortName,
			}).
			Return(link.Record{}, storage.ErrObjectDoesNotExist).
			Once()

		result, err := mocks.service.UpdateLink(id, originalURL, shortName)

		require.Error(t, err)
		assert.Equal(t, Link{}, result)
	})
}

func TestService_DeleteLink(t *testing.T) {
	t.Run("should delete link successfully", func(t *testing.T) {
		const id uint = 42
		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("DeleteByID", id).
			Return(nil).
			Once()

		err := mocks.service.DeleteLink(id)

		require.NoError(t, err)
	})

	t.Run("should return not found when link does not exist", func(t *testing.T) {
		const id uint = 42
		mocks := newServiceMocks(t)
		mocks.linkRepo.
			On("DeleteByID", id).
			Return(storage.ErrObjectDoesNotExist).
			Once()

		err := mocks.service.DeleteLink(id)

		require.Error(t, err)
	})

}