package main

import (
    "context"
    "flag"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/gin-gonic/gin"

    "github.com/crudbox/crudbox/internal/config"
    "github.com/crudbox/crudbox/internal/database"
    "github.com/crudbox/crudbox/internal/handler"
    "github.com/crudbox/crudbox/internal/middleware"
    "github.com/crudbox/crudbox/internal/repository"
    "github.com/crudbox/crudbox/internal/service"
    "github.com/crudbox/crudbox/internal/utils"
)

func main() {
    migrateOnly := flag.Bool("migrate-only", false, "apply database migrations and exit")
    flag.Parse()

    cfg := config.Load()

    db, err := database.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatalf("failed to connect to database: %v", err)
    }
    defer db.Close()

    if err := database.RunMigrations(db); err != nil {
        log.Fatalf("failed to run database migrations: %v", err)
    }

    if *migrateOnly {
        log.Println("database migrations completed")
        return
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

    srv := &http.Server{
        Addr:    ":" + cfg.ServerPort,
        Handler: router,
    }

    serverErrors := make(chan error, 1)
    go func() {
        log.Printf("Crudbox API listening on %s", srv.Addr)
        if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
            serverErrors <- err
        }
    }()

    shutdownSignals := make(chan os.Signal, 1)
    signal.Notify(shutdownSignals, syscall.SIGINT, syscall.SIGTERM)

    select {
    case err := <-serverErrors:
        log.Fatalf("server error: %v", err)
    case sig := <-shutdownSignals:
        log.Printf("received signal %s, shutting down...", sig)
    }

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("failed to shutdown server gracefully: %v", err)
    }

    log.Println("server stopped")
}
