package events

import (
	"context"
	"encoding/json"

	entities "github.com/mikalai-mitsin/example/internal/app/articles/entities/article"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

const (
	topicEventCreated = "example.articles.article.created"
	topicEventUpdated = "example.articles.article.updated"
	topicEventDeleted = "example.articles.article.deleted"
)

type ArticleEventProducer struct {
	producer producer
	logger   logger
}

func NewArticleEventProducer(
	producer producer,
	logger logger,
) *ArticleEventProducer {
	return &ArticleEventProducer{producer: producer, logger: logger}
}

func (p *ArticleEventProducer) Created(ctx context.Context, article entities.Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}
	message := &kafka.Message{
		Topic: topicEventCreated,
		Value: data,
		Key:   article.ID.String(),
	}
	if err := p.producer.Send(ctx, message); err != nil {
		return errs.FromKafkaError(err)
	}
	return nil
}

func (p *ArticleEventProducer) Updated(ctx context.Context, article entities.Article) error {
	data, err := json.Marshal(article)
	if err != nil {
		return err
	}
	message := &kafka.Message{
		Topic: topicEventUpdated,
		Value: data,
		Key:   article.ID.String(),
	}
	if err := p.producer.Send(ctx, message); err != nil {
		return errs.FromKafkaError(err)
	}
	return nil
}

func (p *ArticleEventProducer) Deleted(ctx context.Context, id uuid.UUID) error {
	message := &kafka.Message{
		Topic: topicEventDeleted,
		Value: []byte(id.String()),
		Key:   id.String(),
	}
	if err := p.producer.Send(ctx, message); err != nil {
		return errs.FromKafkaError(err)
	}
	return nil
}
