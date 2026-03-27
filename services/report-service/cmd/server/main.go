package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	pb "github.com/sig-agro/api/proto/report"
)

type ReportHandler struct {
	pb.UnimplementedReportServiceServer
}

func (s *ReportHandler) GenerateReport(ctx context.Context, req *pb.GenerateReportRequest) (*pb.GenerateReportResponse, error) {
	log.Printf("GenerateReport for producer ID: %d, type: %s\n", req.ProducerId, req.ReportType)
	return &pb.GenerateReportResponse{
		ReportId: 1,
		Url:      "/reports/1.pdf",
		Message:  "Report generated successfully",
	}, nil
}

func (s *ReportHandler) GetReport(ctx context.Context, req *pb.GetReportRequest) (*pb.GetReportResponse, error) {
	log.Printf("GetReport for report ID: %d\n", req.ReportId)
	return &pb.GetReportResponse{
		ReportId:   req.ReportId,
		ProducerId: 1,
		ReportType: "summary",
		Url:        "/reports/1.pdf",
	}, nil
}

func (s *ReportHandler) ListReports(ctx context.Context, req *pb.ListReportsRequest) (*pb.ListReportsResponse, error) {
	log.Printf("ListReports for producer ID: %d\n", req.ProducerId)
	return &pb.ListReportsResponse{
		Total: 0,
	}, nil
}

func main() {
	port := 50057
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterReportServiceServer(grpcServer, &ReportHandler{})

	log.Printf("Report Service listening on port %d\n", port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
