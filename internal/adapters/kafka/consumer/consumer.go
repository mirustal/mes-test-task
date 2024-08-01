package consumer

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/IBM/sarama"

	"kafka-app/internal/adapters/db/postgres"
	"kafka-app/internal/adapters/kafka/producer"
	"kafka-app/pkg/config"
)

type ConsumerGroupHandler interface {
	sarama.ConsumerGroupHandler
	WaitReady()
	Reset()
}

type ConsumerGroup struct {
	cg sarama.ConsumerGroup
}


func NewConsumerGroup(broker string, topics []string, group string, handler ConsumerGroupHandler) (*ConsumerGroup, error) {
	ctx := context.Background()
	cfg := sarama.NewConfig()
	cfg.Consumer.Offsets.AutoCommit.Enable = true
	cfg.Consumer.Offsets.AutoCommit.Interval = 5 * time.Second
	
	client, err := sarama.NewConsumerGroup([]string{broker}, group, cfg)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer group: %w", err)
	}

	go func() {
		for {
			fmt.Println("Starting to consume messages")
			if err := client.Consume(ctx, topics, handler); err != nil {
				if err == sarama.ErrClosedConsumerGroup {
					break
				} else {
					fmt.Printf("Error from consumer: %v\n", err)
					time.Sleep(time.Second)
				}
			}
			if ctx.Err() != nil {
				return
			}
			handler.Reset()
		}
	}()

	handler.WaitReady()

	return &ConsumerGroup{
		cg: client,
	}, nil
}

func (c *ConsumerGroup) Close() error {
	return c.cg.Close()
}

func decodeMessage(data []byte) (*producer.Message, error) {
	var msg producer.Message
	err := json.Unmarshal(data, &msg)
	if err != nil {
		return nil, err
	}
	return &msg, nil
}

func StartSyncConsumer(cfg *config.Kafka, db *postgres.PostgresMessageRep) (*ConsumerGroup, error) {
	handler := NewSyncConsumerGroupHandler(func(data []byte) error {
		msg, err := decodeMessage(data)
		if err != nil {
			return fmt.Errorf("failed to decode message: %w", err)
		}
		go func(id string) { 
			time.Sleep(5 * time.Second)
			err := db.MarkAsRead(context.Background(), id)
			if err != nil {
				log.Printf("Failed to mark message as read: %v", err)
			}
		}(msg.ID)

		fmt.Printf("Consumed message: %v\n", msg)
		return nil
	})

	consumer, err := NewConsumerGroup(fmt.Sprintf("%v:%v", cfg.Host, cfg.Port), []string{cfg.Topic}, fmt.Sprintf("sync-consumer-%d", time.Now().Unix()), handler)
	if err != nil {
		return nil, fmt.Errorf("failed to start consumer group: %w", err)
	}
	return consumer, nil
}