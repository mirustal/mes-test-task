package main

import (
	"context"
	"log"
	"net/http"

	"kafka-app/internal/adapters/db/postgres"
	"kafka-app/internal/http-server/interfaces/routes"
	"kafka-app/pkg/config"
	"kafka-app/pkg/logger"
)

func main() {
	ctx := context.Background()

	cfg, err := config.LoadConfig("config", "yaml")
	if err != nil {
		log.Fatal("fail load config: %v", err)
	}

	logger.LogInit("debug")
	r := routes.InitRoutes()
	db, err := postgres.NewMR(ctx, cfg.DB)
	log.Fatal(http.ListenAndServe(":8081", r))
}
