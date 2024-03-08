package server

import (
	"context"
	"github.com/qPyth/mobydev-internship-auth/internal/config"
	"net/http"
)

type Server struct {
	httpServer *http.Server
}

func New(cfg *config.Config, handler http.Handler) *Server {
	return &Server{
		httpServer: &http.Server{
			Addr:         cfg.HTTP.Host + ":" + cfg.HTTP.Port,
			Handler:      handler,
			ReadTimeout:  cfg.HTTP.ReadTimeOut,
			WriteTimeout: cfg.HTTP.WriteTimeOut,
		},
	}
}

func (s *Server) Run() error {
	return s.httpServer.ListenAndServe()
}

func (s *Server) Stop(ctx context.Context) error {
	return s.httpServer.Shutdown(ctx)
}
