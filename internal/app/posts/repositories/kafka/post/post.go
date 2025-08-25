package events

import (
	"context"
	"encoding/json"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/post"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

const (
	topicEventCreated = "example.posts.post.created"
	topicEventUpdated = "example.posts.post.updated"
	topicEventDeleted = "example.posts.post.deleted"
)

type PostEventProducer struct {
	producer producer
	logger   logger
}

func NewPostEventProducer(
	producer producer,
	logger logger,
) *PostEventProducer {
	return &PostEventProducer{producer: producer, logger: logger}
}

func (p *PostEventProducer) Created(ctx context.Context, post entities.Post) error {
	data, err := json.Marshal(post)
	if err != nil {
		return err
	}
	message := &kafka.Message{
		Topic: topicEventCreated,
		Value: data,
	}
	if err := p.producer.Send(ctx, message); err != nil {
		return errs.FromKafkaError(err)
	}
	return nil
}

func (p *PostEventProducer) Updated(ctx context.Context, post entities.Post) error {
	data, err := json.Marshal(post)
	if err != nil {
		return err
	}
	message := &kafka.Message{
		Topic: topicEventUpdated,
		Value: data,
	}
	if err := p.producer.Send(ctx, message); err != nil {
		return errs.FromKafkaError(err)
	}
	return nil
}

func (p *PostEventProducer) Deleted(ctx context.Context, id uuid.UUID) error {
	message := &kafka.Message{
		Topic: topicEventDeleted,
		Value: []byte(id.String()),
	}
	if err := p.producer.Send(ctx, message); err != nil {
		return errs.FromKafkaError(err)
	}
	return nil
}
