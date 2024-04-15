package http

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/ttodoshi/code-typing-text-service/internal/adapters/handler/http/api"
	"github.com/ttodoshi/code-typing-text-service/pkg/logging"
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

	r.log.Info("initializing routes")

	// swagger
	e.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

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
