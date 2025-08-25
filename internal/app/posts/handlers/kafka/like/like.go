package handlers

import (
	"github.com/IBM/sarama"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
)

type LikeHandler struct {
	likeUseCase likeUseCase
	logger      logger
}

func NewLikeHandler(likeUseCase likeUseCase, logger logger) *LikeHandler {
	return &LikeHandler{likeUseCase: likeUseCase, logger: logger}
}
func (h *LikeHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}
func (h *LikeHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h *LikeHandler) ConsumeClaim(
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
