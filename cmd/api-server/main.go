package main

import (
	"3x-api/internal/config"
	"log/slog"
)

func main() {
	cfg := config.MustLoad()
	config.SetupLogger(cfg.Env)
	slog.Info("API-server has started")
	Bootstrap(cfg)

	// TODO: bootstrap VPN configuration

	// TODO: init storage: sqlite

	// TODO: init router: net/http

	// TODO: run server:

}
