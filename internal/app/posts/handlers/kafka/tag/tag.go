package handlers

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
)

type TagHandler struct {
	tagUseCase tagUseCase
	logger     logger
}

func NewTagHandler(tagUseCase tagUseCase, logger logger) *TagHandler {
	return &TagHandler{tagUseCase: tagUseCase, logger: logger}
}
func (h *TagHandler) Created(ctx context.Context, msg *sarama.ConsumerMessage) error {
	logger := h.logger.WithContext(ctx)
	logger.Info(
		"received created message",
		log.String("topic", msg.Topic),
		log.Int32("partition", msg.Partition),
		log.Int64("offset", msg.Offset),
		log.String("key", string(msg.Key)),
		log.String("value", string(msg.Value)),
	)
	return nil
}
func (h *TagHandler) Updated(ctx context.Context, msg *sarama.ConsumerMessage) error {
	logger := h.logger.WithContext(ctx)
	logger.Info(
		"received updated message",
		log.String("topic", msg.Topic),
		log.Int32("partition", msg.Partition),
		log.Int64("offset", msg.Offset),
		log.String("key", string(msg.Key)),
		log.String("value", string(msg.Value)),
	)
	return nil
}
func (h *TagHandler) Deleted(ctx context.Context, msg *sarama.ConsumerMessage) error {
	logger := h.logger.WithContext(ctx)
	logger.Info(
		"received deleted message",
		log.String("topic", msg.Topic),
		log.Int32("partition", msg.Partition),
		log.Int64("offset", msg.Offset),
		log.String("key", string(msg.Key)),
		log.String("value", string(msg.Value)),
	)
	return nil
}
func (h *TagHandler) RegisterKafka(consumer *kafka.Consumer) error {
	consumer.AddHandler(
		kafka.NewHandler("example.posts.tag.created", "example.posts.tag.created", h.Created),
	)
	consumer.AddHandler(
		kafka.NewHandler("example.posts.tag.updated", "example.posts.tag.updated", h.Updated),
	)
	consumer.AddHandler(
		kafka.NewHandler("example.posts.tag.deleted", "example.posts.tag.deleted", h.Deleted),
	)
	return nil
}
