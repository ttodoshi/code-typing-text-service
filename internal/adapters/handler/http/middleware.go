package http

import (
	"code-typing-text-service/internal/core/ports"
	"code-typing-text-service/internal/core/ports/dto"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"os"
	"time"
)

func ErrorHandlerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()

		if len(c.Errors) > 0 {
			err := c.Errors.Last()
			var responseStatus int
			if errors.Is(err, ports.BadRequestError) {
				responseStatus = http.StatusBadRequest
			} else if errors.Is(err, ports.UnauthorizedError) {
				responseStatus = http.StatusUnauthorized
			} else if errors.Is(err, ports.ForbiddenError) {
				responseStatus = http.StatusForbidden
			} else if errors.Is(err, ports.NotFoundError) {
				responseStatus = http.StatusNotFound
			} else {
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

func AuthenticationMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := c.GetHeader("Authorization")

		if tokenString == "" {
			return
		}
		token, err := jwt.Parse(tokenString[7:], func(token *jwt.Token) (interface{}, error) {
			return []byte(os.Getenv("SECRET_KEY")), nil
		})
		if err != nil || !token.Valid {
			return
		}

		if claims, ok := token.Claims.(jwt.MapClaims); ok {
			c.Set("userID", claims["sub"].(string))
			c.Set("nickname", claims["nickname"].(string))
		}

		c.Next()
	}
}
