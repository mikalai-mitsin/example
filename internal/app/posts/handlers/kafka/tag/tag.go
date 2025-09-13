package handlers

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
)

const (
	topicName = "example.posts.tag.v1"
	groupID   = "example.posts.tag"
)

type TagHandler struct {
	tagUseCase tagUseCase
	logger     logger
}

func NewTagHandler(tagUseCase tagUseCase, logger logger) *TagHandler {
	return &TagHandler{tagUseCase: tagUseCase, logger: logger}
}
func (h *TagHandler) Handle(ctx context.Context, msg *sarama.ConsumerMessage) error {
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
func (h *TagHandler) RegisterKafka(consumer *kafka.Consumer) error {
	consumer.AddHandler(kafka.NewHandler(topicName, groupID, h.Handle))
	return nil
}
