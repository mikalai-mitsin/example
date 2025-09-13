package repositories

import (
	"context"

	entities "github.com/mikalai-mitsin/example/internal/app/posts/entities/like"
	"github.com/mikalai-mitsin/example/internal/pkg/errs"
	"github.com/mikalai-mitsin/example/internal/pkg/kafka"
	"google.golang.org/protobuf/proto"
)

const (
	topicName = "example.posts.like.v1"
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

func (p *LikeEventProducer) Send(ctx context.Context, like entities.Like) error {
	data, err := proto.Marshal(decodeLike(like))
	if err != nil {
		return err
	}
	message := &kafka.Message{
		Topic: topicName,
		Value: data,
		Key:   like.ID.String(),
	}
	if err := p.producer.Send(ctx, message); err != nil {
		return errs.FromKafkaError(err)
	}
	return nil
}
