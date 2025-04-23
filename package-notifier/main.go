package main

import (
	"context"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/snavarro/microtracker/package-notifier/proto"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type notificationServer struct {
	pb.UnimplementedNotificationServiceServer
}

func (s *notificationServer) NotifyPackageStatus(ctx context.Context, req *pb.NotificationRequest) (*pb.NotificationResponse, error) {
	log.Printf("Received notification for package %s: %s", req.Status.PackageId, req.Status.Status)

	// Here you would implement actual notification logic
	// For now, we'll just log the notification
	for _, channel := range req.NotificationChannels {
		log.Printf("Sending notification via %s", channel)
	}

	return &pb.NotificationResponse{
		Success: true,
		Message: "Notification sent successfully",
	}, nil
}

func (s *notificationServer) GetNotificationHistory(ctx context.Context, req *pb.GetHistoryRequest) (*pb.NotificationHistory, error) {
	log.Printf("Getting notification history for package %s", req.PackageId)

	// Here you would implement actual history retrieval logic
	// For now, we'll return an empty history
	return &pb.NotificationHistory{
		Statuses: []*pb.PackageStatus{},
	}, nil
}

func main() {
	// Get port from environment variable or use default
	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "50051"
	}

	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8080"
	}

	// Create TCP listener for gRPC
	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, &notificationServer{})

	// Start gRPC server in a goroutine
	go func() {
		log.Printf("Starting gRPC server on port %s", grpcPort)
		if err := grpcServer.Serve(lis); err != nil {
			log.Fatalf("failed to serve gRPC: %v", err)
		}
	}()

	// Create HTTP server for gRPC-Gateway and Swagger UI
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Create gRPC-Gateway mux
	gwmux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	if err := pb.RegisterNotificationServiceHandlerFromEndpoint(ctx, gwmux, "localhost:"+grpcPort, opts); err != nil {
		log.Fatalf("failed to register gateway: %v", err)
	}

	// Create Gin router
	router := gin.Default()

	// Add middleware for request logging
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// Swagger UI endpoint
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Forward all other requests to gRPC-Gateway
	router.NoRoute(func(c *gin.Context) {
		gwmux.ServeHTTP(c.Writer, c.Request)
	})

	// Create HTTP server
	httpServer := &http.Server{
		Addr:    ":" + httpPort,
		Handler: router,
	}

	// Start HTTP server in a goroutine
	go func() {
		log.Printf("Starting HTTP server on port %s", httpPort)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("failed to serve HTTP: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the servers
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down servers...")

	// Shutdown HTTP server
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Fatalf("failed to shutdown HTTP server: %v", err)
	}

	// Shutdown gRPC server
	grpcServer.GracefulStop()
}
