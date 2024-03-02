package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
	_ "speed-typing-text-service/docs"
	"speed-typing-text-service/internal/adapters/handler"
	"speed-typing-text-service/internal/adapters/repository/postgres"
	"speed-typing-text-service/internal/core/servises"
	"speed-typing-text-service/pkg/discovery"
	"speed-typing-text-service/pkg/env"
	"speed-typing-text-service/pkg/logging"
)

const (
	Dev  = "dev"
	Prod = "prod"
)

var (
	codeExampleHandler *handler.CodeExampleHandler
	regularTextHandler *handler.RegularTextHandler
	log                logging.Logger
)

func init() {
	env.LoadEnvVariables()
	if os.Getenv("PROFILE") == Prod {
		gin.SetMode(gin.ReleaseMode)
	}
	log = logging.GetLogger()
	discovery.InitServiceDiscovery()
}

func main() {
	codeExampleRepository := postgres.NewCodeExampleRepository()
	regularTextRepository := postgres.NewRegularTextRepository()
	codeExampleService := servises.NewCodeExampleService(codeExampleRepository, log)
	regularTextService := servises.NewRegularTextService(regularTextRepository, log)
	codeExampleHandler = handler.NewCodeExampleHandler(codeExampleService, log)
	regularTextHandler = handler.NewRegularTextHandler(regularTextService, log)

	initRoutes()
}

//	@title						Text Generation Service API
//	@version					1.0

//	@host						localhost:8080
//	@BasePath					/api/v1

// @externalDocs.description	OpenAPI
func initRoutes() {
	r := gin.Default()

	log.Info("initializing error handling middleware")
	r.Use(handler.ErrorHandlerMiddleware())

	log.Info("initializing handlers")

	// swagger
	r.GET("/swagger-ui/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	apiGroup := r.Group("/api")

	v1ApiGroup := apiGroup.Group("/v1")

	v1TextsGroup := v1ApiGroup.Group("/texts")
	{
		v1TextsGroup.GET("/", regularTextHandler.GetRegularTexts)

		v1TextsGroup.GET("/programming-languages", codeExampleHandler.GetProgrammingLanguages)
		v1TextsGroup.GET("/code-examples/:uuid", codeExampleHandler.GetCodeExampleByUUID)
		v1TextsGroup.GET("/code-examples", codeExampleHandler.GetCodeExamples)
	}

	log.Infof("starting server on port :%s", os.Getenv("PORT"))

	err := r.Run()
	if err != nil {
		log.Fatal("error while running server")
	}
}
