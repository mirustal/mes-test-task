package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"kafka-app/internal/adapters/db/postgres"
	"kafka-app/internal/adapters/kafka/consumer"
	"kafka-app/pkg/config"
)

func main() {
	cfg, err := config.LoadConfig("config", "yaml")
	if err != nil {
		log.Fatal("fail load config %v", err)
	}

		
	db, err := postgres.NewMR(context.Background(), cfg.DB)

	kafka, err := consumer.StartSyncConsumer(cfg.Kafka, db)
	if err != nil {
		log.Fatal("fail load kafka %v", err)
	}
	defer kafka.Close()




	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	fmt.Println("recieved signal", <-c)
}