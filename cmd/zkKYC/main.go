package main

import (
	"zkKYC-backend/internal/app/config"
	"zkKYC-backend/internal/app/server"
)

var cfg config.Config

func main() {

	cfg.Init()

	s := server.NewServer(cfg)
	s.Start(cfg)
}
