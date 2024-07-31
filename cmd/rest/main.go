package main

import (
	"context"
	"log"
	"net/http"

	"kafka-app/internal/adapters/db/postgres"
	"kafka-app/internal/adapters/kafka/producer"
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

	logger.LogInit(cfg.ModeLog)
	rep, err := postgres.NewMR(ctx, cfg.DB)
	if err != nil {
		log.Fatal(err)
	}

	broker := cfg.Kafka.Host + ":" + cfg.Kafka.Port
	kafka, err := producer.NewProducer(broker)

	r := routes.InitRoutes(rep, kafka)
	log.Println("Start serve in port %v", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":" + cfg.ServerPort, r))
}
