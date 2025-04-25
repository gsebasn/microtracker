package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sns/types"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	pb "github.com/snavarro/microtracker/package-notifier/proto"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// NotificationHandler defines the interface for different notification channels
type NotificationHandler interface {
	Send(ctx context.Context, status *pb.PackageStatus, config map[string]string) error
}

// EmailNotificationHandler handles email notifications
type EmailNotificationHandler struct{}

func (h *EmailNotificationHandler) Send(ctx context.Context, status *pb.PackageStatus, config map[string]string) error {
	// TODO: Implement email sending logic
	log.Printf("Sending email notification for package %s to %s", status.PackageId, config["recipient"])
	return nil
}

// SMSNotificationHandler handles SMS notifications
type SMSNotificationHandler struct{}

func (h *SMSNotificationHandler) Send(ctx context.Context, status *pb.PackageStatus, config map[string]string) error {
	// TODO: Implement SMS sending logic
	log.Printf("Sending SMS notification for package %s to %s", status.PackageId, config["phone"])
	return nil
}

// QueueNotificationHandler handles queue notifications
type QueueNotificationHandler struct{}

func (h *QueueNotificationHandler) Send(ctx context.Context, status *pb.PackageStatus, config map[string]string) error {
	// TODO: Implement queue publishing logic
	log.Printf("Publishing notification for package %s to queue %s", status.PackageId, config["queue"])
	return nil
}

// LogNotificationHandler handles log notifications
type LogNotificationHandler struct{}

func (h *LogNotificationHandler) Send(ctx context.Context, status *pb.PackageStatus, config map[string]string) error {
	log.Printf("Logging notification for package %s: %s", status.PackageId, status.Status)
	return nil
}

type notificationServer struct {
	pb.UnimplementedNotificationServiceServer
	snsClient *sns.Client
	topicArn  string
}

func newNotificationServer() (*notificationServer, error) {
	// Load AWS configuration
	cfg, err := config.LoadDefaultConfig(context.Background())
	if err != nil {
		return nil, fmt.Errorf("unable to load AWS SDK config: %v", err)
	}

	// Create SNS client
	snsClient := sns.NewFromConfig(cfg)

	// Get topic ARN from environment variable
	topicArn := os.Getenv("AWS_SNS_TOPIC_ARN")
	if topicArn == "" {
		return nil, fmt.Errorf("AWS_SNS_TOPIC_ARN environment variable is required")
	}

	return &notificationServer{
		snsClient: snsClient,
		topicArn:  topicArn,
	}, nil
}

func (s *notificationServer) NotifyPackageStatus(ctx context.Context, req *pb.NotificationRequest) (*pb.NotificationResponse, error) {
	log.Printf("Received notification for package %s: %s", req.Status.PackageId, req.Status.Status)

	// Create notification message
	message := struct {
		PackageStatus *pb.PackageStatus         `json:"package_status"`
		Channels      []*pb.NotificationChannel `json:"channels"`
		TemplateID    string                    `json:"template_id,omitempty"`
		Metadata      map[string]string         `json:"metadata,omitempty"`
	}{
		PackageStatus: req.Status,
		Channels:      req.Channels,
		TemplateID:    req.TemplateId,
		Metadata:      req.Metadata,
	}

	// Convert message to JSON
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal notification message: %v", err)
	}

	// Publish message to SNS topic
	input := &sns.PublishInput{
		TopicArn: aws.String(s.topicArn),
		Message:  aws.String(string(messageBytes)),
		MessageAttributes: map[string]types.MessageAttributeValue{
			"package_id": {
				DataType:    aws.String("String"),
				StringValue: aws.String(req.Status.PackageId),
			},
			"status": {
				DataType:    aws.String("String"),
				StringValue: aws.String(req.Status.Status),
			},
		},
	}

	result, err := s.snsClient.Publish(ctx, input)
	if err != nil {
		return nil, fmt.Errorf("failed to publish notification: %v", err)
	}

	return &pb.NotificationResponse{
		Success: true,
		Message: "Notification published successfully",
		Results: []*pb.NotificationResult{
			{
				ChannelType: "sns",
				Success:     true,
				Message:     fmt.Sprintf("Message published with ID: %s", *result.MessageId),
			},
		},
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

	// Create notification server
	server, err := newNotificationServer()
	if err != nil {
		log.Fatalf("failed to create notification server: %v", err)
	}

	// Create gRPC server
	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, server)

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
