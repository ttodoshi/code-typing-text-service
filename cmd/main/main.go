package main

import (
	"github.com/gin-gonic/gin"
	"github.com/hashicorp/consul/api"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"os"
	_ "speed-typing-text-service/docs"
	"speed-typing-text-service/internal/adapters/handler"
	"speed-typing-text-service/internal/adapters/repository/postgres"
	"speed-typing-text-service/internal/core/ports"
	"speed-typing-text-service/internal/core/servises"
	"speed-typing-text-service/pkg/env"
	"speed-typing-text-service/pkg/logging"
	"strconv"
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
	initConsul()
	initRoutes()
}

//	@title						Text Generation Service API
//	@version					1.0
//	@host						localhost:8080
//	@BasePath					/api/v1
//	@externalDocs.description	OpenAPI
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
		regularTextHandler := handler.NewRegularTextHandler(regularTextService, log)
		v1TextsGroup.GET("/", regularTextHandler.GetRegularTexts)

		codeExampleHandler := handler.NewCodeExampleHandler(codeExampleService, log)
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

func initConsul() {
	log.Info("initializing consul client")

	consulClient, err := api.NewClient(
		&api.Config{
			Address: os.Getenv("CONSUL_HOST"),
		},
	)
	if err != nil {
		log.Fatal("error creating consul client")
	}

	log.Info("register service in consul")
	agent := consulClient.Agent()
	parsedPort, err := strconv.Atoi(os.Getenv("PORT"))
	if err != nil {
		log.Fatal("port parse error")
	}

	service := &api.AgentServiceRegistration{
		Name:    os.Getenv("CONSUL_SERVICE_NAME"),
		Port:    parsedPort,
		Address: os.Getenv("CONSUL_SERVICE_ADDRESS"),
	}
	err = agent.ServiceRegister(service)
	if err != nil {
		log.Fatalf("error while service registration due to error '%s'", err)
	}
}
