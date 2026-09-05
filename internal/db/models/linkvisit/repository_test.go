package linkvisit

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	gormpostgres "gorm.io/driver/postgres"
	"gorm.io/gorm"

	"shortener/internal/db"
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
	t.Run("should create visit successfully", func(t *testing.T) {
		repository, sqlMock := newRepositoryMock(t)
		sqlMock.ExpectBegin()
		sqlMock.
			ExpectQuery(`INSERT INTO "shortened_link_visits"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id"}).AddRow(1),
			)
		sqlMock.ExpectCommit()

		record, err := repository.CreateOne(Insert{
			LinkID:    1,
			IP:        "127.0.0.1",
			UserAgent: "Mozilla/5.0",
			Status:    200,
			Referrer:  "https://example.com",
		})

		require.NoError(t, err)
		assert.Equal(t, uint(1), record.ID)
		assert.Equal(t, uint(1), record.LinkID)
		assert.Equal(t, "127.0.0.1", record.IP)
		assert.Equal(t, "Mozilla/5.0", record.UserAgent)
		assert.Equal(t, uint(200), record.Status)
		assert.Equal(t, "https://example.com", record.Referrer)
	})
}

func TestRepository_GetMany(t *testing.T) {
	t.Run("should get visits successfully", func(t *testing.T) {
		repository, sqlMock := newRepositoryMock(t)
		opts := ListOptions{
			db.ListOptions{Limit: 10, Offset: 0},
			Filters{LinkIDs: []uint{1}},
		}

		sqlMock.
			ExpectQuery(`SELECT \* FROM "shortened_link_visits"`).
			WillReturnRows(
				sqlmock.NewRows([]string{"id", "link_id", "ip", "user_agent", "status", "referrer"}).
					AddRow(1, 1, "127.0.0.1", "Mozilla/5.0", 200, "https://example.com").
					AddRow(2, 1, "192.168.1.1", "Chrome/89.0", 200, "https://google.com"),
			)

		records, err := repository.GetMany(opts)

		require.NoError(t, err)
		assert.Len(t, records, 2)
		assert.Equal(t, uint(1), records[0].ID)
		assert.Equal(t, uint(1), records[0].LinkID)
		assert.Equal(t, "127.0.0.1", records[0].IP)
		assert.Equal(t, "Mozilla/5.0", records[0].UserAgent)
		assert.Equal(t, uint(200), records[0].Status)
		assert.Equal(t, "https://example.com", records[0].Referrer)
	})
}

func TestRepository_Count(t *testing.T) {
	t.Run("should count visits successfully", func(t *testing.T) {
		repository, sqlMock := newRepositoryMock(t)
		sqlMock.ExpectQuery(`SELECT count\(\*\) FROM "shortened_link_visits"`).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(42))

		count, err := repository.Count()

		require.NoError(t, err)
		assert.Equal(t, 42, count)
	})
}
