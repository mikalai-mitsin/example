package events

import (
	"context"
	"encoding/json"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/dtx"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"github.com/mikalai-mitsin/example/internal/pkg/uuid"
)

const (
	topicEventCreated = "example.posts.tag.created"
	topicEventUpdated = "example.posts.tag.updated"
	topicEventDeleted = "example.posts.tag.deleted"
)

type TagEventProducer struct {
	producer producer
	logger   logger
}

func NewTagEventProducer(
	producer producer,
	logger logger,
) *TagEventProducer {
	return &TagEventProducer{producer: producer, logger: logger}
}

func (p *TagEventProducer) Created(ctx context.Context, _ dtx.TX, tag entities.Tag) error {
	data, err := json.Marshal(tag)
	if err != nil {
		return err
	}
	message := &kafka.Message{
		Topic: topicEventCreated,
		Value: data,
		Key:   tag.ID.String(),
	}
	if err := p.producer.Send(ctx, message); err != nil {
		return errs.FromKafkaError(err)
	}
	return nil
}

func (p *TagEventProducer) Updated(ctx context.Context, _ dtx.TX, tag entities.Tag) error {
	data, err := json.Marshal(tag)
	if err != nil {
		return err
	}
	message := &kafka.Message{
		Topic: topicEventUpdated,
		Value: data,
		Key:   tag.ID.String(),
	}
	if err := p.producer.Send(ctx, message); err != nil {
		return errs.FromKafkaError(err)
	}
	return nil
}

func (p *TagEventProducer) Deleted(ctx context.Context, _ dtx.TX, id uuid.UUID) error {
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
