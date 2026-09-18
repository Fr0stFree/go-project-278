package httperror

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"shortener/internal/services/shortener"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestWriteErrorResponse(t *testing.T) {
	gin.SetMode(gin.TestMode)

	type subTest struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   string
	}

	subTests := []subTest{
		{
			name:           "should return 400 for invalid JSON",
			err:            &json.SyntaxError{},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request"}`,
		},
		{
			name:           "should return 400 for empty body",
			err:            io.EOF,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request"}`,
		},
		{
			name:           "should return 400 for unexpected EOF",
			err:            io.ErrUnexpectedEOF,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `{"error":"invalid request"}`,
		},
		{
			name: "should return 404 for not found error",
			err: &shortener.NotFoundError{
				Message: "link not found",
			},
			expectedStatus: http.StatusNotFound,
			expectedBody:   `{"error":"link not found"}`,
		},
		{
			name: "should return 409 for conflict error",
			err: &shortener.ConflictError{
				Field:   "short_name",
				Message: "already exists",
			},
			expectedStatus: http.StatusConflict,
			expectedBody:   `{"error":"already exists"}`,
		},
		{
			name: "should return 413 for request body too large",
			err: &http.MaxBytesError{
				Limit: 16 * 1024,
			},
			expectedStatus: http.StatusRequestEntityTooLarge,
			expectedBody:   `{"error":"request body too large"}`,
		},
		{
			name: "should return 422 for service validation error",
			err: &shortener.ValidationError{
				Field:   "original_url",
				Message: "invalid URL",
			},
			expectedStatus: http.StatusUnprocessableEntity,
			expectedBody:   `{"errors":{"original_url":"invalid URL"}}`,
		},
		{
			name:           "should return 500 for unknown error",
			err:            errors.New("database unavailable"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `{"error":"something went wrong"}`,
		},
	}

	for _, subTest := range subTests {
		t.Run(subTest.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			ctx, _ := gin.CreateTestContext(recorder)
			ctx.Request = httptest.NewRequest(http.MethodGet, "/test", nil)

			WriteResponse(ctx, subTest.err)

			assert.Equal(t, subTest.expectedStatus, recorder.Code)
			assert.Equal(t, "application/json; charset=utf-8", recorder.Header().Get("Content-Type"))
			assert.JSONEq(t, subTest.expectedBody, recorder.Body.String())
		})
	}
}
