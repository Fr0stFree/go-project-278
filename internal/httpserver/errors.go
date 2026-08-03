package httpserver

import (
	"encoding/json"
	"errors"
	"net/http"

	"shortener/internal/services/shortener"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

// WriteErrorResponse translates a service error into an HTTP response.
func WriteErrorResponse(ctx *gin.Context, err error) {
	var (
		validationErr *shortener.ValidationError
		notFoundErr   *shortener.NotFoundError
		conflictErr   *shortener.ConflictError
		jsonTypeErr   *json.UnmarshalTypeError
		jsonSyntaxErr *json.SyntaxError
		validatorErr  validator.ValidationErrors
	)

	switch {
	case errors.As(err, &validationErr):
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": map[string]string{validationErr.Field: validationErr.Message}})
	case errors.As(err, &notFoundErr):
		ctx.JSON(http.StatusNotFound, gin.H{"error": notFoundErr.Message})
	case errors.As(err, &conflictErr):
		ctx.JSON(http.StatusConflict, gin.H{"error": map[string]string{conflictErr.Field: conflictErr.Message}})
	case errors.As(err, &jsonSyntaxErr) || errors.As(err, &jsonTypeErr):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
	case errors.As(err, &validatorErr):
		ctx.JSON(http.StatusBadRequest, gin.H{"error": validatorErr.Error()})
	default:
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "something went wrong"})
	}
}
