package link

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"shortener/internal/db"
	"shortener/internal/db/models"
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

	database := &db.DataBase{DB: gormDB}

	return NewRepository(database), sqlMock
}

func TestRepository_CreateOne(t *testing.T) {
	t.Run("should create link successfully", func(t *testing.T) {
		repository, sqlMock := newRepositoryMock(t)
		sqlMock.ExpectBegin()
		sqlMock.ExpectQuery(`INSERT INTO "shortened_links"`).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		sqlMock.ExpectCommit()

		record, err := repository.CreateOne(Insert{
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

		record, err := repository.GetByID(1)

		require.NoError(t, err)
		assert.Equal(t, record, Record{
			Model:       gorm.Model{ID: 1},
			OriginalURL: "https://example.com",
			ShortName:   "abc123",
		})
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

		options := ListOptions{
			ListOptions: models.ListOptions{
				Limit:     10,
				SortBy:    "id",
				SortOrder: "asc",
			},
			Filters: Filters{ShortNames: []string{"abc123", "def456"}},
		}

		records, err := repository.GetMany(options)

		require.NoError(t, err)
		require.Len(t, records, 2)
		assert.Equal(t, records[0], Record{
			Model:       gorm.Model{ID: 1},
			OriginalURL: "https://example.com",
			ShortName:   "abc123",
		})
		assert.Equal(t, records[1], Record{
			Model:       gorm.Model{ID: 2},
			OriginalURL: "https://example.org",
			ShortName:   "def456",
		})
	})
}

func TestRepository_Count(t *testing.T) {
	t.Run("should count links successfully", func(t *testing.T) {
		repository, sqlMock := newRepositoryMock(t)
		sqlMock.
			ExpectQuery(`SELECT count\(\*\) FROM "shortened_links"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(42))

		count, err := repository.Count()

		require.NoError(t, err)
		assert.Equal(t, count, 42)
	})
}

func TestRepository_UpdateByID(t *testing.T) {
	t.Run("should update link successfully", func(t *testing.T) {
		repository, sqlMock := newRepositoryMock(t)
		sqlMock.
			ExpectQuery(`SELECT \* FROM "shortened_links"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "original_url", "short_name"}).
					AddRow(1, "https://example.com", "abc123"),
			)
		sqlMock.ExpectBegin()
		sqlMock.
			ExpectExec(`UPDATE "shortened_links" SET`).
			WillReturnResult(sqlmock.NewResult(0, 1))
		sqlMock.ExpectCommit()

		result, err := repository.UpdateByID(1, Update{
			OriginalURL: "https://example.org",
			ShortName:   "def456",
		})

		require.NoError(t, err)
		assert.Equal(t, uint(1), result.ID)
		assert.Equal(t, "https://example.org", result.OriginalURL)
		assert.Equal(t, "def456", result.ShortName)
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

		err := repository.DeleteByID(1)

		require.NoError(t, err)
	})
}
