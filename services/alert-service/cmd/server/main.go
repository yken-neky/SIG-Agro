package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/sig-agro/api/proto/alert"
)

type AlertHandler struct {
	pb.UnimplementedAlertServiceServer
}

func (s *AlertHandler) CreateAlert(ctx context.Context, req *pb.CreateAlertRequest) (*pb.CreateAlertResponse, error) {
	log.Printf("CreateAlert for parcel ID: %d, type: %s\n", req.ParcelId, req.AlertType)
	return &pb.CreateAlertResponse{
		AlertId: 1,
		Message: "Alert created successfully",
	}, nil
}

func (s *AlertHandler) GetAlert(ctx context.Context, req *pb.GetAlertRequest) (*pb.GetAlertResponse, error) {
	log.Printf("GetAlert for alert ID: %d\n", req.AlertId)
	return &pb.GetAlertResponse{
		AlertId:   req.AlertId,
		ParcelId:  1,
		AlertType: "weather",
		Severity:  "medium",
	}, nil
}

func (s *AlertHandler) ListAlerts(ctx context.Context, req *pb.ListAlertsRequest) (*pb.ListAlertsResponse, error) {
	log.Printf("ListAlerts for parcel ID: %d\n", req.ParcelId)
	return &pb.ListAlertsResponse{
		Total: 0,
	}, nil
}

func (s *AlertHandler) EvaluateAlerts(ctx context.Context, req *pb.EvaluateAlertsRequest) (*pb.EvaluateAlertsResponse, error) {
	log.Printf("EvaluateAlerts\n")
	return &pb.EvaluateAlertsResponse{
		Evaluated: 0,
		Triggered: 0,
		Message:   "Alerts evaluated",
	}, nil
}

func main() {
	port := 50055
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterAlertServiceServer(grpcServer, &AlertHandler{})

	log.Printf("Alert Service listening on port %d\n", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
