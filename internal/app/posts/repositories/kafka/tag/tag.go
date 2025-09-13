package repositories

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/tag"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"google.golang.org/protobuf/proto"
)

const (
	topicName = "example.posts.tag.v1"
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

func (p *TagEventProducer) Send(ctx context.Context, tag entities.Tag) error {
	data, err := proto.Marshal(decodeTag(tag))
	if err != nil {
		return err
	}
	message := &kafka.Message{
		Topic: topicName,
		Value: data,
		Key:   tag.ID.String(),
	}
	if err := p.producer.Send(ctx, message); err != nil {
		return errs.FromKafkaError(err)
	}
	return nil
}
