// Package httperror provides utilities for handling HTTP errors in the application.
package httperror

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"shortener/internal/common/textutils"
	"shortener/internal/httpserver/httptools/httpparam"
	"shortener/internal/services/shortener"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// WriteResponse translates a service error into an HTTP response.
func WriteResponse(ctx *gin.Context, err error) {
	var (
		validationErr *shortener.ValidationError
		notFoundErr   *shortener.NotFoundError
		conflictErr   *shortener.ConflictError
		maxBytesErr   *http.MaxBytesError
		validatorErr  validator.ValidationErrors
	)

	switch {
	case isMalformedJSONError(err): // 400
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{"error": "invalid request"},
		)
	case errors.As(err, &notFoundErr): // 404
		ctx.JSON(
			http.StatusNotFound,
			gin.H{"error": notFoundErr.Message},
		)
	case errors.As(err, &conflictErr): // 409
		ctx.JSON(
			http.StatusConflict,
			gin.H{"error": map[string]string{conflictErr.Field: conflictErr.Message}},
		)
	case errors.As(err, &maxBytesErr): // 413
		ctx.JSON(
			http.StatusRequestEntityTooLarge,
			gin.H{"error": "request body too large"},
		)
	case errors.As(err, &validatorErr): // 422
		fieldErrors := make(map[string]string)
		for _, fieldErr := range validatorErr {
			fieldErrors[textutils.ToSnakeCase(fieldErr.Field())] = fieldErr.Error()
		}

		ctx.JSON(
			http.StatusUnprocessableEntity,
			gin.H{"errors": fieldErrors},
		)
	case errors.As(err, &validationErr): // 422
		ctx.JSON(
			http.StatusUnprocessableEntity,
			gin.H{"errors": map[string]string{validationErr.Field: validationErr.Message}},
		)
	default: // 500
		slog.ErrorContext(
			ctx,
			"internal error",
			slog.String("request_id", httpparam.ReadRequestIDContext(ctx)),
			slog.String("reason", err.Error()),
			slog.String("operation", ctx.HandlerName()),
			slog.String("method", ctx.Request.Method),
			slog.String("url", ctx.Request.URL.String()),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
}

func isMalformedJSONError(err error) bool {
	var (
		jsonTypeErr   *json.UnmarshalTypeError
		jsonSyntaxErr *json.SyntaxError
	)

	return errors.As(err, &jsonTypeErr) ||
		errors.As(err, &jsonSyntaxErr) ||
		errors.Is(err, io.EOF) ||
		errors.Is(err, io.ErrUnexpectedEOF)
}
