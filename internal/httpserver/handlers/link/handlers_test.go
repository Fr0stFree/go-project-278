package link

import (
	"net/http"
	"net/http/httptest"
	"shortener/internal/services/shortener"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

type mockShortenerService struct {
	mock.Mock
}

func (m *mockShortenerService) GetRedirectLink(shortName string) (shortener.Link, error) {
	args := m.Called(shortName)

	return args.Get(0).(shortener.Link), args.Error(1)
}

func (m *mockShortenerService) SaveLinkVisit(linkID uint, ip, userAgent, referrer string, status uint) (shortener.LinkVisit, error) {
	args := m.Called(linkID, ip, userAgent, referrer, status)

	return args.Get(0).(shortener.LinkVisit), args.Error(1)
}

func (m *mockShortenerService) CreateLink(originalURL, shortName string) (shortener.Link, error) {
	args := m.Called(originalURL, shortName)

	return args.Get(0).(shortener.Link), args.Error(1)
}

func (m *mockShortenerService) GetLink(id uint) (shortener.Link, error) {
	args := m.Called(id)

	return args.Get(0).(shortener.Link), args.Error(1)
}

func (m *mockShortenerService) ListLinksWithCount(optsBuilder *shortener.LinkListOptionsBuilder) ([]shortener.Link, int, error) {
	args := m.Called(optsBuilder)

	return args.Get(0).([]shortener.Link), args.Int(1), args.Error(2)
}

func (m *mockShortenerService) UpdateLink(id uint, originalURL, shortName string) (shortener.Link, error) {
	args := m.Called(id, originalURL, shortName)

	return args.Get(0).(shortener.Link), args.Error(1)
}

func (m *mockShortenerService) DeleteLink(id uint) error {
	args := m.Called(id)

	return args.Error(0)
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

func TestHandler_redirect(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should redirect to original URL and save link visit", func(t *testing.T) {
		const (
			shortName        = "abc123"
			originalURL      = "https://example.com"
			ip               = "127.0.0.1"
			userAgent        = "Mozilla/5.0"
			referrer         = "https://referrer.com"
			status      uint = http.StatusFound
		)

		mocks := newHandlerMocks(t)
		mocks.shortener.
			On("GetRedirectLink", shortName).
			Return(shortener.Link{
				ID:          1,
				OriginalURL: originalURL,
				ShortName:   shortName,
				ShortURL:    "http://localhost/r/" + shortName,
			}, nil).
			Once()
		mocks.shortener.
			On("SaveLinkVisit", uint(1), ip, userAgent, referrer, status).
			Return(shortener.LinkVisit{
				ID:        1,
				LinkID:    1,
				IP:        ip,
				UserAgent: userAgent,
				Referrer:  referrer,
				Status:    status,
			}, nil).
			Once()

		recorder := httptest.NewRecorder()
		router := newRouter(t, mocks.shortener)
		request := httptest.NewRequest(http.MethodGet, "/r/"+shortName, nil)
		request.RemoteAddr = ip + ":12345"
		request.Header.Set("User-Agent", userAgent)
		request.Header.Set("Referer", referrer)

		router.ServeHTTP(recorder, request)

		assert.Equal(t, http.StatusFound, recorder.Code)
		assert.Equal(t, originalURL, recorder.Header().Get("Location"))
	})

	// TODO: add more test cases
}

func TestHandler_create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should create link successfully", func(t *testing.T) {
		const (
			originalURL = "https://example.com"
			shortName   = "abc123"
			shortURL    = "https://short.example.com/abc123"
		)

		mocks := newHandlerMocks(t)
		mocks.shortener.
			On("CreateLink", originalURL, shortName).
			Return(shortener.Link{
				ID:          1,
				OriginalURL: originalURL,
				ShortName:   shortName,
				ShortURL:    shortURL,
			}, nil).
			Once()

		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(
			http.MethodPost,
			"/api/links",
			strings.NewReader(`{
				"original_url": "https://example.com",
				"short_name": "abc123"
			}`),
		)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusCreated, recorder.Code)
		assert.JSONEq(t, `{
			"id": 1,
			"original_url": "https://example.com",
			"short_name": "abc123",
			"short_url": "https://short.example.com/abc123"
		}`, recorder.Body.String())
	})

	// TODO: add more test cases
}

func TestHandler_get(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should get link successfully", func(t *testing.T) {
		const (
			linkID      = 1
			originalURL = "https://example.com"
			shortName   = "abc123"
			shortURL    = "https://short.example.com/abc123"
		)

		mocks := newHandlerMocks(t)
		mocks.shortener.
			On("GetLink", uint(linkID)).
			Return(shortener.Link{
				ID:          linkID,
				OriginalURL: originalURL,
				ShortName:   shortName,
				ShortURL:    shortURL,
			}, nil).
			Once()

		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/links/1", nil)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		assert.JSONEq(t, `{
			"id": 1,
			"original_url": "https://example.com",
			"short_name": "abc123",
			"short_url": "https://short.example.com/abc123"
		}`, recorder.Body.String())
	})

	// TODO: add more test cases
}

func TestHandler_list(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should list links successfully", func(t *testing.T) {
		mocks := newHandlerMocks(t)
		mocks.shortener.
			On("ListLinksWithCount", mock.Anything).
			Return([]shortener.Link{
				{
					ID:          1,
					OriginalURL: "https://example.com",
					ShortName:   "abc123",
					ShortURL:    "https://short.example.com/abc123",
				},
			}, 1, nil).
			Once()

		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/links", nil)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		assert.JSONEq(t, `[
			{
				"id": 1,
				"original_url": "https://example.com",
				"short_name": "abc123",
				"short_url": "https://short.example.com/abc123"
			}
		]`, recorder.Body.String())
	})

	t.Run("should parse range and set content range header", func(t *testing.T) {
		mocks := newHandlerMocks(t)
		mocks.shortener.
			On("ListLinksWithCount", mock.Anything).
			Return([]shortener.Link{}, 42, nil).
			Once()

		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/links?range=[10,19]", nil)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "links 10-19/42", recorder.Header().Get("Content-Range"))
		assert.JSONEq(t, `[]`, recorder.Body.String())
	})

	// TODO: add more test cases
}

func TestHandler_update(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should update link successfully", func(t *testing.T) {
		const (
			linkID      = 1
			originalURL = "https://example.com"
			shortName   = "abc123"
			shortURL    = "https://short.example.com/abc123"
		)

		mocks := newHandlerMocks(t)
		mocks.shortener.
			On("UpdateLink", uint(linkID), originalURL, shortName).
			Return(shortener.Link{
				ID:          linkID,
				OriginalURL: originalURL,
				ShortName:   shortName,
				ShortURL:    shortURL,
			}, nil).
			Once()

		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(
			http.MethodPut,
			"/api/links/1",
			strings.NewReader(`{
				"original_url": "https://example.com",
				"short_name": "abc123"
			}`),
		)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		assert.JSONEq(t, `{
			"id": 1,
			"original_url": "https://example.com",
			"short_name": "abc123",
			"short_url": "https://short.example.com/abc123"
		}`, recorder.Body.String())
	})

	// TODO: add more test cases
}

func TestHandler_delete(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should delete link successfully", func(t *testing.T) {
		const linkID = 1

		mocks := newHandlerMocks(t)
		mocks.shortener.
			On("DeleteLink", uint(linkID)).
			Return(nil).
			Once()

		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodDelete, "/api/links/1", nil)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusNoContent, recorder.Code)
	})

	// TODO: add more test cases
}
