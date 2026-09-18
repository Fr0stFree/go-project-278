//go:build integration

package integration

import (
	"bytes"
	"encoding/json"
	"net/http"
	"shortener/internal/services/shortener"
	"strconv"
	"testing"

	"context"
	"database/sql"
	"fmt"
	"shortener/internal/app"
	"shortener/internal/config"
	"shortener/internal/db"
	"shortener/internal/httpserver"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

const (
	dbImage       = "postgres:18-alpine"
	dbUser        = "shortener"
	dbPassword    = "shortener"
	dbName        = "shortener"
	migrationsDir = "../../../db/migrations"
)

type linkResponse struct {
	ID          uint   `json:"id"`
	OriginalURL string `json:"original_url"`
	ShortName   string `json:"short_name"`
	ShortURL    string `json:"short_url"`
}

func applyMigrations(t *testing.T, databaseURL string) {
	t.Helper()
	database, err := sql.Open("pgx", databaseURL)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = database.Close()
	})
	require.NoError(t, database.PingContext(t.Context()))
	require.NoError(t, goose.SetDialect("postgres"))
	require.NoError(t, goose.Up(database, migrationsDir))
}

func runDB(t *testing.T) string {
	pg, err := postgres.Run(
		t.Context(),
		dbImage,
		postgres.WithDatabase(dbName),
		postgres.WithUsername(dbUser),
		postgres.WithPassword(dbPassword),
		postgres.BasicWaitStrategies(),
	)
	require.NoError(t, err)

	t.Cleanup(func() {
		_ = pg.Terminate(context.Background())
	})
	dbURL, err := pg.ConnectionString(t.Context(), "sslmode=disable")
	require.NoError(t, err)

	return dbURL
}

func startApp(t *testing.T, databaseURL string) string {
	t.Helper()
	t.Setenv("DATABASE_URL", databaseURL)

	cfg, err := config.New()
	require.NoError(t, err)
	cfg.HTTP.Sentry.IsEnabled = false

	database, err := db.New(&cfg.Database)
	require.NoError(t, err)

	service := shortener.NewService(
		database.Links,
		database.LinkVisits,
	)

	server := httpserver.New(service, &cfg.HTTP, cfg.App.BaseURL)
	application := app.New(server, database, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 1)

	go func() {
		errCh <- application.Run(ctx)
	}()

	t.Cleanup(func() {
		cancel()

		select {
		case err := <-errCh:
			require.NoError(t, err)
		case <-time.After(5 * time.Second):
			t.Error("application did not stop within 5 seconds")
		}
	})

	baseURL := fmt.Sprintf("http://localhost:%d", cfg.HTTP.Port)

	require.Eventually(t, func() bool {
		resp, err := http.Get(baseURL + "/ping")
		if err != nil {
			return false
		}
		defer resp.Body.Close()

		return resp.StatusCode == http.StatusOK
	}, 5*time.Second, 50*time.Millisecond)

	return baseURL
}

func initApp(t *testing.T) string {
	databaseURL := runDB(t)
	applyMigrations(t, databaseURL)
	appURL := startApp(t, databaseURL)

	return appURL
}

func TestLinkCreation(t *testing.T) {
	appURL := initApp(t)
	url := appURL + "/api/links"
	body := []byte(`{
		"original_url": "https://example.com",
		"short_name": "example"
	}`)

	resp, err := http.Post(url, "application/json", bytes.NewReader(body))
	require.NoError(t, err)
	defer resp.Body.Close()

	var link linkResponse
	require.Equal(t, http.StatusCreated, resp.StatusCode)
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&link))
	require.Equal(t, "https://example.com", link.OriginalURL)
	require.Equal(t, "example", link.ShortName)
	require.NotEmpty(t, link.ShortURL)

	resp, err = http.Get(url + "/" + strconv.Itoa(int(link.ID)))
	require.NoError(t, err)
	defer resp.Body.Close()

	var retrievedLink linkResponse
	require.Equal(t, http.StatusOK, resp.StatusCode)
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&retrievedLink))
	require.Equal(t, link.ID, retrievedLink.ID)
	require.Equal(t, link.OriginalURL, retrievedLink.OriginalURL)
	require.Equal(t, link.ShortName, retrievedLink.ShortName)
	require.Equal(t, link.ShortURL, retrievedLink.ShortURL)
}
