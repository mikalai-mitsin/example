package handlers

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
)

const (
	topicName = "example.posts.post.v1"
	groupID   = "example.posts.post"
)

type PostHandler struct {
	postUseCase postUseCase
	logger      logger
}

func NewPostHandler(postUseCase postUseCase, logger logger) *PostHandler {
	return &PostHandler{postUseCase: postUseCase, logger: logger}
}
func (h *PostHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
	logger := h.logger.WithContext(ctx)
	logger.Info(
		"received message",
		log.String("topic", msg.Topic),
		log.Int32("partition", msg.Partition),
		log.Int64("offset", msg.Offset),
		log.String("key", string(msg.Key)),
		log.String("value", string(msg.Value)),
	)
	return nil
}
func (h *PostHandler) RegisterKafka(consumer *kafka.Consumer) error {
	consumer.AddHandler(kafka.NewHandler(topicName, groupID, h.Handle))
	return nil
}
