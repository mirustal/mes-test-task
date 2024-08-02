package consumer

import (
	"encoding/json"
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/IBM/sarama"
	"github.com/stretchr/testify/assert"

	"kafka-app/internal/adapters/kafka/producer"
)

var (
	testTopicPrefix = "test-practice-topic"
	testBroker      = "localhost:19092"
)

func testProduce(topic string, limit int) <-chan struct{} {
	var produceDone = make(chan struct{})

	p, err := sarama.NewAsyncProducer([]string{"localhost:19092"}, sarama.NewConfig())
	if err != nil {
		return nil
	}

	go func() {
		defer close(produceDone)

		for i := 0; i < limit; i++ {
			msg := producer.Message{fmt.Sprintf("%d", i), fmt.Sprintf("%d", i)}
			msgBytes, err := json.Marshal(msg)
			if err != nil {
				continue
			}
			select {
			case p.Input() <- &sarama.ProducerMessage{
				Topic: topic,
				Value: sarama.ByteEncoder(msgBytes),
			}:
			case err := <-p.Errors():
				fmt.Printf("Failed to send message to kafka, err: %s, msg: %s\n", err, msgBytes)
			}
		}
	}()

	return produceDone
}

func TestSyncConsumer(t *testing.T) {
	limit := 4500

	topic := testTopicPrefix + fmt.Sprintf("%d", time.Now().Unix())

	produceDone := testProduce(topic, 5000)

	var consumeMsgMap = make(map[string]struct{})
	var resChan = make(chan string)
	go func() {
		for r := range resChan {
			consumeMsgMap[r] = struct{}{}
		}
	}()

	handler := NewSyncConsumerGroupHandler(func(data []byte) error {
		var msg producer.Message
		err := json.Unmarshal(data, &msg)
		if err != nil {
			return err
		}
		resChan <- msg.Text
		return nil
	})
	consumer, err := NewConsumerGroup(testBroker, []string{topic}, fmt.Sprintf("%d", time.Now().Unix()), handler)
	if err != nil {
		return
	}
	defer consumer.Close()

	<-produceDone

	time.Sleep(1 * time.Second)

	log.Println(len(consumeMsgMap))
	assert.Equal(t, limit, len(consumeMsgMap))
}
