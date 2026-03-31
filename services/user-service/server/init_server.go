package server

import (
	"github.com/sig-agro/services/user-service/domain"
	"github.com/sig-agro/services/user-service/internal/repository/database"
	"github.com/sig-agro/services/user-service/internal/service"
	"github.com/sig-agro/services/user-service/internal/transport/grpc"
	"github.com/sig-agro/services/user-service/internal/transport/http2"
)

func RunServer(configFilePath string) {
	conf := domain.New(configFilePath)

	dbConn := InitPostgresDb(conf)
	rdClient := InitRedis(conf)

	defer dbConn.Close()

	postgresRepo := database.NewPostgresRepository(dbConn)
	cacheRepo := database.NewRedisRepository(rdClient)

	userService := service.NewUserService(postgresRepo, cacheRepo, conf.JWTSecret)

	grpcServer := &grpc.UserGRPCServer{
		Port:    conf.ListenPort,
		Service: userService,
	}

	http2s := http2.Server{
		Configuration: conf,
		UserService:   userService,
	}
	http2s.New()
	go http2s.Start()

	grpcServer.Serve()
}
