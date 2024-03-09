package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"speed-typing-text-service/pkg/logging"
)

type Router struct {
	log logging.Logger
	*RegularTextHandler
	*CodeExampleHandler
}

func NewRouter(log logging.Logger, regularTextHandler *RegularTextHandler, codeExampleHandler *CodeExampleHandler) *Router {
	return &Router{
		log:                log,
		RegularTextHandler: regularTextHandler,
		CodeExampleHandler: codeExampleHandler,
	}
}

func (r *Router) InitRoutes(e *gin.Engine) {
	r.log.Info("initializing error handling middleware")
	e.Use(ErrorHandlerMiddleware())

	r.log.Info("initializing routes")

	// swagger
	e.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup := e.Group("/api")

	v1ApiGroup := apiGroup.Group("/v1")

	v1TextsGroup := v1ApiGroup.Group("/texts")
	{
		v1TextsGroup.GET("/", r.GetRegularTexts)

		v1TextsGroup.GET("/programming-languages", r.GetProgrammingLanguages)
		v1TextsGroup.GET("/code-examples/:uuid", r.GetCodeExampleByUUID)
		v1TextsGroup.GET("/code-examples", r.GetCodeExamples)
	}
}
