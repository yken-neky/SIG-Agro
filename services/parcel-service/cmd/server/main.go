package main

import (
	"fmt"
	"log"
	"net"

	"github.com/sig-agro/services/parcel-service/internal/config"
	"github.com/sig-agro/services/parcel-service/internal/handler"
	"github.com/sig-agro/services/parcel-service/internal/repository"
	"google.golang.org/grpc"

	pb "github.com/sig-agro/api/proto/parcel"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Initialize database
	db, err := repository.NewDatabase(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize repository
	repo := repository.NewRepository(db)

	// Initialize handler
	svc := handler.NewParcelService(repo)

	// Start gRPC server
	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", cfg.Port))
	if err != nil {
		log.Fatalf("Failed to listen on port %d: %v", cfg.Port, err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterParcelServiceServer(grpcServer, svc)

	log.Printf("Parcel Service listening on port %d\n", cfg.Port)
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
