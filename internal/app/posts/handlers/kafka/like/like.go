package handlers

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
)

type LikeHandler struct {
	likeUseCase likeUseCase
	logger      logger
}

func NewLikeHandler(likeUseCase likeUseCase, logger logger) *LikeHandler {
	return &LikeHandler{likeUseCase: likeUseCase, logger: logger}
}
func (h *LikeHandler) Created(ctx context.Context, msg *sarama.ConsumerMessage) error {
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
func (h *LikeHandler) Updated(ctx context.Context, msg *sarama.ConsumerMessage) error {
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
func (h *LikeHandler) Deleted(ctx context.Context, msg *sarama.ConsumerMessage) error {
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
func (h *LikeHandler) RegisterKafka(consumer *kafka.Consumer) error {
	consumer.AddHandler(
		kafka.NewHandler("example.posts.like.created", "example.posts.like.created", h.Created),
	)
	consumer.AddHandler(
		kafka.NewHandler("example.posts.like.updated", "example.posts.like.updated", h.Updated),
	)
	consumer.AddHandler(
		kafka.NewHandler("example.posts.like.deleted", "example.posts.like.deleted", h.Deleted),
	)
	return nil
}
