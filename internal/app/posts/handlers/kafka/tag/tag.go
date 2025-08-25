package handlers

import (
	"github.com/IBM/sarama"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
)

type TagHandler struct {
	tagUseCase tagUseCase
	logger     logger
}

func NewTagHandler(tagUseCase tagUseCase, logger logger) *TagHandler {
	return &TagHandler{tagUseCase: tagUseCase, logger: logger}
}
func (h *TagHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (h *TagHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *TagHandler) ConsumeClaim(
	session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		logger := h.logger
		logger.Info(
			"received message",
			log.String("topic", msg.Topic),
			log.Int32("partition", msg.Partition),
			log.Int64("offset", msg.Offset),
			log.String("key", string(msg.Key)),
			log.String("value", string(msg.Value)),
		)
		session.MarkMessage(msg, "")
	}
	return nil
}
