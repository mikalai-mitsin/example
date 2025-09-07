package handlers

import (
	"context"

	"github.com/IBM/sarama"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/log"
)

type ArticleHandler struct {
	articleUseCase articleUseCase
	logger         logger
}

func NewArticleHandler(articleUseCase articleUseCase, logger logger) *ArticleHandler {
	return &ArticleHandler{articleUseCase: articleUseCase, logger: logger}
}
func (h *ArticleHandler) Created(ctx context.Context, msg *sarama.ConsumerMessage) error {
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
func (h *ArticleHandler) Updated(ctx context.Context, msg *sarama.ConsumerMessage) error {
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
func (h *ArticleHandler) Deleted(ctx context.Context, msg *sarama.ConsumerMessage) error {
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
func (h *ArticleHandler) RegisterKafka(consumer *kafka.Consumer) error {
	consumer.AddHandler(
		kafka.NewHandler(
			"example.articles.article.created",
			"example.articles.article.created",
			h.Created,
		),
	)
	consumer.AddHandler(
		kafka.NewHandler(
			"example.articles.article.updated",
			"example.articles.article.updated",
			h.Updated,
		),
	)
	consumer.AddHandler(
		kafka.NewHandler(
			"example.articles.article.deleted",
			"example.articles.article.deleted",
			h.Deleted,
		),
	)
	return nil
}
