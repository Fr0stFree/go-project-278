package link

import (
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"shortener/internal/services/shortener"
)

func newRepositoryMock(t *testing.T) (*Repository, sqlmock.Sqlmock) {
	t.Helper()

	sqlDB, sqlMock, err := sqlmock.New()
	require.NoError(t, err)

	t.Cleanup(func() {
		sqlMock.ExpectClose()
		require.NoError(t, sqlDB.Close())
		require.NoError(t, sqlMock.ExpectationsWereMet())
	})

	gormDB, err := gorm.Open(
		gormpostgres.New(gormpostgres.Config{Conn: sqlDB}),
		&gorm.Config{TranslateError: true},
	)
	require.NoError(t, err)

	return NewRepository(gormDB), sqlMock
}

func TestRepository_CreateOne(t *testing.T) {
	t.Run("should create link successfully", func(t *testing.T) {
		repository, sqlMock := newRepositoryMock(t)
		sqlMock.ExpectBegin()
		sqlMock.ExpectQuery(`INSERT INTO "shortened_links"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		sqlMock.ExpectCommit()

		record, err := repository.CreateOne(t.Context(), shortener.CreateLinkParams{
			OriginalURL: "https://example.com",
			ShortName:   "abc123",
		})

		require.NoError(t, err)
		assert.Equal(t, uint(1), record.ID)
		assert.Equal(t, "https://example.com", record.OriginalURL)
		assert.Equal(t, "abc123", record.ShortName)
	})
}

func TestRepository_GetByID(t *testing.T) {
	t.Run("should get link by ID successfully", func(t *testing.T) {
		repository, sqlMock := newRepositoryMock(t)
		sqlMock.ExpectQuery(`SELECT \* FROM "shortened_links" WHERE "shortened_links"\."id" = \$1`).
			WillReturnRows(sqlmock.NewRows([]string{"id", "original_url", "short_name"}).
				AddRow(1, "https://example.com", "abc123"))

		record, err := repository.GetByID(t.Context(), 1)

		require.NoError(t, err)
		assert.Equal(t, shortener.Link{
			ID:          1,
			OriginalURL: "https://example.com",
			ShortName:   "abc123",
		}, record)
	})
}

func TestRepository_GetMany(t *testing.T) {
	t.Run("should get links by short names successfully", func(t *testing.T) {
		repository, sqlMock := newRepositoryMock(t)

		sqlMock.
			ExpectQuery(`SELECT \* FROM "shortened_links" WHERE short_name IN`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "original_url", "short_name"}).
					AddRow(1, "https://example.com", "abc123").
					AddRow(2, "https://example.org", "def456"),
			)

		options := shortener.LinkListOptions{
			ListOptions: shortener.ListOptions{
				Limit:     10,
				SortOrder: "asc",
			},
			SortBy:     shortener.LinkSortByID,
			ShortNames: []string{"abc123", "def456"},
		}

		records, err := repository.GetMany(t.Context(), options)

		require.NoError(t, err)
		require.Len(t, records, 2)
		assert.Equal(t, shortener.Link{
			ID:          1,
			OriginalURL: "https://example.com",
			ShortName:   "abc123",
		}, records[0])
		assert.Equal(t, shortener.Link{
			ID:          2,
			OriginalURL: "https://example.org",
			ShortName:   "def456",
		}, records[1])
	})
}

func TestRepository_Count(t *testing.T) {
	t.Run("should count links successfully", func(t *testing.T) {
		repository, sqlMock := newRepositoryMock(t)
		sqlMock.
			ExpectQuery(`SELECT count\(\*\) FROM "shortened_links"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(42))

		count, err := repository.Count(t.Context())

		require.NoError(t, err)
		assert.Equal(t, 42, count)
	})
}

