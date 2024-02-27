package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"speed-typing-text-service/internal/adapters/dto"
	"speed-typing-text-service/internal/core/errors"
	"time"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			var responseStatus int
			switch err.Err.(type) {
			case *errors.NotFoundError:
				responseStatus = http.StatusNotFound
			case *errors.MappingError:
				responseStatus = http.StatusInternalServerError
			default:
				responseStatus = http.StatusInternalServerError
			}
			c.JSON(responseStatus, dto.ErrorResponseDto{
				Timestamp: time.Now(),
				Status:    responseStatus,
				Error:     http.StatusText(responseStatus),
				Message:   err.Err.Error(),
				Path:      c.Request.URL.Path,
			})
			c.Abort()
			return
		}
	}
}
