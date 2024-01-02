package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"speed-typing-text-service/internal/core/errors"
	"time"
)

type errorResponse struct {
	Timestamp time.Time `json:"timestamp"`
	Status    int       `json:"status,omitempty"`
	Error     string    `json:"error,omitempty"`
	Message   string    `json:"message,omitempty"`
	Path      string    `json:"path,omitempty"`
}

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
			c.JSON(responseStatus, errorResponse{
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
