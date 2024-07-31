package producer

import (
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"

	"kafka-app/pkg/config"
)

type Producer struct {
	p   sarama.AsyncProducer
	cfg *config.Kafka
}

func NewProducer(cfg *config.Kafka) (*Producer, error) {
	broker := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	producer, err := sarama.NewAsyncProducer([]string{broker}, sarama.NewConfig())
	if err != nil {
		return nil, err
	}
	return &Producer{
		p:   producer,
		cfg: cfg,
	}, nil
}

type Message struct {
	Text string `json:"id"`
}

func (p *Producer) ProduceMessage(text string) error {
	msg := Message{text}
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal json fail")
	}

	select {
	case p.p.Input() <- &sarama.ProducerMessage{
		Topic: p.cfg.Topic,
		Value: sarama.ByteEncoder(msgBytes),
	}:
		return nil
	case err := <-p.p.Errors():
		return fmt.Errorf("Failed to send message to kafka, err: %s, msg: %s\n", err, msg.Text)
	}
}

func (p *Producer) Close() error {
	if p != nil {
		return p.p.Close()
	}
	return nil
}
