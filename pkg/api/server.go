package api

import (
	"fmt"
	"net/http"
	"time"

	"github/Igo87/crypt/config"
	"github/Igo87/crypt/pkg/logger"
)

//go:generate mockgen -source=server.go -destination=mocks/mock_server.go -package=mocks
type Server struct {
	server *http.Server
}

// ServeHTTP serves the HTTP request by directly calling the handler.
func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.server.Handler.ServeHTTP(w, r)
}

// New returns a new instance of the Server interface.
func New(handler http.Handler, addr string) *Server {
	return &Server{
		server: &http.Server{
			Addr:              addr,
			Handler:           handler,
			ReadHeaderTimeout: time.Duration(config.Cfg.Waiting) * time.Second,
			IdleTimeout:       time.Duration(config.Cfg.Waiting) * time.Second,
			MaxHeaderBytes:    1 << 20,
			WriteTimeout:      time.Duration(config.Cfg.Waiting) * time.Second,
		},
	}
}

func (s *Server) Start() error {
	msgErr := fmt.Errorf("failed to start server: %w", http.ErrServerClosed)

	go func() {
		if err := s.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.LogStart().Error("Failed to start server: %v", msgErr)
		}
	}()
	return nil
}
