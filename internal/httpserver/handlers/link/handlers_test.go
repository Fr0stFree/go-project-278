package link

import (
	"context"
	"errors"
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

type mockService struct {
	mock.Mock
}

func (m *mockService) GetRedirectLink(ctx context.Context, shortName string) (shortener.Link, error) {
	args := m.Called(ctx, shortName)

	return args.Get(0).(shortener.Link), args.Error(1)
}

func (m *mockService) SaveLinkVisit(ctx context.Context, linkID uint, ip, userAgent, referrer string, status uint) (shortener.LinkVisit, error) {
	args := m.Called(ctx, linkID, ip, userAgent, referrer, status)

	return args.Get(0).(shortener.LinkVisit), args.Error(1)
}

func (m *mockService) CreateLink(ctx context.Context, originalURL, shortName string) (shortener.Link, error) {
	args := m.Called(ctx, originalURL, shortName)

	return args.Get(0).(shortener.Link), args.Error(1)
}

func (m *mockService) GetLink(ctx context.Context, id uint) (shortener.Link, error) {
	args := m.Called(ctx, id)

	return args.Get(0).(shortener.Link), args.Error(1)
}

func (m *mockService) ListLinksWithCount(ctx context.Context, optsBuilder *shortener.LinkListOptionsBuilder) ([]shortener.Link, int, error) {
	args := m.Called(ctx, optsBuilder)

	return args.Get(0).([]shortener.Link), args.Int(1), args.Error(2)
}

func (m *mockService) UpdateLink(ctx context.Context, id uint, originalURL, shortName string) (shortener.Link, error) {
	args := m.Called(ctx, id, originalURL, shortName)

	return args.Get(0).(shortener.Link), args.Error(1)
}

func (m *mockService) DeleteLink(ctx context.Context, id uint) error {
	args := m.Called(ctx, id)

	return args.Error(0)
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

func newRouter(t *testing.T, shortener Service) *gin.Engine {
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
			On("GetRedirectLink", mock.Anything, shortName).
			Return(shortener.Link{
				ID:          1,
				OriginalURL: originalURL,
				ShortName:   shortName,
				ShortURL:    "http://localhost/r/" + shortName,
			}, nil).
			Once()
		mocks.shortener.
			On("SaveLinkVisit", mock.Anything, uint(1), ip, userAgent, referrer, status).
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

	t.Run("should redirect even when saving link visit fails", func(t *testing.T) {
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
			On("GetRedirectLink", mock.Anything, shortName).
			Return(shortener.Link{
				ID:          1,
				OriginalURL: originalURL,
				ShortName:   shortName,
				ShortURL:    "http://localhost/r/" + shortName,
			}, nil).
			Once()

		mocks.shortener.
			On("SaveLinkVisit", mock.Anything, uint(1), ip, userAgent, referrer, status).
			Return(shortener.LinkVisit{}, errors.New("database is unavailable")).
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
}

func TestHandler_create(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should create link successfully", func(t *testing.T) {
		const (
			originalURL = "https://example.com"
			shortName   = "abc-123"
			shortURL    = "https://short.example.com/abc-123"
		)

		mocks := newHandlerMocks(t)
		mocks.shortener.
			On("CreateLink", mock.Anything, originalURL, shortName).
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
				"short_name": "abc-123"
			}`),
		)
		request.Header.Set("Content-Type", "application/json")

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusCreated, recorder.Code)
		assert.JSONEq(t, `{
			"id": 1,
			"original_url": "https://example.com",
			"short_name": "abc-123",
			"short_url": "https://short.example.com/abc-123"
		}`, recorder.Body.String())
	})

	t.Run("should handle conflict error when creating a link with an existing short name", func(t *testing.T) {
		const (
			originalURL = "https://example.com"
			shortName   = "abc123"
		)

		mocks := newHandlerMocks(t)
		mocks.shortener.
			On("CreateLink", mock.Anything, originalURL, shortName).
			Return(shortener.Link{}, &shortener.ConflictError{
				Field:   "short_name",
				Message: "short name already exists",
			}).
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

		require.Equal(t, http.StatusConflict, recorder.Code)
		assert.JSONEq(t, `{"error": {"short_name": "short name already exists"}}`, recorder.Body.String())
	})

	t.Run("should handle invalid request data properly", func(t *testing.T) {
		mocks := newHandlerMocks(t)

		subtests := []struct {
			name             string
			body             string
			expectedResponse string
			expectedStatus   int
		}{
			{
				name:             "invalid json",
				body:             "invalid json",
				expectedResponse: `{"error": "invalid request"}`,
				expectedStatus:   http.StatusBadRequest,
			},
			{
				name:             "missing original_url",
				body:             `{"short_name": "abc123"}`,
				expectedResponse: `{"errors": {"original_url": "Key: 'createLinkRequestBody.OriginalURL' Error:Field validation for 'OriginalURL' failed on the 'required' tag"}}`,
				expectedStatus:   http.StatusUnprocessableEntity,
			},
			{
				name:             "short_name containing a path separator",
				body:             `{"original_url": "https://example.com", "short_name": "abr/123"}`,
				expectedResponse: `{"errors": {"short_name": "Key: 'createLinkRequestBody.ShortName' Error:Field validation for 'ShortName' failed on the 'excludes' tag"}}`,
				expectedStatus:   http.StatusUnprocessableEntity,
			},
			{
				name:             "invalid original_url",
				body:             `{"original_url": "invalid-url", "short_name": "abc123"}`,
				expectedResponse: `{"errors": {"original_url": "Key: 'createLinkRequestBody.OriginalURL' Error:Field validation for 'OriginalURL' failed on the 'http_url' tag"}}`,
				expectedStatus:   http.StatusUnprocessableEntity,
			},
			{
				name:             "empty body",
				body:             "",
				expectedResponse: `{"error": "invalid request"}`,
				expectedStatus:   http.StatusBadRequest,
			},
		}
		for _, subtest := range subtests {
			t.Run(subtest.name, func(t *testing.T) {
				router := newRouter(t, mocks.shortener)
				recorder := httptest.NewRecorder()
				request := httptest.NewRequest(http.MethodPost, "/api/links", strings.NewReader(subtest.body))
				request.Header.Set("Content-Type", "application/json")

				router.ServeHTTP(recorder, request)

				require.Equal(t, subtest.expectedStatus, recorder.Code)
				assert.JSONEq(t, subtest.expectedResponse, recorder.Body.String())
			})
		}
	})
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
			On("GetLink", mock.Anything, uint(linkID)).
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

	t.Run("should handle not found error when getting a non-existing link", func(t *testing.T) {
		const linkID = 999

		mocks := newHandlerMocks(t)
		mocks.shortener.
			On("GetLink", mock.Anything, uint(linkID)).
			Return(shortener.Link{}, &shortener.NotFoundError{
				Message: "link not found",
			}).
			Once()

		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/links/999", nil)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusNotFound, recorder.Code)
		assert.JSONEq(t, `{"error": "link not found"}`, recorder.Body.String())
	})

	t.Run("should handle invalid link ID parameter", func(t *testing.T) {
		mocks := newHandlerMocks(t)

		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/links/invalid-id", nil)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
		assert.JSONEq(t, `{"error": {"link_id": "invalid positive integer: invalid-id"}}`, recorder.Body.String())
	})
}

func TestHandler_list(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("should list links successfully", func(t *testing.T) {
		mocks := newHandlerMocks(t)
		mocks.shortener.
			On("ListLinksWithCount", mock.Anything, mock.Anything).
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

	t.Run("should return empty content range when requested page is empty", func(t *testing.T) {
		mocks := newHandlerMocks(t)
		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/links?range=[10,19]", nil)

		mocks.shortener.
			On("ListLinksWithCount", mock.Anything,
				mock.MatchedBy(func(builder *shortener.LinkListOptionsBuilder) bool {
					from, to := builder.Range()

					return from == 10 && to == 19
				}),
			).
			Return([]shortener.Link{}, 42, nil).
			Once()

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "links */42", recorder.Header().Get("Content-Range"))
		assert.JSONEq(t, `[]`, recorder.Body.String())
	})

	t.Run("should build content range from actual number of returned links", func(t *testing.T) {
		mocks := newHandlerMocks(t)
		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/links?range=[10,19]", nil)

		mocks.shortener.
			On("ListLinksWithCount", mock.Anything, mock.MatchedBy(func(builder *shortener.LinkListOptionsBuilder) bool {
				from, to := builder.Range()

				return from == 10 && to == 19
			})).
			Return([]shortener.Link{{ID: 11}, {ID: 12}, {ID: 13}, {ID: 14}}, 14, nil).
			Once()

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		assert.Equal(t, "links 10-13/14", recorder.Header().Get("Content-Range"))
	})

	t.Run("should handle invalid range parameter", func(t *testing.T) {
		mocks := newHandlerMocks(t)
		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, "/api/links?range=invalid-range", nil)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
		assert.JSONEq(t, `{"error": {"range": "invalid range format: invalid-range"}}`, recorder.Body.String())
	})

	t.Run("should handle sort parameter", func(t *testing.T) {
		mocks := newHandlerMocks(t)
		mocks.shortener.
			On("ListLinksWithCount", mock.Anything, mock.Anything).
			Return([]shortener.Link{}, 0, nil).
			Once()

		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, `/api/links?sort=["original_url","asc"]`, nil)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusOK, recorder.Code)
		assert.JSONEq(t, `[]`, recorder.Body.String())
	})

	t.Run("should handle invalid sort parameter", func(t *testing.T) {
		mocks := newHandlerMocks(t)
		router := newRouter(t, mocks.shortener)
		recorder := httptest.NewRecorder()
		request := httptest.NewRequest(http.MethodGet, `/api/links?sort=invalid-sort`, nil)

		router.ServeHTTP(recorder, request)

		require.Equal(t, http.StatusUnprocessableEntity, recorder.Code)
		assert.JSONEq(t, `{"error": {"sort": "invalid sort format: invalid-sort"}}`, recorder.Body.String())
	})
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
			On("UpdateLink", mock.Anything, uint(linkID), originalURL, shortName).
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
			On("DeleteLink", mock.Anything, uint(linkID)).
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
