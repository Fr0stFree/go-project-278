// Package httperror provides utilities for handling HTTP errors in the application.
package httperror

import (
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"net/http"

	"shortener/internal/common/textutils"
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
		jsonTypeErr   *json.UnmarshalTypeError
		jsonSyntaxErr *json.SyntaxError
		maxBytesErr   *http.MaxBytesError
		validatorErr  validator.ValidationErrors
	)

	switch {
	// Hexlet integration requires that the API returns a 400 status code for invalid JSON requests, and a 422 status code for validation errors.
	case errors.As(err, &jsonSyntaxErr) ||
		errors.As(err, &jsonTypeErr) ||
		errors.Is(err, io.EOF) ||
		errors.Is(err, io.ErrUnexpectedEOF): // 400
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
	case errors.As(err, &notFoundErr): // 404
		ctx.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Message})
	case errors.As(err, &conflictErr): // 409
		ctx.JSON(http.StatusConflict, gin.H{"error": map[string]string{conflictErr.Field: conflictErr.Message}})
	case errors.As(err, &maxBytesErr): // 413
		ctx.JSON(http.StatusRequestEntityTooLarge, gin.H{"error": "request body too large"})
	case errors.As(err, &validatorErr): // 422
		fieldErrors := make(map[string]string)
		for _, fieldErr := range validatorErr {
			fieldErrors[textutils.ToSnakeCase(fieldErr.Field())] = fieldErr.Error()
		}

		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"errors": fieldErrors})
	case errors.As(err, &validationErr): // 422
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": map[string]string{validationErr.Field: validationErr.Message}})
	default:
		slog.ErrorContext(
			ctx,
			"internal error",
			slog.String("reason", err.Error()),
			slog.String("operation", ctx.HandlerName()),
			slog.String("method", ctx.Request.Method),
			slog.String("url", ctx.Request.URL.String()),
		)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
}
