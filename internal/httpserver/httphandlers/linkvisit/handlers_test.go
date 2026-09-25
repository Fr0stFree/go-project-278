package linkvisit

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"shortener/internal/services/shortener"
)

type mockService struct {
	mock.Mock
}

func (m *mockService) ListLinkVisitsWithCount(ctx context.Context, optsBuilder *shortener.LinkVisitListOptionsBuilder) ([]shortener.LinkVisit, int, error) {
	args := m.Called(ctx, optsBuilder)

	return args.Get(0).([]shortener.LinkVisit), args.Int(1), args.Error(2)
}

type handlerMocks struct {
	shortener *mockService
	handler   *handler
}

func newHandlerMocks(t *testing.T) handlerMocks {
	t.Helper()

	shortener := new(mockService)
	router := gin.New()

	handler := &handler{shortener}
	RegisterRoutes(shortener, router)

	t.Cleanup(func() {
		shortener.AssertExpectations(t)
	})

	return handlerMocks{
		handler:   handler,
		shortener: shortener,
	}
}

func newRouter(t *testing.T, service Service) *gin.Engine {
	t.Helper()

	router := gin.New()
	RegisterRoutes(service, router)

	return router
}

func TestHandler_list(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should list link visits with count successfully", func(t *testing.T) {
		mocks := newHandlerMocks(t)
		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/link_visits", nil)

		mocks.shortener.
			On("ListLinkVisitsWithCount", mock.Anything, mock.Anything).
			Return([]shortener.LinkVisit{
				{
					ID:        1,
					LinkID:    1,
					CreatedAt: time.Date(2026, 9, 4, 12, 0, 0, 0, time.UTC),
					IP:        "127.0.0.1",
					UserAgent: "Mozilla/5.0",
					Status:    http.StatusFound,
				},
			}, 1, nil).
			Once()

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		assert.JSONEq(t, `[
			{
				"id": 1,
				"link_id": 1,
				"created_at": "2026-09-04T12:00:00Z",
				"ip": "127.0.0.1",
				"user_agent": "Mozilla/5.0",
				"status": 302
			}
		]`, recorder.Body.String())
	})

	t.Run("should return empty content range when requested page is empty", func(t *testing.T) {
		mocks := newHandlerMocks(t)
		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/link_visits?range=%5B10%2C19%5D", nil)

		mocks.shortener.
			On("ListLinkVisitsWithCount", mock.Anything, mock.Anything).
			Return([]shortener.LinkVisit{}, 42, nil).
			Once()

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "link_visits */42", recorder.Header().Get("Content-Range"))
		assert.JSONEq(t, `[]`, recorder.Body.String())
	})

	t.Run("should build content range from actual number of returned visits", func(t *testing.T) {
		mocks := newHandlerMocks(t)
		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/link_visits?range=%5B10%2C19%5D", nil)

		mocks.shortener.
			On("ListLinkVisitsWithCount", mock.Anything, mock.Anything).
			Return([]shortener.LinkVisit{
				{ID: 11, LinkID: 1},
				{ID: 12, LinkID: 1},
				{ID: 13, LinkID: 1},
				{ID: 14, LinkID: 1},
			}, 14, nil).
			Once()

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "link_visits 10-13/14", recorder.Header().Get("Content-Range"))
	})

	t.Run("should return bad request for invalid range", func(t *testing.T) {
		mocks := newHandlerMocks(t)
		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/link_visits?range=invalid", nil)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	})

	t.Run("should handle sort parameter", func(t *testing.T) {
		mocks := newHandlerMocks(t)
		mocks.shortener.
			On("ListLinkVisitsWithCount", mock.Anything, mock.Anything).
			Return([]shortener.LinkVisit{}, 0, nil).
			Once()

		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, `/api/link_visits?sort=["status","asc"]`, nil)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		assert.JSONEq(t, `[]`, recorder.Body.String())
	})

	t.Run("should reject invalid sort format", func(t *testing.T) {
		mocks := newHandlerMocks(t)
		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/link_visits?sort=invalid", nil)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
		assert.JSONEq(t, `{"errors":{"sort":"invalid sort format: invalid"}}`, recorder.Body.String())
	})

	t.Run("should reject unsupported sort field", func(t *testing.T) {
		service := shortener.NewService(nil, nil)
		router := newRouter(t, service)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, `/api/link_visits?sort=["referrer","ASC"]`, nil)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
		assert.JSONEq(t, `{"errors":{"sort":"unsupported sort field: \"referrer\""}}`, recorder.Body.String())
	})

	t.Run("should reject unsupported sort direction", func(t *testing.T) {
		service := shortener.NewService(nil, nil)
		router := newRouter(t, service)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, `/api/link_visits?sort=["status","sideways"]`, nil)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
		assert.JSONEq(t, `{"errors":{"sort":"unsupported sort order: \"SIDEWAYS\""}}`, recorder.Body.String())
	})
}
