package httpserver

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func setupRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)

	router := newRouter()

	return router
}

func TestPingHandler(t *testing.T) {
	router := setupRouter()
	request := httptest.NewRequest(http.MethodGet, "/ping", nil)
	response := httptest.NewRecorder()

	router.ServeHTTP(response, request)

	require.Equal(t, http.StatusOK, response.Code)
	require.Equal(t, "pong", response.Body.String())
	require.Equal(t, "text/plain; charset=utf-8", response.Header().Get("Content-Type"))
}
