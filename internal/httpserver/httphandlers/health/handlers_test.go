package health

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type readinessCheckerMock struct {
	err error
}

func (s readinessCheckerMock) PingContext(context.Context) error {
	return s.err
}

func TestPing(t *testing.T) {
	gin.SetMode(gin.TestMode)

	recorder := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(recorder)
	h := &handler{}

	h.ping(ctx)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "pong", recorder.Body.String())
}

func TestHealth(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type subTest struct {
		name           string
		checker        ReadinessChecker
		expectedStatus int
		expectedBody   string
	}

	subTests := []subTest{
		{
			name:           "should report ready when PostgreSQL is available",
			checker:        readinessCheckerMock{},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"ok":true}`,
		},
		{
			name:           "should report not ready without exposing the database error",
			checker:        readinessCheckerMock{err: errors.New("some database error")},
			expectedStatus: http.StatusServiceUnavailable,
			expectedBody:   `{"ok":false}`,
		},
	}

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/health", nil)

			h := &handler{
				checker:            subTest.checker,
				healthcheckTimeout: time.Second,
			}

			h.health(ctx)

			require.Equal(t, subTest.expectedStatus, recorder.Code)
			assert.JSONEq(t, subTest.expectedBody, recorder.Body.String())
			assert.NotContains(t, recorder.Body.String(), "some database error")
		})
	}
}
