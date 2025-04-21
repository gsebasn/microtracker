package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/snavarro/microtracker/config"
	"github.com/snavarro/microtracker/internal/handler"
	"github.com/snavarro/microtracker/internal/repository/mongo"
	"github.com/snavarro/microtracker/internal/service"
)

// @title Package Tracking API
// @version 1.0
// @description A microservice for tracking packages
// @host localhost:8080
// @BasePath /api/v1
func main() {
	// Initialize configuration
	cfg := config.NewConfig()

	// Connect to MongoDB
	db, err := config.ConnectDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	// Set Gin mode
	if os.Getenv("GIN_MODE") != "debug" {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize repositories
	packageRepo := mongo.NewPackageRepository(db)

	// Initialize services
	packageService := service.NewPackageService(packageRepo)

	// Initialize handlers
	packageHandler := handler.NewPackageHandler(packageService)

	// Initialize router
	router := gin.Default()

	// Add middleware for request logging
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// API routes
	api := router.Group("/api/v1")
	{
		packages := api.Group("/packages")
		{
			packages.GET("", packageHandler.ListPackages)
			packages.GET("/search", packageHandler.SearchPackages)
			packages.GET("/:id", packageHandler.GetPackage)
			packages.POST("", packageHandler.CreatePackage)
			packages.PUT("/:id", packageHandler.UpdatePackage)
			packages.DELETE("/:id", packageHandler.DeletePackage)
		}
	}

	// Create HTTP server
	srv := &http.Server{
		Addr:    cfg.ServerAddress,
		Handler: router,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("Server starting on %s", cfg.ServerAddress)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	log.Println("Server exiting")
}
