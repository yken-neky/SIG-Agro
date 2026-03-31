package http2

import (
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"

	
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"

	"github.com/sig-agro/services/user-service/domain"
	"github.com/sig-agro/services/user-service/internal/service"
)

type Server struct {
	Configuration *domain.Config
	UserService   *service.UserService
	server        *http.Server
}

func (s *Server) New() {
	r := gin.New()
	ginLogger := gin.Logger()
	if !s.Configuration.GinLogs {
		ginLogger = gin.LoggerWithConfig(gin.LoggerConfig{
			Output: io.Discard,
		})
	}
	r.Use(ginLogger)
	r.Use(gin.Recovery())

	RegisterHandlers(r, s.UserService)

	h2s := &http2.Server{}
	s.server = &http.Server{
		Addr:    fmt.Sprintf("0.0.0.0:%d", s.Configuration.Http2ListenPort),
		Handler: h2c.NewHandler(r, h2s),
	}
}

func (s *Server) Start() {
	if err := s.server.ListenAndServe(); err != nil {
		log.Fatal().Msgf("error: %v", err)
	}
}
