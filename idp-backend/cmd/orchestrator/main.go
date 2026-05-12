package main

import (
	"log"

	"idp-backend/internal/config"
	"idp-backend/internal/handler"
	"idp-backend/internal/repository"
	"idp-backend/internal/service"
	"idp-backend/pkg/db"
	"idp-backend/pkg/k8s"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	// 0. Load .env file
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 1. Load Config
	cfg := config.LoadConfig()

	// 2. Initialize Dependencies
	pgDB := db.NewPostgresDB(cfg.DBURL)
	k8sClient := k8s.NewK8sClient(cfg.KubeConfig, cfg.K8SHost)
	kafkaSvc := service.NewKafkaService(cfg.KafkaBrokers)
	defer kafkaSvc.Close()

	// 3. Setup Architecture Layers
	repo := repository.NewProjectRepository(pgDB)
	k8sSvc := service.NewK8sService(k8sClient)
	buildSvc := service.NewBuildService(cfg.RegistryURL, k8sSvc)
	projectSvc := service.NewProjectService(repo, kafkaSvc, k8sSvc, buildSvc, cfg.BaseDomain)
	projectHandler := handler.NewProjectHandler(projectSvc)

	// 4. Start Kafka Consumer (Background Worker)
	kafkaSvc.StartConsumer(projectSvc.ProcessDeployment)

	// 5. Setup Gin Router
	r := gin.Default()

	// CORS middleware
	r.Use(cors.New(cors.Config{
		// AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001", "http://192.168.2.202:3000", "https://idp."},
		AllowAllOrigins: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		api.POST("/projects", projectHandler.CreateProject)
		api.GET("/projects", projectHandler.GetProjects)
		api.GET("/projects/:id", projectHandler.GetProjectByID)
		api.GET("/projects/:id/logs", projectHandler.GetProjectLogs)
	}

	// 6. Start Server
	log.Printf("Orchestrator starting on port %s...", cfg.ServerPort)
	if err := r.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}


