package middleware

import (
	"github.com/gin-gonic/gin"
	"github.com/snisid/platform/backend/internal/platform/errors"
	"github.com/snisid/platform/backend/internal/platform/validation"
)

// DTOContextKey is a generic context key used to store the parsed object
type dtoContextKey struct {
	Type string
}

// GetValidated extracts the strongly-typed DTO from the gin Context.
// Must be used downstream of ValidateBody or ValidateQuery.
func GetValidated[T any](c *gin.Context) (T, bool) {
	var zero T
	val, exists := c.Get("validated_dto")
	if !exists {
		return zero, false
	}
	dto, ok := val.(T)
	return dto, ok
}

// ValidateBody reads the JSON body, binds it to generic type T, runs schema validation,
// and saves the parsed struct into the context for the handler to use.
func ValidateBody[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.ShouldBindJSON(&req); err != nil {
			// Bind failure or Validation failure
			platformErr := validation.TranslateError(err, "middleware.ValidateBody")
			errors.RespondWithError(c, platformErr)
			return
		}

		// Also run structural validation explicitly just in case Gin skipped it
		if err := validation.Struct(req); err != nil {
			platformErr := validation.TranslateError(err, "middleware.ValidateBody")
			errors.RespondWithError(c, platformErr)
			return
		}

		c.Set("validated_dto", req)
		c.Next()
	}
}

// ValidateQuery reads the URL Query string, binds it to generic type T,
// runs schema validation, and saves the parsed struct into the context.
func ValidateQuery[T any]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var req T
		if err := c.ShouldBindQuery(&req); err != nil {
			platformErr := validation.TranslateError(err, "middleware.ValidateQuery")
			errors.RespondWithError(c, platformErr)
			return
		}

		if err := validation.Struct(req); err != nil {
			platformErr := validation.TranslateError(err, "middleware.ValidateQuery")
			errors.RespondWithError(c, platformErr)
			return
		}

		c.Set("validated_dto", req)
		c.Next()
	}
}
