package main

import (
	"sfit-platform-web-backend/internal/config"
	"sfit-platform-web-backend/internal/server"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	server := server.NewServer(cfg)

	if err := server.Run(); err != nil {
		panic(err)
	}
}
