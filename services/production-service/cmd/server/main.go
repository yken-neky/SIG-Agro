package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/sig-agro/api/proto/production"
)

type ProductionHandler struct {
	pb.UnimplementedProductionServiceServer
}

func (s *ProductionHandler) RecordActivity(ctx context.Context, req *pb.RecordActivityRequest) (*pb.RecordActivityResponse, error) {
	log.Printf("RecordActivity for parcel ID: %d, type: %s\n", req.ParcelId, req.ActivityType)
	return &pb.RecordActivityResponse{
		ActivityId: 1,
		Message:    "Activity recorded successfully",
	}, nil
}

func (s *ProductionHandler) GetActivity(ctx context.Context, req *pb.GetActivityRequest) (*pb.GetActivityResponse, error) {
	log.Printf("GetActivity for activity ID: %d\n", req.ActivityId)
	return &pb.GetActivityResponse{
		ActivityId:   req.ActivityId,
		ParcelId:     1,
		ActivityType: "harvest",
	}, nil
}

func (s *ProductionHandler) ListActivities(ctx context.Context, req *pb.ListActivitiesRequest) (*pb.ListActivitiesResponse, error) {
	log.Printf("ListActivities for parcel ID: %d\n", req.ParcelId)
	return &pb.ListActivitiesResponse{
		Total: 0,
	}, nil
}

func main() {
	port := 50054
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterProductionServiceServer(grpcServer, &ProductionHandler{})

	log.Printf("Production Service listening on port %d\n", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
