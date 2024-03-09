package main

import (
	"github.com/gin-gonic/gin"
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

func init() {
	env.LoadEnvVariables()
	if os.Getenv("PROFILE") == Prod {
		gin.SetMode(gin.ReleaseMode)
	}
	discovery.InitServiceDiscovery()
}

//	@title		Text Generation Service API
//	@version	1.0

// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	log := logging.GetLogger()

	r := gin.Default()
	router := initRouter(log)
	router.InitRoutes(r)

	log.Fatalf("error while running server due to: %s", r.Run())
}

func initRouter(log logging.Logger) *handler.Router {
	codeExampleRepository := postgres.NewCodeExampleRepository(log)
	regularTextRepository := postgres.NewRegularTextRepository(log)
	codeExampleService := servises.NewCodeExampleService(codeExampleRepository, log)
	regularTextService := servises.NewRegularTextService(regularTextRepository, log)
	return handler.NewRouter(
		log,
		handler.NewRegularTextHandler(
			regularTextService, log,
		), handler.NewCodeExampleHandler(
			codeExampleService, log,
		),
	)
}
