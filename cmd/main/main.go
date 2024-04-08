package main

import (
	_ "code-typing-text-service/docs"
	"code-typing-text-service/internal/adapters/handler/http"
	"code-typing-text-service/internal/adapters/handler/http/api"
	"code-typing-text-service/internal/adapters/repository/postgres"
	"code-typing-text-service/internal/core/servises"
	"code-typing-text-service/pkg/discovery"
	"code-typing-text-service/pkg/env"
	"code-typing-text-service/pkg/logging"
	"github.com/gin-gonic/gin"
	"os"
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

func initRouter(log logging.Logger) *http.Router {
	codeExampleRepository := postgres.NewCodeExampleRepository(log)
	codeExampleService := servises.NewCodeExampleService(codeExampleRepository, log)
	return http.NewRouter(
		log,
		api.NewCodeExampleHandler(
			codeExampleService, log,
		),
	)
}
