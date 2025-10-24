package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/ctonew/mockapi/internal/config"
	"github.com/ctonew/mockapi/internal/database"
	"github.com/ctonew/mockapi/internal/handler"
	"github.com/ctonew/mockapi/internal/middleware"
	"github.com/ctonew/mockapi/internal/repository"
	"github.com/ctonew/mockapi/internal/service"
	"github.com/ctonew/mockapi/internal/utils"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	defer db.Close()

	if err := database.RunMigrations(db); err != nil {
		log.Fatalf("failed to run database migrations: %v", err)
	}

	userRepo := repository.NewUserRepository(db)
	organisationRepo := repository.NewOrganisationRepository(db)
	projectRepo := repository.NewProjectRepository(db)
	endpointRepo := repository.NewEndpointRepository(db)

	passwordManager := utils.NewPasswordManager()
	tokenManager := utils.NewTokenManager(cfg.JWTSecret)

	authService := service.NewAuthService(userRepo, passwordManager, tokenManager)
	organisationService := service.NewOrganisationService(organisationRepo)
	projectService := service.NewProjectService(projectRepo, organisationRepo)
	endpointService := service.NewEndpointService(endpointRepo, projectRepo)

	router := gin.Default()

	api := router.Group("/api/v1")
	authHandler := handler.NewAuthHandler(authService)
	authHandler.RegisterRoutes(api.Group("/auth"))

	authMiddleware := middleware.NewAuthMiddleware(tokenManager, userRepo)
	protected := api.Group("")
	protected.Use(authMiddleware.Handler())

	organisationHandler := handler.NewOrganisationHandler(organisationService)
	organisationHandler.RegisterRoutes(protected)

	projectHandler := handler.NewProjectHandler(projectService)
	projectHandler.RegisterRoutes(protected)

	endpointHandler := handler.NewEndpointHandler(endpointService)
	endpointHandler.RegisterRoutes(protected)

	mockHandler := handler.NewMockHandler(endpointService)
	router.Any("/:projectCode/*endpointPath", mockHandler.Handle)

	if err := router.Run(":" + cfg.ServerPort); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
