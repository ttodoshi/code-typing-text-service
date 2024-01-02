package main

import (
	"github.com/gin-gonic/gin"
	"os"
	"speed-typing-text-service/internal/adapters/handler"
	"speed-typing-text-service/internal/adapters/repository/postgres"
	"speed-typing-text-service/internal/core/ports"
	"speed-typing-text-service/internal/core/servises"
	"speed-typing-text-service/pkg/env"
	"speed-typing-text-service/pkg/logging"
)

const (
	Dev  = "dev"
	Prod = "prod"
)

var (
	codeExampleService ports.CodeExampleService
	regularTextService ports.RegularTextService
	log                logging.Logger
)

func init() {
	env.LoadEnvVariables()
	if os.Getenv("PROFILE") == Prod {
		gin.SetMode(gin.ReleaseMode)
	}
	log = logging.GetLogger()
}

func main() {
	codeExampleRepository := postgres.NewCodeExampleRepository()
	regularTextRepository := postgres.NewRegularTextRepository()
	codeExampleService = servises.NewCodeExampleService(codeExampleRepository, log)
	regularTextService = servises.NewRegularTextService(regularTextRepository, log)
	initRoutes()
}

func initRoutes() {
	r := gin.Default()

	log.Info("initializing error handling middleware")
	r.Use(handler.ErrorHandlerMiddleware())

	log.Info("initializing handlers")

	apiGroup := r.Group("/api")

	v1ApiGroup := apiGroup.Group("/v1")

	v1RegularTextsGroup := v1ApiGroup.Group("/texts")
	{
		regularTextHandler := handler.NewRegularTextHandler(regularTextService, log)
		v1RegularTextsGroup.GET("/", regularTextHandler.GetRegularTexts)
	}

	codeExampleHandler := handler.NewCodeExampleHandler(codeExampleService, log)
	v1ProgrammingLanguagesGroup := v1ApiGroup.Group("/programming-languages")
	{
		v1ProgrammingLanguagesGroup.GET("/", codeExampleHandler.GetProgrammingLanguages)
	}

	v1CodeExamplesGroup := v1ApiGroup.Group("/codes")
	{
		v1CodeExamplesGroup.GET("/", codeExampleHandler.GetCodeExamples)
	}

	log.Infof("starting server on port :%s", os.Getenv("PORT"))

	err := r.Run()
	if err != nil {
		log.Fatalf("error while running server")
	}
}