func TestRepository_UpdateByID(t *testing.T) {
	t.Run("should update link successfully", func(t *testing.T) {
		repository, sqlMock := newRepositoryMock(t)
		sqlMock.ExpectBegin()
		sqlMock.
			ExpectQuery(`UPDATE "shortened_links" SET .* WHERE id = \$[0-9]+ AND "shortened_links"."deleted_at" IS NULL RETURNING \*`).
			WithArgs(sqlmock.AnyArg(), "https://example.org", "def456", uint(1)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "original_url", "short_name"}).AddRow(1, "https://example.org", "def456"))
		sqlMock.ExpectCommit()

		result, err := repository.UpdateByID(t.Context(), 1, shortener.UpdateLinkParams{
			OriginalURL: "https://example.org",
			ShortName:   "def456",
		})

		require.NoError(t, err)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "https://example.org", result.OriginalURL)
		assert.Equal(t, "def456", result.ShortName)
	})

	t.Run("should return not found when link does not exist", func(t *testing.T) {
		repository, sqlMock := newRepositoryMock(t)
		sqlMock.ExpectBegin()
		sqlMock.
			ExpectQuery(`UPDATE "shortened_links" SET .* WHERE id = \$[0-9]+ AND "shortened_links"."deleted_at" IS NULL RETURNING \*`).
			WithArgs(sqlmock.AnyArg(), "https://example.org", "def456", uint(999)).
			WillReturnRows(sqlmock.NewRows([]string{"id", "original_url", "short_name"}))
		sqlMock.ExpectCommit()

		result, err := repository.UpdateByID(t.Context(), 999, shortener.UpdateLinkParams{
			OriginalURL: "https://example.org",
			ShortName:   "def456",
		})

		require.ErrorIs(t, err, shortener.ErrRepositoryNotFound)
		assert.Equal(t, shortener.Link{}, result)
	})

	t.Run("should propagate unexpected database error", func(t *testing.T) {
		repository, sqlMock := newRepositoryMock(t)
		dbErr := errors.New("database unavailable")

		sqlMock.ExpectBegin()
		sqlMock.ExpectQuery(`UPDATE "shortened_links" SET .* WHERE id = \$[0-9]+ AND "shortened_links"."deleted_at" IS NULL RETURNING \*`).
			WithArgs(sqlmock.AnyArg(), "https://example.org", "def456", uint(1)).
			WillReturnError(dbErr)
		sqlMock.ExpectRollback()

		result, err := repository.UpdateByID(t.Context(), 1, shortener.UpdateLinkParams{
			OriginalURL: "https://example.org",
			ShortName:   "def456",
		})

		require.ErrorIs(t, err, dbErr)
		assert.Equal(t, shortener.Link{}, result)
	})
}

func TestRepository_DeleteByID(t *testing.T) {
	t.Run("should delete link successfully", func(t *testing.T) {
		repository, sqlMock := newRepositoryMock(t)
		sqlMock.ExpectBegin()
		sqlMock.
			ExpectExec(`UPDATE "shortened_links" SET "deleted_at"`).
			WithArgs(sqlmock.AnyArg(), 1).
			WillReturnResult(sqlmock.NewResult(0, 1))
		sqlMock.ExpectCommit()

		err := repository.DeleteByID(t.Context(), 1)

		require.NoError(t, err)
	})
}

func TestToSortOrder(t *testing.T) {
	type subTest struct {
		name      string
		field     shortener.LinkSortField
		direction shortener.SortDirection
		expected  string
	}

	subTests := []subTest{
		{name: "id", field: shortener.LinkSortByID, direction: shortener.SortAscending, expected: "id ASC"},
		{name: "original url", field: shortener.LinkSortByOriginalURL, direction: shortener.SortDescending, expected: "original_url DESC, id DESC"},
		{name: "short name", field: shortener.LinkSortByShortName, direction: shortener.SortAscending, expected: "short_name ASC, id ASC"},
		{name: "created at", field: shortener.LinkSortByCreatedAt, direction: shortener.SortDescending, expected: "created_at DESC, id DESC"},
	}

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {
			assert.Equal(t, subTest.expected, toSortOrder(subTest.field, subTest.direction))
		})
	}
}
