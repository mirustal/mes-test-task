package consumer

import "github.com/IBM/sarama"

type syncConsumerGroupHandler struct {
	ready chan bool

	cb func([]byte) error
}

func NewSyncConsumerGroupHandler(cb func([]byte) error) ConsumerGroupHandler {
	handler := syncConsumerGroupHandler{
		ready: make(chan bool, 0),
		cb:    cb,
	}
	return &handler
}

func (h *syncConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(h.ready)
	return nil
}


func (h *syncConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *syncConsumerGroupHandler) WaitReady() {
	<-h.ready
	return
}

func (h *syncConsumerGroupHandler) Reset() {
	h.ready = make(chan bool, 0)
	return
}

func (h *syncConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	claimMsgChan := claim.Messages()

	for message := range claimMsgChan {
		if h.cb(message.Value) == nil {
			session.MarkMessage(message, "")
		}
	}

	return nil
}