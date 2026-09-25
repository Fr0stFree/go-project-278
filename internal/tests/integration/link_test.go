//go:build integration

package integration

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/postgres"

	"shortener/internal/app"
	"shortener/internal/config"
	"shortener/internal/db"
	"shortener/internal/httpserver"
	"shortener/internal/services/shortener"
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

type linkVisitResponse struct {
	ID        uint   `json:"id"`
	LinkID    uint   `json:"link_id"`
	CreatedAt string `json:"created_at"`
	IP        string `json:"ip"`
	UserAgent string `json:"user_agent"`
	Status    uint   `json:"status"`
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
	t.Helper()
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

	service := shortener.NewService(database.Links, database.LinkVisits)
	server := httpserver.New(service, database, &cfg.HTTP, cfg.App.BaseURL)
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
		resp, err := http.Get(baseURL + "/health")
		if err != nil {
			return false
		}
		defer func() {
			err := resp.Body.Close()
			require.NoError(t, err)
		}()

		return resp.StatusCode == http.StatusOK
	}, 5*time.Second, 50*time.Millisecond)

	return baseURL
}

func setupApp(t *testing.T) string {
	t.Helper()
	databaseURL := runDB(t)
	applyMigrations(t, databaseURL)
	appURL := startApp(t, databaseURL)

	return appURL
}

type testClient struct {
	baseURL   string
	client    *http.Client
	userAgent string
}

func newTestClient(baseURL string, userAgent string) *testClient {
	return &testClient{
		baseURL:   baseURL,
		client:    &http.Client{},
		userAgent: userAgent,
	}
}

func (c *testClient) request(
	t *testing.T,
	method string,
	path string,
	body any,
) *http.Response {
	t.Helper()

	var reader io.Reader

	if body != nil {
		data, err := json.Marshal(body)
		require.NoError(t, err)

		reader = bytes.NewReader(data)
	}

	req, err := http.NewRequest(method, c.baseURL+path, reader)
	req.Header.Set("User-Agent", c.userAgent)
	require.NoError(t, err)

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	resp, err := c.client.Do(req)
	require.NoError(t, err)

	return resp
}

func decodeJSON[T any](t *testing.T, resp *http.Response) T {
	t.Helper()

	var result T
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&result))

	t.Cleanup(func() {
		err := resp.Body.Close()
		require.NoError(t, err)
	})

	return result
}

func TestLinkRedirectFlow(t *testing.T) {
	appURL := setupApp(t)
	expectedUserAgent := "integration-test"
	client := newTestClient(appURL, expectedUserAgent)

	var (
		expectedID          uint = 1
		expectedShortName        = "example"
		expectedOriginalURL      = fmt.Sprintf("%s/ping", appURL)
		listLinksPath            = "/api/links"
		listVisitsPath           = "/api/link_visits"
		redirectPath             = fmt.Sprintf("/r/%s", expectedShortName)
		getLinkPath              = fmt.Sprintf("%s/%d", listLinksPath, expectedID)
	)

	t.Run("should create link", func(t *testing.T) {
		var respBody linkResponse

		resp := client.request(t, http.MethodPost, listLinksPath, map[string]string{
			"original_url": expectedOriginalURL,
			"short_name":   expectedShortName,
		})

		require.Equal(t, http.StatusCreated, resp.StatusCode)

		respBody = decodeJSON[linkResponse](t, resp)
		require.NotZero(t, respBody.ID)
		require.Equal(t, expectedOriginalURL, respBody.OriginalURL)
		require.Equal(t, expectedShortName, respBody.ShortName)
		require.NotEmpty(t, respBody.ShortURL)
	})

	t.Run("should retrieve created link", func(t *testing.T) {
		var respBody linkResponse

		resp := client.request(t, http.MethodGet, getLinkPath, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)

		respBody = decodeJSON[linkResponse](t, resp)
		require.NotZero(t, respBody.ID)
		require.Equal(t, expectedOriginalURL, respBody.OriginalURL)
		require.Equal(t, expectedShortName, respBody.ShortName)
		require.NotEmpty(t, respBody.ShortURL)
	})

	t.Run("should follow redirect", func(t *testing.T) {
		resp := client.request(t, http.MethodGet, redirectPath, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, expectedOriginalURL, resp.Request.URL.String())
		require.NotNil(t, resp.Request.Response)
		require.Equal(t, http.StatusFound, resp.Request.Response.StatusCode)
		require.Equal(t, expectedOriginalURL, resp.Request.Response.Header.Get("Location"))
	})

	t.Run("should find redirect", func(t *testing.T) {
		var respBody []linkVisitResponse

		resp := client.request(t, http.MethodGet, listVisitsPath, nil)

		require.Equal(t, http.StatusOK, resp.StatusCode)
		require.Equal(t, "link_visits 0-0/1", resp.Header.Get("Content-Range"))

		respBody = decodeJSON[[]linkVisitResponse](t, resp)
		require.Len(t, respBody, 1)

		visit := respBody[0]
		require.NotZero(t, visit.ID)
		require.Equal(t, expectedID, visit.LinkID)
		require.Equal(t, uint(http.StatusFound), visit.Status)
		require.Equal(t, expectedUserAgent, visit.UserAgent)
		require.NotEmpty(t, visit.IP)
		_, err := time.Parse(time.RFC3339, visit.CreatedAt)
		require.NoError(t, err)
	})
}
