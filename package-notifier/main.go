package main

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	pb "github.com/snavarro/microtracker/package-notifier/proto"
	"google.golang.org/grpc"
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
	port := os.Getenv("NOTIFIER_PORT")
	if port == "" {
		port = "50051"
	}

	// Create TCP listener
	lis, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// Create gRPC server
	s := grpc.NewServer()
	pb.RegisterNotificationServiceServer(s, &notificationServer{})

	// Start server in a goroutine
	go func() {
		log.Printf("Starting notification server on port %s", port)
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")
	s.GracefulStop()
}
