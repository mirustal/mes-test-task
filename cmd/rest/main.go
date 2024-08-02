package main
// @title Message API
// @version 1.0
// @description This is a sample server for the Message API.
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
	kafkaProduce, err := producer.NewProducer(cfg.Kafka)
	if err != nil {
		log.Fatal(err)
	}
	defer kafkaProduce.Close()

	// _, consumerHandler, err := consumer.StartSyncConsumer(cfg.Kafka)
	// if err != nil {
	// 	log.Fatal("fail load kafka %v", err)
	// }

	r := routes.InitRoutes(rep, kafkaProduce)
	log.Println("Start serve in port ", cfg.ServerPort)
	log.Fatal(http.ListenAndServe(":"+cfg.ServerPort, r))
}
