package main

import (
	"github.com/qPyth/mobydev-internship-auth/internal/app"
	"github.com/qPyth/mobydev-internship-auth/internal/config"
)

func main() {
	cfg := config.Load()
	app.Run(cfg)
}
