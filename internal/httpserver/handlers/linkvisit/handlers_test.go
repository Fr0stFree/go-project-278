package linkvisit

import (
	"net/http"
	"net/http/httptest"
	"shortener/internal/services/shortener"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockShortenerService struct {
	mock.Mock
}

func (m *mockShortenerService) ListLinkVisitsWithCount(optsBuilder *shortener.LinkVisitListOptionsBuilder) ([]shortener.LinkVisit, int, error) {
	args := m.Called(optsBuilder)

	return args.Get(0).([]shortener.LinkVisit), args.Int(1), args.Error(2)
}

type handlerMocks struct {
	shortener *mockShortenerService
	handler   *handler
}

func newHandlerMocks(t *testing.T) handlerMocks {
	t.Helper()

	shortener := new(mockShortenerService)
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

func newRouter(t *testing.T, shortener shortenerService) *gin.Engine {
	t.Helper()

	router := gin.New()
	RegisterRoutes(shortener, router)

	return router
}

func TestHandler_list(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should list link visits with count successfully", func(t *testing.T) {
		mocks := newHandlerMocks(t)
		mocks.shortener.
			On("ListLinkVisitsWithCount", mock.Anything).
			Return([]shortener.LinkVisit{
				{
					ID:        1,
					LinkID:    1,
					CreatedAt: "2026-09-04T12:00:00Z",
					IP:        "127.0.0.1",
					UserAgent: "Mozilla/5.0",
					Status:    http.StatusFound,
				},
			}, 1, nil).
			Once()

		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/link_visits", nil)

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

	t.Run("should parse range successfully", func(t *testing.T) {
		mocks := newHandlerMocks(t)
		mocks.shortener.
			On("ListLinkVisitsWithCount", mock.Anything).
			Return([]shortener.LinkVisit{}, 42, nil).
			Once()

		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/link_visits?range=%5B10%2C19%5D", nil)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "link_visits 10-19/42", recorder.Header().Get("Content-Range"))
		assert.JSONEq(t, `[]`, recorder.Body.String())
	})

	t.Run("should return bad request for invalid range", func(t *testing.T) {
		mocks := newHandlerMocks(t)
		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/link_visits?range=invalid", nil)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
	})
}
