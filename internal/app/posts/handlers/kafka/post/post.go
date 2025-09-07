package handlers

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
)

type PostHandler struct {
	postUseCase postUseCase
	logger      logger
}

func NewPostHandler(postUseCase postUseCase, logger logger) *PostHandler {
	return &PostHandler{postUseCase: postUseCase, logger: logger}
}
func (h *PostHandler) Created(ctx context.Context, msg *sarama.ConsumerMessage) error {
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
func (h *PostHandler) Updated(ctx context.Context, msg *sarama.ConsumerMessage) error {
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
func (h *PostHandler) Deleted(ctx context.Context, msg *sarama.ConsumerMessage) error {
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
func (h *PostHandler) RegisterKafka(consumer *kafka.Consumer) error {
	consumer.AddHandler(
		kafka.NewHandler("example.posts.post.created", "example.posts.post.created", h.Created),
	)
	consumer.AddHandler(
		kafka.NewHandler("example.posts.post.updated", "example.posts.post.updated", h.Updated),
	)
	consumer.AddHandler(
		kafka.NewHandler("example.posts.post.deleted", "example.posts.post.deleted", h.Deleted),
	)
	return nil
}
