package app

import (
	_ "github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
	"github.com/qPyth/mobydev-internship-auth/internal/config"
	"github.com/qPyth/mobydev-internship-auth/internal/server"
	"github.com/qPyth/mobydev-internship-auth/internal/services"
	"github.com/qPyth/mobydev-internship-auth/internal/storage/sqlite"
	"github.com/qPyth/mobydev-internship-auth/internal/transport/http"
	"github.com/qPyth/mobydev-internship-auth/pkg/auth"
	"log/slog"
	"os"
)

var jwtSecret = os.Getenv("JWT_SECRET")

func Run(cfg *config.Config) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	storage := sqlite.New(cfg.StoragePath)
	tokenManager := auth.NewManager(jwtSecret, cfg.TokenTTL)
	userService := services.NewUserService(storage, tokenManager)

	h := http.NewHandler(logger, userService)

	srv := server.New(cfg, h.Init())
	logger.Info("starting server on port: ", "port", cfg.Port)
	if err := srv.Run(); err != nil {
		logger.Error("failed to run server: ", "error", err.Error())
	}

}
