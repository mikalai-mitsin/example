package events

import (
	"context"
	"encoding/json"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

const (
	topicEventCreated = "example.posts.like.created"
	topicEventUpdated = "example.posts.like.updated"
	topicEventDeleted = "example.posts.like.deleted"
)

type LikeEventProducer struct {
	producer producer
	logger   logger
}

func NewLikeEventProducer(
	producer producer,
	logger logger,
) *LikeEventProducer {
	return &LikeEventProducer{producer: producer, logger: logger}
}

func (p *LikeEventProducer) Created(ctx context.Context, like entities.Like) error {
	data, err := json.Marshal(like)
	if err != nil {
		return err
	}
	message := &kafka.Message{
		Topic: topicEventCreated,
		Value: data,
		Key:   like.ID.String(),
	}
	if err := p.producer.Send(ctx, message); err != nil {
		return errs.FromKafkaError(err)
	}
	return nil
}

func (p *LikeEventProducer) Updated(ctx context.Context, like entities.Like) error {
	data, err := json.Marshal(like)
	if err != nil {
		return err
	}
	message := &kafka.Message{
		Topic: topicEventUpdated,
		Value: data,
		Key:   like.ID.String(),
	}
	if err := p.producer.Send(ctx, message); err != nil {
		return errs.FromKafkaError(err)
	}
	return nil
}

func (p *LikeEventProducer) Deleted(ctx context.Context, id uuid.UUID) error {
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
