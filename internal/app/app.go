package app

import (
	"github.com/qPyth/mobydev-internship-auth/internal/config"
	"github.com/qPyth/mobydev-internship-auth/internal/server"
	"github.com/qPyth/mobydev-internship-auth/internal/transport/http"
	"log/slog"
	"os"
)

func Run(cfg *config.Config) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	h := http.NewHandler(logger)
	srv := server.New(cfg, h.Init())

}
