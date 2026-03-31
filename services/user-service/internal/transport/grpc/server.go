package grpc

import (
	"fmt"
	"net"

	"github.com/rs/zerolog/log"
	pb "github.com/sig-agro/api/proto/user"
	"google.golang.org/grpc"
)

func (s *UserGRPCServer) Serve() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Port))
	if err != nil {
		log.Fatal().Msgf("failed to listen: %v", err)
	}
	grpcs := grpc.NewServer()
	pb.RegisterUserServiceServer(grpcs, s)
	log.Info().Msgf("Starting gRPC user service in port: %v", s.Port)
	if err := grpcs.Serve(lis); err != nil {
		log.Fatal().Msgf("failed to serve: %v", err)
	}
}
