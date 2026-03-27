package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/sig-agro/api/proto/notification"
)

type NotificationHandler struct {
	pb.UnimplementedNotificationServiceServer
}

func (s *NotificationHandler) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	log.Printf("SendNotification for user ID: %d, channel: %s\n", req.UserId, req.Channel)
	return &pb.SendNotificationResponse{
		NotificationId: 1,
		Message:        "Notification sent successfully",
	}, nil
}

func (s *NotificationHandler) GetNotification(ctx context.Context, req *pb.GetNotificationRequest) (*pb.GetNotificationResponse, error) {
	log.Printf("GetNotification for notification ID: %d\n", req.NotificationId)
	return &pb.GetNotificationResponse{
		NotificationId: req.NotificationId,
		UserId:         1,
		Title:          "Test",
		Message:        "Test notification",
	}, nil
}

func (s *NotificationHandler) ListNotifications(ctx context.Context, req *pb.ListNotificationsRequest) (*pb.ListNotificationsResponse, error) {
	log.Printf("ListNotifications for user ID: %d\n", req.UserId)
	return &pb.ListNotificationsResponse{
		Total: 0,
	}, nil
}

func (s *NotificationHandler) MarkAsRead(ctx context.Context, req *pb.MarkAsReadRequest) (*pb.MarkAsReadResponse, error) {
	log.Printf("MarkAsRead for notification ID: %d\n", req.NotificationId)
	return &pb.MarkAsReadResponse{
		Success: true,
		Message: "Notification marked as read",
	}, nil
}

func main() {
	port := 50056
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterNotificationServiceServer(grpcServer, &NotificationHandler{})

	log.Printf("Notification Service listening on port %d\n", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
