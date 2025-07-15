package http

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ttodoshi/code-typing-text-service/internal/adapters/handler/http/api"
	"github.com/ttodoshi/code-typing-text-service/pkg/logging"
	"net/http"
)

type Router struct {
	log logging.Logger
	*api.CodeExampleHandler
}

func NewRouter(log logging.Logger, codeExampleHandler *api.CodeExampleHandler) *Router {
	return &Router{
		log:                log,
		CodeExampleHandler: codeExampleHandler,
	}
}

func (r *Router) InitRoutes(e *gin.Engine) {
	r.log.Info("initializing error handling middleware")
	e.Use(ErrorHandlerMiddleware())
	e.Use(AuthenticationMiddleware())
	e.Use(cors.New(cors.Config{
		AllowOriginFunc: func(origin string) bool {
			return true
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS", "HEAD"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.log.Info("initializing routes")

	// swagger
	e.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// healthcheck
	e.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	apiGroup := e.Group("/api")

	v1ApiGroup := apiGroup.Group("/v1")

	v1TextsGroup := v1ApiGroup.Group("/texts")
	{
		v1TextsGroup.GET("/programming-languages", r.GetProgrammingLanguages)
		v1TextsGroup.GET("/code-examples/:uuid", r.GetCodeExampleByUUID)
		v1TextsGroup.GET("/code-examples", r.GetCodeExamples)

		v1TextsGroup.POST("/code-examples", r.CreateCodeExample)
		v1TextsGroup.DELETE("/code-examples/:uuid", r.DeleteCodeExample)
	}
}
